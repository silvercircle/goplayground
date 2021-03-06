package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"testgo/lib"
	"testgo/lib/db"
	"testgo/lib/lang"
	"time"
)

var store = sessions.NewCookieStore([]byte("adadfasdfadafsd"))
var rescanLock sync.RWMutex

type FCGISrv struct{}

func (s FCGISrv) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	HandleRequest(resp, req)
}

var mustRescanTemplates = false

/*
 * the main function that handle each request. It does:
 * o) set up the lib.Data that holds all per-request specific data
 * o) connect to the database
 * o) determines what to do, based on the matched route
 * o) outputs all templates that have been registered in Data.Templates
 */
func HandleRequest(resp http.ResponseWriter, req *http.Request) {
	var rdat lib.Data
	var dberr error

	if mustRescanTemplates {
		rescanLock.Lock()
		if mustRescanTemplates {
			mustRescanTemplates = false
			lib.PreloadTemplates("", true)
		}
		rescanLock.Unlock()
	}
	rdat.BeginRequest = time.Now()
	req.ParseForm()
	rdat.TheDB, dberr = db.Connect()
	rdat.Context = map[string]interface{}{}
	rdat.StringContext = map[string]string{}
	rdat.Out = resp
	rdat.Req = req
	rdat.Lang = lang.Languages[lib.SysConf.Settings.Settings.Language]
	rdat.Context["L"] = rdat.Lang
	rdat.Context["C"] = rdat.StringContext

	if req.FormValue("xml") != "" {
		rdat.ResponseTypeXML = true
	} else if req.FormValue("json") != "" {
		rdat.ResponseTypeJSON = true
	}

	if ok := lib.SysConf.Router.Match(req, &rdat.RouteMatch); ok {
		rdat.Context["matched_route"] = rdat.RouteMatch.Route.GetName()
		rdat.CurrentRoute = rdat.Context["matched_route"].(string)
	} else {
		rdat.Context["matched_route"] = "no match"
		rdat.CurrentRoute = "index"
	}

	if req.Proto[0:5] == "HTTPS" {
		rdat.BaseURL = "https://" + req.Host + lib.SysConf.Settings.Server.URLPrefix
	} else {
		rdat.BaseURL = "http://" + req.Host + lib.SysConf.Settings.Server.URLPrefix
	}
	rdat.Session, _ = store.Get(req, "appsession")

	// we were not able to connect to our database, inform the user and bail out early
	if dberr != nil || rdat.TheDB == nil {
		//t, _ := lib.LoadTemplate("errors/dberror")
		rdat.Templates = append(rdat.Templates, "errors/dberror")
		rdat.Context["dberror"] = dberr.Error()
	} else {
		defer rdat.TheDB.Close()
		rdat.Dispatch()
	}
	// must be done before content is sent
	rdat.Session.Save(req, resp)
	// output any optional http headers
	if rdat.ResponseTypeXML {
		resp.Header().Add("Content-Type", "text/xml; charset=UTF-8")
	} else if rdat.ResponseTypeJSON {
		resp.Header().Add("Content-Type", "Application/json; charset=UTF-8")
	} else {
		resp.Header().Add("Content-Type", "text/html; charset=UTF-8")
	}
	if v, ok := rdat.Context["httpheaders"]; ok {
		for key, value := range v.(map[string]string) {
			resp.Header().Add(key, value)
		}
	}

	// begin content output. If no header/footer templates have been loaded
	// by the handler, do this now and output the header template. Unless it's XML,
	// then just output the XML header.
	if !rdat.ResponseTypeXML && !rdat.ResponseTypeJSON {
		if rdat.HeaderTemplate == "" {
			//rdat.HeaderTemplate, _ = lib.LoadTemplate("header")
			rdat.HeaderTemplate = "header"
		}
		if rdat.FooterTemplate == "" {
			//rdat.FooterTemplate, _ = lib.LoadTemplate("footer")
			rdat.FooterTemplate = "footer"
		}
		lib.T.ExecuteTemplate(resp, rdat.HeaderTemplate, rdat.Context)
	} else if rdat.ResponseTypeXML {
		resp.Write([]byte(`<?xml version="1.0" encoding="UTF-8" ?>`))
	}
	// output all templates that were registered by the handler(s) (if any)
	for _, t := range rdat.Templates {
		lib.T.ExecuteTemplate(resp, t, rdat.Context)
	}
	// and finally the footer (but again, not for xml
	rdat.EndRequest = time.Now()
	rdat.Context["loadtime"] = fmt.Sprintf("%.2f", float32(rdat.EndRequest.Sub(rdat.BeginRequest))/1000000.0) + "ms"
	if !rdat.ResponseTypeXML && !rdat.ResponseTypeJSON {
		lib.T.ExecuteTemplate(resp, rdat.FooterTemplate, rdat.Context)
	}
}

func main() {
	var srv *FCGISrv
	var listener net.Listener

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	store.Options = &sessions.Options{
		Path:   "/app",
		MaxAge: 86400 * 365,
	}

	lib.SysConf.Homepath = dir
	lib.SysConf.Settings.Init()
	lang.Init()

	if lib.SysConf.Settings.Server.URLPrefix != "" {
		lib.SysConf.Router = mux.NewRouter().PathPrefix(lib.SysConf.Settings.Server.URLPrefix).Subrouter()
	} else {
		lib.SysConf.Router = mux.NewRouter()
	}
	lib.SysConf.Router.StrictSlash(true)
	lib.SysConf.Router.Path("/test/{id:[0-9]+}").Name("testroute")
	lib.SysConf.Router.Path("/index").Name("index")
	lib.SysConf.Router.Path("/profile.{id:[0-9]+}").Name("profile")

	//lib.Templates = map[string]*lib.TemplateNode{}

	if _, err := os.Stat(lib.SysConf.Settings.Server.FCGISock); err == nil {
		fmt.Println("Stale socket found, removing")
		os.Remove(lib.SysConf.Settings.Server.FCGISock)
	}
	lib.PreloadTemplates("", true)

	if lib.SysConf.Settings.Server.FCGI {
		listener, _ = net.Listen("unix", lib.SysConf.Settings.Server.FCGISock)
		defer listener.Close()
		os.Chmod(lib.SysConf.Settings.Server.FCGISock, 0777)
		srv = new(FCGISrv)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for sig := range c {
			log.Printf("captured %v, exiting", sig)
			if lib.SysConf.Settings.Server.FCGI {
				listener.Close()
				os.Remove(lib.SysConf.Settings.Server.FCGISock)
			}
			os.Exit(1)
		}
	}()
	go func() {
		var currentMostRecent int64
		for {
			time.Sleep(lib.SysConf.Settings.Templates.MTimeCheckInterval * 1000 * time.Millisecond)
			if !mustRescanTemplates {
				currentMostRecent = lib.SysConf.TemplateMostRecent
				lib.TemplateLock.Lock()
				lib.DoPreloadTemplates("", true)
				lib.TemplateLock.Unlock()
				if lib.SysConf.TemplateMostRecent > currentMostRecent {
					rescanLock.Lock()
					mustRescanTemplates = true
					rescanLock.Unlock()
				}
			}
		}
	}()
	if lib.SysConf.Settings.Server.FCGI {
		fcgi.Serve(listener, srv)
	} else {
		http.HandleFunc("/", HandleRequest)
		http.ListenAndServe(lib.SysConf.Settings.Server.HttpHost+":"+lib.SysConf.Settings.Server.HttpPort, nil)
	}
}

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
    "syscall"
    "testgo/lib"
    "testgo/lib/db"
    "testgo/lib/lang"
    "time"
)

var store = sessions.NewCookieStore([]byte("adadfasdfadafsd"))

type FCGISrv struct{}

func (s FCGISrv) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
    HandleRequest(resp, req)
}

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

    rdat.BeginRequest = time.Now()
    req.ParseForm()
    rdat.TheDB, dberr = db.Connect()
    rdat.Context = map[string]interface{}{}
    rdat.Out = resp
    rdat.Req = req
    rdat.Lang = lang.Languages["default"]
    rdat.Context["L"] = rdat.Lang

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

    session, _ := store.Get(req, "appsession")

    if dberr != nil || rdat.TheDB == nil {
        t, _ := lib.LoadTemplate("errors/dberror")
        rdat.Templates = append(rdat.Templates, t)
        rdat.Context["dberror"] = dberr.Error()
    } else {
        defer rdat.TheDB.Close()
        rdat.Dispatch()
    }
    // must be done before content is sent
    session.Save(req, resp)
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
        if rdat.HeaderTemplate == nil {
            rdat.HeaderTemplate, _ = lib.LoadTemplate("header")
        }
        if rdat.FooterTemplate == nil {
            rdat.FooterTemplate, _ = lib.LoadTemplate("footer")
        }
        rdat.HeaderTemplate.Execute(resp, rdat.Context)
    } else if rdat.ResponseTypeXML {
        resp.Write([]byte(`<?xml version="1.0" encoding="UTF-8" ?>`))
    }
    // output all templates that were registered by the handler(s) (if any)
    for _, t := range rdat.Templates {
        t.Execute(resp, rdat.Context)
    }
    // and finally the footer (but again, not for xml
    rdat.EndRequest = time.Now()
    rdat.Context["loadtime"] = rdat.EndRequest.Sub(rdat.BeginRequest)
    if !rdat.ResponseTypeXML && !rdat.ResponseTypeJSON {
        rdat.FooterTemplate.Execute(resp, rdat.Context)
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

    lib.Templates = map[string]*lib.TemplateNode{}

    if _, err := os.Stat(lib.SysConf.Settings.Server.FCGISock); err == nil {
        fmt.Println("Stale socket found, removing")
        os.Remove(lib.SysConf.Settings.Server.FCGISock)
    }
    lib.PreloadTemplates("")

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
    if lib.SysConf.Settings.Server.FCGI {
        fcgi.Serve(listener, srv)
    } else {
        http.HandleFunc("/", HandleRequest)
        http.ListenAndServe(lib.SysConf.Settings.Server.HttpHost+":"+lib.SysConf.Settings.Server.HttpPort, nil)
    }
}

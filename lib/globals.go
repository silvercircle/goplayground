package lib

import (
    "github.com/gorilla/mux"
    "github.com/jmoiron/sqlx"
	"github.com/gorilla/sessions"    
    "io"
    "net/http"
    "time"
)

// global app configuration object
var SysConf struct {
    Homepath string								// where we reside in the filesystem (needed to build relative paths)
    Router   *mux.Router						// the router object, holding all known routes
    Settings Settings							// Settings object, populated with defaults and overrides from appconfig.ini
    TemplateMostRecent int64
}

// per request private data
type Data struct {
    TheDB            *sqlx.DB					// db connection
    Context          map[string]interface{}		// request context
    Out              io.Writer					// response must go there...
    Req              *http.Request				// the request 
    Lang             map[string]string			// points to the language table
    Templates        []string					// templates loaded during request processing (HandleRequest() must output them)
    HeaderTemplate   string						// allows for custom header and footer template(s) for this particular request
    FooterTemplate   string						// if left empty, HandleRequest() will use the default ones
    ResponseTypeXML  bool						// request has the xml parameter set - response should be xml, used for Ajax requests
    ResponseTypeJSON bool						// request has the json parameter set - response should be json
    RouteMatch       mux.RouteMatch			 	// used to match the current route
    CurrentRoute     string						// name of current route (if any). defaults to "index"
    BeginRequest     time.Time					// for timing the request
    EndRequest       time.Time
    Session			 *sessions.Session			// session object
}

package lib

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"io"
	"net/http"
	"time"
)

// global app configuration object
var SysConf struct {
	Homepath           string      // where we reside in the filesystem (needed to build relative paths)
	Router             *mux.Router // the router object, holding all known routes
	Settings           Settings    // Settings object, populated with defaults and overrides from appconfig.ini
	TemplateMostRecent int64
}

// per request private data
type Data struct {
	TheDB            *sqlx.DB               // db connection
	// The output context is organized as some kind of namespace. The root node may
	// hold any type of data - strings, arrays, structs, arrays of structs, maps
	// anything that is required.
	// sub-namespaces exist. Context["L"] holds the current language definition 
	// (a map of strings). Context["C"] is another map of strings, mainly for generic
	// data. Context["S"] holds the current session and Context["U"] the logged-in
	// user (if any, otherwise its a default struct for a guest user)
	Context          map[string]interface{} // request context
	// string context is a shortcut and is mapped to Context["C"]
	StringContext	 map[string]string
	Out              io.Writer              // response must go there...
	Req              *http.Request          // the request
	// Lang is one of Languages[] a map of strings. For template output
	// it is mapped to Context["L"]
	Lang             map[string]string      // points to the language table
	Templates        []string               // templates loaded during request processing (HandleRequest() must output them)
	HeaderTemplate   string                 // allows for custom header and footer template(s) for this particular request
	FooterTemplate   string                 // if left empty, HandleRequest() will use the default ones
	ResponseTypeXML  bool                   // request has the xml parameter set - response should be xml, used for Ajax requests
	ResponseTypeJSON bool                   // request has the json parameter set - response should be json
	RouteMatch       mux.RouteMatch         // used to match the current route
	CurrentRoute     string                 // name of current route (if any). defaults to "index"
	BeginRequest     time.Time              // for timing the request
	EndRequest       time.Time
	Session          *sessions.Session // session object
	BaseURL			 string
}

package lib

import (
    "github.com/gorilla/mux"
    "github.com/jmoiron/sqlx"
    "html/template"
    "io"
    "net/http"
    "time"
)

var SysConf struct {
    TheDB    *sqlx.DB
    Homepath string
    Router   *mux.Router
    Settings Settings
}

type Data struct {
    TheDB          *sqlx.DB
    Context        map[string]interface{}
    Out            io.Writer
    Req            *http.Request
    Lang           map[string]string
    Templates      []*template.Template
    HeaderTemplate *template.Template
    FooterTemplate *template.Template
    ResponseTypeXML bool
    ResponseTypeJSON bool
    RouteMatch		mux.RouteMatch
    CurrentRoute	string
    BeginRequest	time.Time
    EndRequest		time.Time
}

type Ctx map[string]interface{}

var Globalvar = 10

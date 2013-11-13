package lib

import (
    "code.google.com/p/gcfg"
    "os"
)

type Settings struct {
	Server struct {
		FCGI		 bool
		FCGISock	 string
		HttpPort	 string
		HttpHost	 string
		URLPrefix    string
	}
    Database struct {
        DBUser       string
        DBPass       string
        DBName       string
        DBHost       string
        DBPort       string
        DBMethod     string
        DBSocketname string
    }
    Templates struct {
    	Preload		 bool
    }
}

func (s *Settings) Init() {
    s.Database.DBUser = "dbuser"
    s.Database.DBPass = "dbpass"
    s.Database.DBName = "golang"
    s.Database.DBHost = "localhost"
    s.Database.DBPort = "3306"
    s.Database.DBMethod = "unix"
    s.Database.DBSocketname = "/var/run/mysqld/mysqld.sock"
	
	s.Templates.Preload = false
	
	s.Server.FCGI = true
	s.Server.FCGISock = SysConf.Homepath + "/app.sock"
	s.Server.HttpPort = "8080"
	s.Server.HttpHost = ""
	s.Server.URLPrefix = ""
	
    if _, err := os.Lstat(SysConf.Homepath + "/appconfig.ini"); err == nil {
        gcfg.ReadFileInto(s, SysConf.Homepath+"/appconfig.ini")
    }
}

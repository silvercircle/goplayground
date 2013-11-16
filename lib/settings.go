// implement app settings
package lib

import (
	"code.google.com/p/gcfg"
	"os"
	"time"
	"log"
)

type Settings struct {
	Server struct {
		FCGI      bool
		FCGISock  string
		HttpPort  string
		HttpHost  string
		URLPrefix string
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
		MTimeCheckInterval time.Duration
	}
	Settings struct {
		Language string
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

	// check every 60 *seconds* for changed templates
	// can be customized via appconfig.ini for a more convenient value on a dev system
	s.Templates.MTimeCheckInterval = 60 

	s.Server.FCGI = true
	s.Server.FCGISock = SysConf.Homepath + "/app.sock"
	s.Server.HttpPort = "8080"
	s.Server.HttpHost = ""
	s.Server.URLPrefix = ""

	s.Settings.Language = "default"

	if _, err := os.Lstat(SysConf.Homepath + "/appconfig.ini"); err == nil {
		err = gcfg.ReadFileInto(s, SysConf.Homepath+"/appconfig.ini")
		if err != nil {
			log.Fatal(err)
		}
	}
	
	if s.Server.URLPrefix == "" {
		s.Server.URLPrefix = "/"
	}
}

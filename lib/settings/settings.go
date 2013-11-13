package settings

import (
    "fmt"
    "os"
    "testgo/lib"
)

type Settings struct {
    DBUser string
    DBPass string
    DBName string
}

func (s *Settings) Init() {
    s.DBUser = "berl"
    s.DBPass = "idefix9"
    s.DBName = "berl_shop"
    if fi, err := os.Lstat(SysConf.Homepath + "/appconfig.ini"); !err {

    } else {
        fmt.Println("Config file not found")
    }
}

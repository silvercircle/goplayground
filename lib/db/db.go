package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"testgo/lib"
)

func Connect() (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	if lib.SysConf.Settings.Database.DBMethod == "unix" {
		db, err = sqlx.Open("mysql", lib.SysConf.Settings.Database.DBUser+":"+lib.SysConf.Settings.Database.DBPass+"@unix("+lib.SysConf.Settings.Database.DBSocketname+")/"+lib.SysConf.Settings.Database.DBName)
	} else {
		db, err = sqlx.Open("mysql", lib.SysConf.Settings.Database.DBUser+":"+lib.SysConf.Settings.Database.DBPass+"@tcp("+lib.SysConf.Settings.Database.DBHost+":"+lib.SysConf.Settings.Database.DBPort+")/"+lib.SysConf.Settings.Database.DBName)
	}

	err = db.Ping()
	return db, err
}

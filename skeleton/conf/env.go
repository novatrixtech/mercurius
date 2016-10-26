package conf

import (
	"github.com/jmoiron/sqlx"
	"gopkg.in/ini.v1"
	"gopkg.in/macaron.v1"
)

type Database interface {
	DSN() string
	DB() (*sqlx.DB, error)
}

var Cfg *ini.File

func init() {
	var err error
	Cfg, err = macaron.SetConfig("conf/app.ini")
	if err != nil {
		panic(err)
	}
}

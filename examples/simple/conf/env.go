package conf

import (
	"gopkg.in/ini.v1"
	"gopkg.in/macaron.v1"
)

/*
Cfg - configuration file content
*/
var Cfg *ini.File

func init() {
	var err error
	Cfg, err = macaron.SetConfig("conf/app.ini")
	if err != nil {
		panic(err)
	}
}

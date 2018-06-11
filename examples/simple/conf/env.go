package conf

import (
	"log"

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
		if isDbConnParamsInEnvVariables() {
			log.Printf("[conf/Init] Error during app.ini reading. Error: %s\n", err.Error())
		} else {
			log.Fatalf("[conf/Init] Error during app.ini reading. Error: %s\n", err.Error())
		}
	}
}

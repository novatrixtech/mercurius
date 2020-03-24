package cache

import (
	"log"
	"os"
	"strconv"

	"github.com/novatrixtech/mercurius/examples/simple/conf"
)

//GetEnvironmentParameter get environment parameter from OS variables or from ini file
func GetEnvironmentParameter(paramName string) (data interface{}) {
	var tmp string
	tmp = os.Getenv(paramName)
	if len(tmp) > 0 {
		data = tmp
		return
	}
	tmp = conf.Cfg.Section("").Key(paramName).Value()
	data = tmp
	return
}

//GetEnvironmentParameterString get string environment parameter from OS variables or from ini file
func GetEnvironmentParameterString(paramName string) (data string) {
	var tmp interface{}
	tmp = GetEnvironmentParameter(paramName)
	defer func() {
		if err := recover(); err != nil {
			log.Println("GetEnvironmentParameterString - panic occurred:")
			return
		}
	}()
	data = tmp.(string)
	return
}

//GetEnvironmentParameterInt get int environment parameter from OS variables or from ini file
func GetEnvironmentParameterInt(paramName string) (data int) {
	var tmp interface{}
	tmp = GetEnvironmentParameter(paramName)
	defer func() {
		if err := recover(); err != nil {
			log.Println("GetEnvironmentParameterInt - panic occurred:")
			return
		}
	}()
	data, err := strconv.Atoi(tmp.(string))
	if err != nil {
		log.Println("GetEnvironmentParameterInt - ", tmp, " is not an int")
		return
	}
	return
}

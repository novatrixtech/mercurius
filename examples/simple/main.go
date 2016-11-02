// app example generated using Mercurius
package main

import (
	"gopkg.in/macaron.v1"

	config "github.com/novatrixtech/mercurius/examples/simple/conf"
	conf "github.com/novatrixtech/mercurius/examples/simple/conf/app"
	"os"
	"strconv"
)

func main() {
	app := macaron.New()
	conf.SetupMiddlewares(app)
	conf.SetupRoutes(app)
	app.Run(port())
}

func port() int {
	port, err := config.Cfg.Section("").Key("http_port").Int()
	if err != nil {
		panic(err)
	}
	if macaron.Env != macaron.DEV {
		if i, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
			port = i
		}
	}
	return port
}

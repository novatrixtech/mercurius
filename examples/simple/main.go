// app example generated using Mercurius
package main

import (
	"gopkg.in/macaron.v1"

	conf "github.com/novatrixtech/mercurius/examples/simple/conf/app"
	config "github.com/novatrixtech/mercurius/examples/simple/conf"
)

func main() {
	app := macaron.New()
	conf.SetupMiddlewares(app)
	conf.SetupRoutes(app)
	port, err := config.Cfg.Section("").Key("http_port").Int()
	if err != nil {
		panic(err)
	}
	app.Run(port)
}

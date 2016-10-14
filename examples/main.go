package main

import (
	"gopkg.in/macaron.v1"

	conf "github.com/novatrixtech/mercurius/examples/conf/app"
)

func main() {
	app := macaron.New()
	conf.SetupMiddlewares(app)
	conf.SetupRoutes(app)
	app.Run(8080)
}

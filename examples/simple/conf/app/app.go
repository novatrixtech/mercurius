package app

import (
	"github.com/go-macaron/binding"
	mcache "github.com/go-macaron/cache"
	"github.com/go-macaron/gzip"
	"github.com/go-macaron/i18n"
	"github.com/go-macaron/jade"
	"github.com/go-macaron/session"
	"github.com/go-macaron/toolbox"
	"github.com/novatrixtech/mercurius/examples/simple/conf"
	"github.com/novatrixtech/mercurius/examples/simple/handler"
	"github.com/novatrixtech/mercurius/examples/simple/lib/auth"
	"github.com/novatrixtech/mercurius/examples/simple/lib/cache"
	"github.com/novatrixtech/mercurius/examples/simple/lib/context"
	"github.com/novatrixtech/mercurius/examples/simple/lib/template"
	"gopkg.in/macaron.v1"
)

func SetupMiddlewares(app *macaron.Macaron) {
	app.Use(macaron.Logger())
	app.Use(macaron.Recovery())
	app.Use(gzip.Gziper())
	app.Use(toolbox.Toolboxer(app, toolbox.Options{
		HealthCheckers: []toolbox.HealthChecker{
			new(handler.AppChecker),
		},
	}))
	app.Use(macaron.Static("public"))
	app.Use(i18n.I18n(i18n.Options{
		Directory: "locale",
		Langs:     []string{"pt-BR", "en-US"},
		Names:     []string{"Português do Brasil", "American English"},
	}))
	app.Use(jade.Renderer(jade.Options{
		Directory: "public/templates",
		Funcs:     template.FuncMaps(),
	}))
	app.Use(macaron.Renderer(macaron.RenderOptions{
		Directory: "public/templates",
		Funcs:     template.FuncMaps(),
	}))
	app.Use(mcache.Cacher(
		cache.Option(conf.Cfg.Section("").Key("cache_adapter").Value()),
	))
	app.Use(session.Sessioner())
	app.Use(context.Contexter())
}

func SetupRoutes(app *macaron.Macaron) {
	app.Group("", func() {
		app.Get("/", handler.ListAccessPage)
		app.Get("/list", handler.ListAccessBy)
		app.Get("/logout", handler.Logout)
	}, auth.LoginRequired)
	app.Get("/login", handler.LoginPage)
	app.Post("/login", binding.BindIgnErr(context.Login{}), handler.BasicAuth)
	app.Group("/api/v1", func() {
		app.Post("/oauth/token", handler.Oauth)
		app.Get("/list", auth.LoginRequiredApi, handler.ListAccessForApi)
	})

}

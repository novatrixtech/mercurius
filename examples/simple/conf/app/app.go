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
	"github.com/novatrixtech/mercurius/examples/simple/lib/contx"
	"github.com/novatrixtech/mercurius/examples/simple/lib/cors"
	"github.com/novatrixtech/mercurius/examples/simple/lib/template"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/macaron.v1"
)

//SetupMiddlewares configures the middlewares using in each web request
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
	//cache in memory
	app.Use(mcache.Cacher(
		cache.Option(conf.Cfg.Section("").Key("cache_adapter").Value()),
	))
	/*
		Redis Cache
		Add this lib to import session: _ "github.com/go-macaron/cache/redis"
		Later replaces the cache in memory instructions for the lines below:
		optCache := mcache.Options{
				Adapter:       conf.Cfg.Section("").Key("cache_adapter").Value(),
				AdapterConfig: conf.Cfg.Section("").Key("cache_adapter_config").Value(),
			}
		app.Use(mcache.Cacher(optCache))
	*/
	app.Use(session.Sessioner())
	app.Use(contx.Contexter())
	app.Use(cors.Cors())
}

//SetupRoutes defines the routes the Web Application will respond
func SetupRoutes(app *macaron.Macaron) {
	app.Group("", func() {
		app.Get("/", handler.ListAccessPage)
		app.Get("/list", handler.ListAccessBy)
		app.Get("/logout", handler.Logout)
	}, auth.LoginRequired)
	//HealthChecker
	app.Get("/health", handler.HealthCheck)
	app.Get("/metrics", promhttp.Handler())
	app.Get("/login", handler.LoginPage)
	app.Post("/login", binding.BindIgnErr(contx.Login{}), handler.BasicAuth)
	app.Group("/api/v1", func() {
		app.Post("/oauth/token", handler.Oauth)
		app.Get("/list", auth.LoginRequiredApi, handler.ListAccessForApi)
	})
	/*
		//An example to test DB connection
		app.Get("", func() string {
			db, err := conf.GetDB()
			if err != nil {
				return err.Error()
			}
			err = db.Ping()
			if err != nil {
				return err.Error()
			}
			col, err := conf.GetMongoCollection("teste")
			if err != nil {
				return err.Error()
			}
			defer col.Database.Session.Close()
			teste := Teste{Status: "OK"}
			err = col.Insert(teste)
			return "Mercurius Works!"
		})

		//Include this struct after import session
		type Teste struct {
			Status string
		}
	*/
}

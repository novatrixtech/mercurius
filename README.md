# mercurius
Mercurius gives you speed when create 'Go' applications. It lets you being focused at business.
<br /><br />
Get a Go web application template with: Internationalization, Routers, Logging, Cache, Database, Jade Template Render Engine, JWT, oAuth 2.0. Built upon Macaron, all those items are all configured and ready to use.

# Getting Started

```go get -v github.com/novatrixtech/mercurius```

```go install github.com/novatrixtech/mercurius```

```mercurius new github.com/novatrixtech/newproject```

# Structure of a Project
```
/conf 
Application configuration including environment-specific configs

/conf/app
Middlewares and routes configuration

/handler
HTTP handlers

/locale
Language specific content bundles

/lib
Common libraries to be used across your app

/model
Models

/public
Web resources that are publicly available

/public/templates
Jade templates

/repository
Database comunication following repository pattern

main.go
Application entry
```
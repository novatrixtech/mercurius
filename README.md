# mercurius
[![GoDoc](https://godoc.org/github.com/novatrixtech/mercurius?status.svg)](https://godoc.org/github.com/novatrixtech/mercurius)
[![license](https://img.shields.io/badge/license-Apache%20License%202.0-blue.svg?style=flat)](http://www.apache.org/licenses/LICENSE-2.0)

Mercurius gives you speed to create new 'Go' applications. It let you more focused on your business than in your backend.
<br/><br/>
Get a web application template for Golang that includes: Internationalization, Routers, Logging, Cache, Database, Jade Template Render Engine, JWT, oAuth 2.0. Built on top of Macaron, all items are configured and ready to use.

# Getting Started

```go get -v github.com/novatrixtech/mercurius/...```

```go install github.com/novatrixtech/mercurius```

```mercurius new```

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

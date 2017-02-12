package context

import (
	"github.com/go-macaron/cache"
	"github.com/go-macaron/jade"
	"github.com/go-macaron/session"
	"github.com/go-macaron/toolbox"
	"gopkg.in/macaron.v1"
	"net/http"
	"encoding/json"
    "log"
)

var ctx *Context

// Context representation
type Context struct {
	*macaron.Context
	render  jade.Render
	Session session.Store
	Flash   *session.Flash
	Cache   cache.Cache
	Toolbox toolbox.Toolbox
}

// HasError return true if request has errors
func (ctx *Context) HasError() bool {
	hasErr, ok := ctx.Data["HasError"]
	if !ok {
		return false
	}
	ctx.Flash.ErrorMsg = ctx.Data["ErrorMsg"].(string)
	ctx.Data["flash"] = ctx.Flash
	return hasErr.(bool)
}

func (ctx *Context) withErr(msg string, userForm interface{}) {
	if userForm != nil {
		AssignForm(userForm, ctx.Data)
	}
	ctx.Flash.ErrorMsg = msg
	ctx.Data["flash"] = ctx.Flash
}

// RenderWithErr render view and add error message using jade
func (ctx *Context) RenderWithErr(msg string, tpl string, userForm interface{}) {
	ctx.withErr(msg, userForm)
	ctx.HTML(http.StatusOK, tpl)
}
 
// NativeRenderWithErr render view and add error message using Go engine
func (ctx *Context) NativeRenderWithErr(msg string, tpl string, userForm interface{}) {
	ctx.withErr(msg, userForm)
	ctx.NativeHTML(http.StatusOK, tpl)
}

// Contexter middleware
func Contexter() macaron.Handler {
	return func(c *macaron.Context, r jade.Render, session session.Store, flash *session.Flash, cache cache.Cache, toolbox toolbox.Toolbox) {
		ctx = &Context{
			Context: c,
			render:  r,
			Session: session,
			Flash:   flash,
			Cache:   cache,
			Toolbox: toolbox,
		}
		c.Map(ctx)
	}
}

// HTML render using jade
func (ctx *Context) HTML(status int, name string) {
	ctx.render.HTML(status, name, ctx.Data)
}

// NativeHTML render using go engine
func (ctx *Context) NativeHTML(status int, name string) {
	ctx.Context.HTML(status, name, ctx.Data)
}

// JSONWithoutEscape render json without escape
func (ctx *Context) JSONWithoutEscape(status int, obj interface{}) {
	ctx.Header().Set("Content-Type", "application/json")
	ret, err := json.Marshal(&obj)
	if err != nil {
		log.Print("[JSONWithoutEscape]" + err.Error())
		http.Error(ctx.Resp, "{'errors':'JSON Marshaling Error = "+err.Error()+"'}", 500)
		return
	}
	ctx.Status(status)
	log.Println("[JSONWithoutEscape] Returned object: " + string(ret))
	ctx.Resp.Write(ret)
}

// I18n view func
func I18n(key string) string {
	return ctx.Tr(key)

}

/*
GetContext Get system context
*/
func GetContext() *Context {
    return ctx
}


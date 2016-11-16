package handler

import (
	"net/http"

	"github.com/novatrixtech/mercurius/examples/simple/conf"
	"github.com/novatrixtech/mercurius/examples/simple/lib/auth"
	"github.com/novatrixtech/mercurius/examples/simple/lib/context"
)

func LoginPage(ctx *context.Context) {
	ctx.HTML(http.StatusOK, "login")
}

func BasicAuth(ctx *context.Context, form context.Login) {
	if ctx.HasError() {
		ctx.HTML(http.StatusOK, "login")
		return
	}
	if form.Username == "maquinataxi" && form.Password == "maquina1234" {
		auth.CreateJWTCookie(form.Username, "", ctx)
		ctx.Redirect("/")
		return
	} else {
		ctx.RenderWithErr(ctx.Tr("auth_fail"), "login", &form)
	}
}

func Oauth(ctx *context.Context) {
	key := conf.Cfg.Section("").Key("oauth_key").Value()
	id, secret, _ := ctx.Req.BasicAuth()
	credentials := auth.Oauth{
		Id:     id,
		Secret: secret,
	}
	appName, appID, _ := auth.ClientDecrypter(key, id, secret)
	data := auth.DB[credentials]
	if data.Id == appID && data.Name == appName {
		token := auth.GenerateJWTToken(data, ctx)
		ctx.JSON(http.StatusOK, map[string]string{"token": token})
		return
	}
	ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid Credentials"})
	return

}

func Logout(ctx *context.Context) {
	auth.InvalidateJWTToken(ctx)
	ctx.Redirect("/login")
}

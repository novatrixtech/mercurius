package auth

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/novatrixtech/mercurius/examples/lib/context"
)

const cookie_name = "mercuriusAuth"

type Oauth struct {
	Id     string `json:"id"`
	Secret string `json:"secret"`
}

type Claims struct {
	Ip string `json:"ip"`
	jwt.StandardClaims
}

type App struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

var DB map[Oauth]*App = make(map[Oauth]*App)

func LoginRequired(ctx *context.Context) {
	cookie, err := ctx.Req.Cookie(cookie_name)
	if err != nil {
		ctx.Redirect("/login")
		return
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, parse)
	if err != nil {
		ctx.Redirect("/login")
		return
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid && ctx.RemoteAddr() == claims.Ip {
		ctx.Data["jwt"] = *claims
	} else {
		ctx.Redirect("/login")
	}
}

func LoginRequiredApi(ctx *context.Context) {
	header := ctx.Req.Header.Get("Authorization")
	if header != "" {
		value := strings.Split(header, " ")[1]
		token, err := jwt.ParseWithClaims(value, &Claims{}, parse)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			return
		}
		if claims, ok := token.Claims.(*Claims); ok && token.Valid && ctx.RemoteAddr() == claims.Ip {
			ctx.Data["jwt"] = *claims
			return
		} else {
			ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			return
		}
	}
	ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
}

package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/novatrixtech/mercurius/examples/simple/lib/contx"
	"github.com/novatrixtech/mercurius/examples/simple/repo"
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

func LoginRequired(ctx *contx.Context) {
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

func LoginRequiredApi(ctx *contx.Context) {
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

/*
ValidateAuthorizationHeader checks if the Authorization Header contains an Access Token and if it is still valid
*/
func ValidateAuthorizationHeader(authHeader http.Header, funcName string, rolesAllowed []int) (contatoID int, err error) {
	err = nil
	var ac string
	items, ok := authHeader["Authorization"]
	if !ok || len(items) < 1 {
		err = errors.New("Cabeçalho de autorização não encontrado")
		return
	}
	ac = strings.TrimSpace(items[0])
	//log.Printf("[ValidateAuthorizationHeader] AccessToken Completo: [%s]\n", ac)
	if !strings.Contains(ac, "Bearer") {
		err = errors.New("Formato de cabeçalho de autorização invalido. Sem Bearer")
		return
	}
	ac = strings.TrimPrefix(ac, "Bearer ")
	ac = strings.TrimSpace(ac)
	//log.Println("[ValidateAuthorizationHeader] AccessToken: ", ac)
	acObj, ok := AccessTokenCache[ac]
	if !ok {
		err = errors.New("Access Token não encontrado no cache")
		return
	}
	allowed := false
	for _, role := range rolesAllowed {
		if acObj.RoleLevel == role {
			allowed = true
			break
		}
	}
	if !allowed {
		err = errors.New("Access Token não tem permissão")
		return
	}
	if acObj.ValidUntil < int(time.Now().Local().Unix()) {
		err = errors.New("Access Token expirado")
		return
	}
	go repo.AddAccessTokenAccessLog(ac, funcName)
	contatoID = acObj.ContatoID
	return
}

package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/novatrixtech/mercurius/examples/simple/lib/contx"
	"github.com/dgrijalva/jwt-go"
)

const cookieName = "mercurius-sample"

//Oauth estrutura de chave de autenticacao
type Oauth struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

// Claims dados a serem recuperados da chave
type Claims struct {
	IP string `json:"ip"`
	jwt.StandardClaims
}

// App aplicacao a ser validada o acesso
type App struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

//DB mapea o acesso a banco de dados que validara as credenciais
var DB = make(map[Oauth]*App)

// LoginRequired valida a existencia ou nao de cookies para acesso a telas
func LoginRequired(ctx *contx.Context) {
	cookie, err := ctx.Req.Cookie(cookieName)
	if err != nil {
		ctx.Redirect("/login")
		log.Println(err)
		return
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, parse)
	if err != nil {
		ctx.Redirect("/login")
		log.Println(err)
		return
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid && ctx.RemoteAddr() == claims.IP {
		ctx.Data["jwt"] = *claims
	} else {
		ctx.Redirect("/login")
		log.Println("Cause: Invalid token")
	}
}

// LoginRequiredAPI valida login para realizacao de chamadas de API
func LoginRequiredAPI(ctx *contx.Context) {
	header := ctx.Req.Header.Get("Authorization")
	if header != "" {
		splitted := strings.Split(header, " ")
		if len(splitted) < 2 {
			ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Malformed request header"})
			return
		}
		value := splitted[1]
		token, err := jwt.ParseWithClaims(value, &Claims{}, parse)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			return
		}
		if claims, ok := token.Claims.(*Claims); ok && token.Valid && ctx.RemoteAddr() == claims.IP {
			ctx.Data["jwt"] = *claims
			return
		}
		ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}
	ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
}

//LoginRequiredAPISystem verifica se as credenciais para executar chamadas na API no nivel sistema estao OK
func LoginRequiredAPISystem(ctx *contx.Context) {
	var err error
	rolesAllowed := []string{"12"}
	_, err = ValidateAuthorizationHeader(ctx.Req.Header, "APIETHOperations", rolesAllowed)
	if err != nil {
		log.Println("[InsereInscricaoCursos] Erro na autorização do AC: " + err.Error())
		ctx.JSON(http.StatusUnauthorized, "{'error':'Access Token invalid'}")
		return
	}
}

/*
ValidateAuthorizationHeader checks if the Authorization Header contains an Access Token and if it is still valid
*/
func ValidateAuthorizationHeader(authHeader http.Header, funcName string, rolesAllowed []string) (contatoID int, err error) {
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
	go AddAccessTokenAccessLog(ac, funcName)
	contatoID = acObj.ContatoID
	return
}

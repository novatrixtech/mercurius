package auth

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/novatrixtech/cryptonx"

	"github.com/jeffprestes/curso-go-web/conf"
	"github.com/jeffprestes/curso-go-web/lib/contx"
)

//InitializeUserCredentials stores user's credentials at OAuth Database
func InitializeUserCredentials(ctx *contx.Context) {
	body, err := ctx.Req.Body().Bytes()
	defer ctx.Req.Body().ReadCloser()
	if err != nil {
		log.Println("[GetUserCredentials] Erro ao transformar o JSON em array de bytes: " + err.Error())
		ctx.JSON(http.StatusBadRequest, "{'error':'Invalid body message'}")
		return
	}
	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("[GetUserCredentials] Erro ao fazer o binding de [%s] em objeto Request: %s\n", string(body), err.Error())
		ctx.JSON(http.StatusBadRequest, "{'error':'Invalid body message'}")
		return
	}
	user.ClientID, user.Secret, err = generateUserCredentials(user, ctx.Req.RemoteAddr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Error generating credentials")
		return
	}
	err = AddCredentialsToUser(user, "12")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Error storing credentials")
		return
	}
	ctx.JSON(http.StatusOK, user)
}

//GetAccessToken Generates OAuth Access Token
func GetAccessToken(ctx *contx.Context) {
	u, p, ok := ctx.Req.BasicAuth()
	if !ok {
		log.Println("[GetAccessToken] Deu ruim a autenticacao...")
		ctx.JSON(http.StatusUnauthorized, "")
		return
	}
	decodedClientID, err := decodeClientID(u)
	if err != nil {
		log.Printf("[GetAccessToken] Erro decodificar dados do ClientID: [%s] - Erro: [%s]\n", u, err.Error())
		ctx.JSON(http.StatusUnauthorized, "Invalid credentials")
		return
	}
	contactName, nonce, err := getDataFromClientID(decodedClientID)
	if err != nil {
		log.Printf("[GetAccessToken] Erro obter dados do ClientID: [%s] - Erro: [%s]\n", u, err.Error())
		ctx.JSON(http.StatusUnauthorized, "Invalid credentials")
		return
	}
	secretDecoded, err := decodeSecret(p, nonce)
	if err != nil {
		log.Println("[GetAccessToken] Erro ao decodar o secret. Erro: ", err.Error())
		ctx.JSON(http.StatusUnauthorized, "Invalid credentials")
		return
	}
	dataDoSecret, contatoID, IPDoSecret, err := getAndValidateDataFromSecret(secretDecoded)
	if err != nil {
		log.Println("[GetAccessToken] Erro ao obter os dados do secret. Erro: ", err.Error())
		ctx.JSON(http.StatusUnauthorized, "Invalid credentials")
		return
	}
	role, err := GetUserRoleByContactID(contatoID)
	if err != nil {
		log.Println("[GetAccessToken] Erro ao obter o role do usuario do banco de dados. Erro: ", err.Error())
		ctx.JSON(http.StatusUnauthorized, "Invalid credentials")
		return
	}
	ac := AccessTokenData{}
	ac.ContatoID = contatoID
	ac.RoleLevel = role
	ac.ValidUntil = int(time.Now().Local().Add(time.Hour * time.Duration(4)).Unix())
	ac.DateWhenSecretWasGenerated = dataDoSecret
	ac.IPUsedToGenerateSecret = IPDoSecret
	ip, _, err := net.SplitHostPort(ctx.Req.RemoteAddr)
	if err != nil {
		//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
		log.Printf("userip: %q is not IP:port", ctx.Req.RemoteAddr)
	}
	ipRemotoOrigem := net.ParseIP(ip)
	ac.IPUsedToGenerateAC = ipRemotoOrigem.String()
	_, acID, err := cryptonx.Encrypter(conf.Cfg.Section("").Key("oauth_key").Value(), contactName)
	if err != nil {
		log.Println("[GetAccessToken] Erro ao gerar Access Token: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, "{'error': 'Could not generate Access Token'}")
		return
	}
	RemoveUnusedAC(contatoID)
	AccessTokenCache[acID] = ac
	go AddAccessTokenRequestLog(acID, ac.ContatoID)
	acPub := AccessTokenPublic{}
	acPub.AccessToken = acID
	acPub.ExpiresIn = int((time.Duration(4) * time.Hour).Seconds())
	acPub.TokenType = "Bearer"
	ctx.JSON(http.StatusOK, acPub)
}

//GetOauthUserCredentials gets user's clientID and Secret
func GetOauthUserCredentials(ctx *contx.Context) {
	var err error
	rolesAllowed := []string{"3", "12"}
	_, err = ValidateAuthorizationHeader(ctx.Req.Header, "GetUserCredentials", rolesAllowed)
	if err != nil {
		log.Println("[GetUserCredentials] Erro na autorização do AC: " + err.Error())
		ctx.JSON(http.StatusUnauthorized, "{'error':'Access Token invalid'}")
		return
	}
	user := User{}
	user.ID = ctx.ParamsInt(":idclient")
	user.Name, err = GetUserNameByContactID(user.ID)
	if err != nil {
		log.Printf("[GetUserCredentials] Erro ao buscar o nome do usuario de ID [%d]. Erro: %s\n", user.ID, err.Error())
		ctx.JSON(http.StatusBadRequest, "{'error':'Invalid user ID'}")
		return
	}
	clientID, secret, err := GetUserCredentials(user)
	if err != nil {
		log.Printf("[GetUserCredentials] Erro ao obter as credenciais do usuario [%d] no banco. Erro: %s\n", user.ID, err.Error())
		ctx.JSON(http.StatusInternalServerError, "{'error': 'Could not obtain credentials'}")
		return
	}
	user.ClientID = clientID
	user.Secret = secret
	ctx.JSON(http.StatusCreated, user)
	return
}

//CheckFormUserCredentials handle user's authentication via Login Form
func CheckFormUserCredentials(ctx *contx.Context, user User) {
	sha512 := sha512.New()
	sha512.Write([]byte(user.Secret))
	password := fmt.Sprintf("%x", sha512.Sum(nil))
	user, err := GetUserCredentialsByLogin(user.ClientID, password)
	if err != nil {
		log.Println("CheckFormUserCredentials - Error: ", err.Error())
		ctx.Redirect("/login")
		return
	}
	log.Printf("CheckFormUserCredentials - user found: %+v\n", user)
	err = CreateJWTCookie(strconv.Itoa(user.ID), jwtIssuerName, 360, ctx)
	if err != nil {
		log.Println("CheckFormUserCredentials - error generating JWT Cookie - Error: ", err.Error())
		ctx.NativeHTML(http.StatusInternalServerError, "erro")
		return
	}
	ctx.Redirect("/")
}

// LoginRequired valida a existencia ou nao de cookies para acesso a telas
func LoginRequired(ctx *contx.Context) {
	cookie, err := ctx.Req.Cookie(cookieName)
	if err != nil {
		ctx.Redirect("/login")
		log.Println("LoginRequired error restoring cookie: ", err)
		return
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, parse)
	if err != nil {
		ctx.Redirect("/login")
		log.Println(err)
		return
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid && ctx.RemoteAddr() == claims.IP {
		log.Printf("LoginRequired Claims %+v\n", claims)
		intUserID, err := strconv.Atoi(claims.StandardClaims.Id)
		if err != nil {
			ctx.Redirect("/login")
			log.Println("LoginRequired error parsing cookie - user id: ", err)
			return
		}
		user, err := GetUserByID(intUserID)
		if err != nil {
			ctx.Redirect("/login")
			log.Println("LoginRequired error getting user data", err)
			return
		}
		ctx.Data["user"] = user
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

//LogoutForm handles the request to logout users that has been authenticated via HTTP form
func LogoutForm(ctx *contx.Context) {
	InvalidateJWTToken(ctx)
	ctx.Redirect("/login")
}

//IndexLogin opens login page
func IndexLogin(ctx *contx.Context) {
	ctx.NativeHTML(http.StatusOK, "login")
}

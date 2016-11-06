package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/novatrixtech/mercurius/examples/simple/conf"
	"github.com/novatrixtech/mercurius/examples/simple/lib/context"
	"net/http"
	"strings"
	"time"
)

func CreateJWTCookie(jwtID string, issuer string, ctx *context.Context) {
	ip := ctx.RemoteAddr()
	expireCookie := time.Now().Add(time.Hour * 1)
	signedToken := generateJWTToken(jwtID, ip, issuer)
	cookie := http.Cookie{Name: cookie_name, Value: signedToken, Expires: expireCookie, HttpOnly: true}
	http.SetCookie(ctx.Resp, &cookie)

}

func GenerateJWTToken(app *App, ctx *context.Context) string {
	ip := ctx.RemoteAddr()
	return generateJWTToken(app.Id, ip, app.Name)
}

func InvalidateJWTToken(ctx *context.Context) {
	deleteCookie := http.Cookie{Name: cookie_name, Value: "none", Expires: time.Now()}
	http.SetCookie(ctx.Resp, &deleteCookie)
}

func generateJWTToken(jwtID string, ip string, issuer string) string {
	expireToken := time.Now().Add(time.Hour * 1).Unix()

	if issuer == "" {
		issuer = "localhost:8080"
	}
	claims := Claims{
		Ip: ip,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    issuer,
			Id:        jwtID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(conf.Cfg.Section("").Key("oauth_key").Value()))
	return signedToken
}

func ClientDecrypter(key, clientID, clientSecret string) (name, id string, err error) {
	secret, _ := hex.DecodeString(clientSecret)
	cid, _ := hex.DecodeString(clientID)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	text, err := aesgcm.Open(nil, secret, cid, nil)
	if err != nil {
		return "", "", err
	}
	values := strings.Split(string(text), "|")
	name = values[0]
	id = values[1]
	return
}

func parse(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("Unexpected Signing method")
	}
	return []byte(conf.Cfg.Section("").Key("oauth_key").Value()), nil
}

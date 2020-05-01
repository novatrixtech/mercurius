package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

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

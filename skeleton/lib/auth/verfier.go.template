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
		err = errors.New("[ValidateAuthorizationHeader] Authorization header not found")
		return
	}
	ac = strings.TrimSpace(items[0])
	//log.Printf("[ValidateAuthorizationHeader] AccessToken Completo: [%s]\n", ac)
	if !strings.Contains(ac, "Bearer") {
		err = errors.New("[ValidateAuthorizationHeader] Authorization header invalid. No Bearer")
		return
	}
	ac = strings.TrimPrefix(ac, "Bearer ")
	ac = strings.TrimSpace(ac)
	//log.Println("[ValidateAuthorizationHeader] AccessToken: ", ac)
	acObj, ok := AccessTokenCache[ac]
	if !ok {
		err = errors.New("[ValidateAuthorizationHeader] Access Token not found in cache")
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
		err = errors.New("[ValidateAuthorizationHeader] Access Token does not have permission")
		return
	}
	if acObj.ValidUntil < int(time.Now().Local().Unix()) {
		err = errors.New("[ValidateAuthorizationHeader] Access Token is expired")
		return
	}
	go AddAccessTokenAccessLog(ac, funcName)
	contatoID = acObj.ContatoID
	return
}

package auth

const cookieName = "exemplo_cliente"
const jwtIssuerName = "exemploClienteApp"

// AccessTokenCache stores AccessToken generated and when it was generated
var AccessTokenCache map[string]AccessTokenData

func init() {
	AccessTokenCache = make(map[string]AccessTokenData, 0)
}

//RemoveUnusedAC remove from cache an Access Token generated before
func RemoveUnusedAC(contatoID int) {
	for key, ac := range AccessTokenCache {
		if ac.ContatoID == contatoID {
			delete(AccessTokenCache, key)
			break
		}
	}
}

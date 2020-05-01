package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

//AcessoAPI acesso à API
type AcessoAPI struct {
	ID           int    `json:"id" db:"contato_id"`
	Nome         string `json:"nome" db:"nome_contato"`
	Email        string `json:"email" db:"email1"`
	UltimoAcesso string `json:"ultimo_acesso" db:"logacr_when"`
	Acessos      int    `json:"acessos" db:"acessos"`
}

//AccessTokenData stores information about access token session
type AccessTokenData struct {
	ValidUntil                 int
	ContatoID                  int
	RoleLevel                  string
	DateWhenSecretWasGenerated time.Time
	IPUsedToGenerateSecret     string
	IPUsedToGenerateAC         string
}

//AccessTokenPublic access token information to be returned to the user
type AccessTokenPublic struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

//User represents user associated to ClientID and Secret
type User struct {
	ID                int    `json:"user_id" db:"logcli_id"`
	LegacyOrPartnerID string `json:"legacy_partner_id" db:"logcli_clientlegacyid"`
	Name              string `json:"name,omitempty" db:"logcli_clientname"`
	ClientID          string `json:"clientID,omitempty" form:"formEmail" db:"logcli_clientid"`
	Secret            string `json:"secret,omitempty" form:"formSenha" db:"logcli_secret"`
	Role              string `db:"logcli_role"`
	LastUpdate        string `db:"logcli_lastupdate"`
}

// Claims dados a serem recuperados da chave
type Claims struct {
	IP string `json:"ip"`
	jwt.StandardClaims
}

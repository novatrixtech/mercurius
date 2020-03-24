package auth

import "time"

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

//OAuthUser represents user associated to ClientID and Secret
type OAuthUser struct {
	UserID   int    `json:"user_id"`
	Name     string `json:"name,omitempty"`
	ClientID string `json:"clientID,omitempty"`
	Secret   string `json:"secret,omitempty"`
}

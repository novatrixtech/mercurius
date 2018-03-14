package model

import "time"

//AccessTokenData stores information about access token session
type AccessTokenData struct {
	ValidUntil                 int
	ContatoID                  int
	RoleLevel                  int
	DateWhenSecretWasGenerated time.Time
	IPUsedToGenerateSecret     string
	IPUsedToGenerateAC         string
}

//AccessTokenPublic access token information to be returned to the user
type AccessTokenPublic struct {
	AccessToken string `json:"Access-Token"`
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

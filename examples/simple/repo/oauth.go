package repo

import (
	"database/sql"
	"log"
	"strings"

	"github.com/novatrixtech/mercurius/examples/simple/conf"
	"github.com/novatrixtech/mercurius/examples/simple/model"
)

//AddAccessTokenAccessLog adds log when an user call the API using an Access Token
func AddAccessTokenAccessLog(accessToken string, funcName string) (err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[AddAccessTokenAccessLog] Erro Pegando o DB: " + err.Error())
		return
	}
	sql := "INSERT INTO logac_accesstokenacessos (logac_accesstoken, logac_function, logac_when) VALUES (?,?,NOW())"
	_, err = db.Exec(sql, accessToken, funcName)
	return
}

//AddAccessTokenRequestLog logs a new access token generation
func AddAccessTokenRequestLog(accessToken string, contatoID int) (err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[AddAccessTokenRequestLog] Erro Pegando o DB: " + err.Error())
		return
	}
	sql := "INSERT INTO logacr_logaccesstokenrequest (logacr_accesstoken, sys_id, logacr_when) VALUES (?,?,NOW())"
	_, err = db.Exec(sql, accessToken, contatoID)
	if err != nil {
		log.Println("[AddAccessTokenRequestLog] Erro ao executar a query: ", sql, " com os parametros: [", accessToken, "] e [", contatoID, "] - Erro: ", err.Error())
	}
	return
}

//GetUserRoleByContactID get active user's role based on contato_id (user ID)
func GetUserRoleByContactID(contactID int) (role int, err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[GetUserRoleByContactID] Erro Pegando o DB: " + err.Error())
		return
	}
	sql := `select 
				a.rol_id
			from
				rol_roles a,
				sys_systems b,
				rsys_rolesystems c
			where
				a.rol_id = c.rol_id and
				b.sys_id = c.sys_id and 
				b.sys_active = 1 and
				b.sys_id = ?`
	err = db.Get(&role, sql, contactID)
	if err != nil {
		log.Println("[GetUserRoleByContactID] Erro ao executar a query: ", sql, " com o parametro: ", contactID, " - Erro: ", err.Error())
		return
	}
	return
}

//GetUserNameByContactID get active user's name based on contato_id (user ID)
func GetUserNameByContactID(contactID int) (name string, err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[GetUserNameByContactID] Erro Pegando o DB: " + err.Error())
		return
	}
	sql := `select 
				sys_name
			FROM 
				sys_systems  
			where 
				sys_id=?`
	err = db.Get(&name, sql, contactID)
	if err != nil {
		log.Println("[GetUserNameByContactID] Erro ao executar a query: ", sql, " com o parametro: ", contactID, " - Erro: ", err.Error())
		return
	}
	return
}

//AddCredentialsToUser insert user's record adding her clientID and secret
func AddCredentialsToUser(user model.OAuthUser, clientID string, secret string) (err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[AddCredentialsToUser] Erro Pegando o DB: " + err.Error())
		return
	}
	sql := "INSERT INTO sys_systems (sys_name, sys_active, sys_clientid, sys_secret) VALUES (?, 1, ?, ?)"
	res, err := db.Exec(sql, user.Name, clientID, secret)
	if err != nil {
		log.Println("[AddCredentialsToUser] Erro ao executar a query: ", sql, " com os parametros: [", user.UserID, "] e [", clientID, "] e [", secret, "] - Erro: ", err.Error())
		return
	}
	sysID, err := res.LastInsertId()
	if err != nil {
		log.Println("[AddCredentialsToUser] Erro ao executar ao obter o ID do novo registro de credencial: ", sql, " com os parametros: [", user.UserID, "] e [", clientID, "] e [", secret, "] - Erro: ", err.Error())
		return
	}
	sql = "INSERT INTO rsys_rolesystems (rol_id, sys_id) VALUES (1, ?)"
	res, err = db.Exec(sql, int(sysID))
	if err != nil {
		log.Println("[AddCredentialsToUser] Erro ao executar a query: ", sql, " com os parametros: [", user.UserID, "] e [", clientID, "] e [", secret, "] - Erro: ", err.Error())
		return
	}
	//log.Println("[AddCredentialsToUser] executado a query: ", sql, " com os parametros: [", user.UserID, "] e [", clientID, "] e [", secret, "] ")
	return
}

//UpdateUserCredentials updates user's record adding her clientID and secret
func UpdateUserCredentials(user model.OAuthUser, clientID string, secret string) (err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[UpdateUserCredentials] Erro Pegando o DB: " + err.Error())
		return
	}
	sql := "UPDATE sys_systems SET sys_clientid = ?, sys_secret = ? WHERE sys_id = ?"
	_, err = db.Exec(sql, clientID, secret, user.UserID)
	if err != nil {
		log.Println("[UpdateUserCredentials] Erro ao executar a query: ", sql, " com os parametros: [", clientID, "] e [", secret, "] e [", user.UserID, "] - Erro: ", err.Error())
		return
	}
	//log.Println("[UpdateUserCredentials] Executado a query: ", sql, " com os parametros: [", user.UserID, "] e [", clientID, "] e [", secret, "]")
	return
}

//GetUserCredentials gets user credentials
func GetUserCredentials(user model.OAuthUser) (clientID string, secret string, err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[StatusUserCredentials] Erro Pegando o DB: " + err.Error())
		return
	}
	strSQL := "SELECT coalesce(sys_clientid, '') as clientID, coalesce(sys_secret, '') as secret FROM sys_systems WHERE sys_id = ?"
	row := db.QueryRow(strSQL, user.UserID)
	err = row.Scan(&clientID, &secret)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		log.Println("[GetUserCredentials] Erro ao executar a query: ", strSQL, " com os parametros: [", user.UserID, "]  - Erro: ", err.Error())
		return
	}
	return
}

//StatusUserCredentials checks if user's has role set and if she has whether credentials is defined or not
func StatusUserCredentials(user model.OAuthUser) (hasRole bool, hasCredentials bool, err error) {
	clientID, secret, err := GetUserCredentials(user)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Println("[StatusUserCredentials] Erro ao pegar o clientID e Secret: " + err.Error())
		return
	}
	hasRole = true
	if len(strings.TrimSpace(clientID)) > 5 && len(strings.TrimSpace(secret)) > 5 {
		hasCredentials = true
	}
	return
}

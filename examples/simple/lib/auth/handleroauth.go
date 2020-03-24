package auth

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/novatrixtech/cryptonx"

	"github.com/novatrixtech/mercurius/examples/simple/conf"
	"github.com/novatrixtech/mercurius/examples/simple/lib/contx"
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
	user := OAuthUser{}
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
	user := OAuthUser{}
	user.UserID = ctx.ParamsInt(":idclient")
	user.Name, err = GetUserNameByContactID(user.UserID)
	if err != nil {
		log.Printf("[GetUserCredentials] Erro ao buscar o nome do usuario de ID [%d]. Erro: %s\n", user.UserID, err.Error())
		ctx.JSON(http.StatusBadRequest, "{'error':'Invalid user ID'}")
		return
	}
	clientID, secret, err := GetUserCredentials(user)
	if err != nil {
		log.Printf("[GetUserCredentials] Erro ao obter as credenciais do usuario [%d] no banco. Erro: %s\n", user.UserID, err.Error())
		ctx.JSON(http.StatusInternalServerError, "{'error': 'Could not obtain credentials'}")
		return
	}
	user.ClientID = clientID
	user.Secret = secret
	ctx.JSON(http.StatusCreated, user)
	return
}

func decodeClientID(origem string) (texto string, err error) {
	err = nil
	tmp, err := hex.DecodeString(origem)
	if err != nil {
		log.Println("[decodeClientID] Erro ao decodar o clientID: ", origem, " - Erro: ", err.Error())
		return
	}
	texto = string(tmp)
	return
}

func getDataFromClientID(clientIDDecoded string) (contactName string, nonce string, err error) {
	err = nil
	if !strings.Contains(clientIDDecoded, "|") {
		err = errors.New("ClientID inválido. Não há o pipe, portanto não há como obter o nonce")
		return
	}
	tmpClientID := strings.Split(clientIDDecoded, "|")
	contactName = tmpClientID[0]
	nonce = tmpClientID[1]
	return
}

func decodeSecret(origem string, nonce string) (texto string, err error) {
	err = nil
	texto, err = cryptonx.Decrypter(conf.Cfg.Section("").Key("oauth_key").Value(), nonce, origem)
	if err != nil {
		log.Println("[decodeSecret] Erro ao decodar o secret: ", origem, " - Erro: ", err.Error())
		return
	}
	return
}

func getAndValidateDataFromSecret(secret string) (data time.Time, contatoID int, IP string, err error) {
	err = nil
	if !strings.Contains(secret, "|") {
		err = errors.New("Secret inválido. Não há o pipe, portanto não há como obter o nonce")
		return
	}
	tmp := strings.Split(secret, "|")
	if len(tmp) < 3 {
		err = errors.New("Secret inválido. Não há elementos suficientes nos dados")
		return
	}
	dataNum, err := strconv.ParseInt(tmp[0], 10, 64)
	if err != nil {
		log.Println("[getInfoFromSecret] Erro ao fazer o parse do timestamp: ", tmp[0], " - Erro: ", err.Error())
		return
	}
	if dataNum < 1505740412 {
		err = errors.New("Secret inválido. Data definida é menor que 2017-09-17")
		return
	}
	data, err = parseDateFromUnixTimestamp(tmp[0])
	if err != nil {
		log.Println("[getInfoFromSecret] Erro ao fazer o parse da data: ", tmp[0], " - Erro: ", err.Error())
		return
	}

	contatoID, err = strconv.Atoi(tmp[1])
	if err != nil {
		log.Println("[getInfoFromSecret] Erro ao fazer o parse do contatoID: ", tmp[1], " - Erro: ", err.Error())
		return
	}
	if contatoID < 1 {
		err = errors.New("ContatoID inválido")
		return
	}
	if len(tmp[2]) < 3 {
		err = errors.New("IP inválido. Numero de itens insuficientes")
		log.Println("[getInfoFromSecret] ", tmp[2], " - Erro: ", err.Error())
		return
	}
	IP = tmp[2]
	return
}

func parseDateFromUnixTimestamp(origem string) (data time.Time, err error) {
	err = nil
	i, err := strconv.ParseInt(origem, 10, 64)
	if err != nil {
		log.Println("[parseDateFromUnixTimestamp] Erro ao fazer o parse do timestamp: ", origem, " - Erro: ", err.Error())
		return
	}
	data = time.Unix(i, 0)
	return
}

func decodeClientIDAndSecret(clientID string, secret string) {
	clientIDDecoded, err := decodeClientID(clientID)
	if err != nil {
		log.Println("[GenerateCredentials] Erro ao decodar o clientID. Erro: ", err.Error())
		return
	}
	log.Println("clientIDDecodado: ", clientIDDecoded)
	_, nonce, err := getDataFromClientID(clientIDDecoded)
	if err != nil {
		log.Println("[GenerateCredentials] Erro ao obter nonce do clientID. Erro: ", err.Error())
		return
	}
	secretDecoded, err := decodeSecret(secret, nonce)
	if err != nil {
		log.Println("[GenerateCredentials] Erro ao decodar o secret. Erro: ", err.Error())
		return
	}
	log.Println("SecretDecodado: ", secretDecoded)
	dataDoSecret, contatoID, IPDoSecret, err := getAndValidateDataFromSecret(secretDecoded)
	if err != nil {
		log.Println("[GenerateCredentials] Erro ao obter os dados do secret. Erro: ", err.Error())
		return
	}
	log.Println("Data: ", dataDoSecret, " - ContatoID: ", contatoID, " - IP: ", IPDoSecret)
}

func generateUserCredentials(user OAuthUser, remoteAddr string) (clientID string, secret string, err error) {
	err = nil
	nomeContatoOrigem := strings.Replace(user.Name, " ", "", -1)
	dataOrigem := time.Now().Unix()
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
		log.Printf("[generateUserCredentials] Erro ao split o host e a porta. userip: %q is not IP:port", remoteAddr)
	}
	ipRemotoOrigem := net.ParseIP(ip)
	if ipRemotoOrigem == nil {
		errStr := fmt.Sprintf("[generateUserCredentials] Erro ao fazer o parse do userip: %q is not IP:port", ip)
		log.Println(errStr)
		err = errors.New(errStr)
		return
	}
	//log.Println("Dados: ", nomeContatoOrigem, "|", contatoIDOrigem, "|", dataOrigem, "|", ipRemotoOrigem)
	secretAntesCrypto := strconv.Itoa(int(dataOrigem)) + "|" + strconv.Itoa(user.UserID) + "|" + ipRemotoOrigem.String()
	//log.Println("[generateUserCredentials] secretAntesCrypto: ", secretAntesCrypto)
	secret, nonce, err := cryptonx.Encrypter(conf.Cfg.Section("").Key("oauth_key").Value(), secretAntesCrypto)
	if err != nil {
		log.Println("[GenerateCredentials] Erro ao encriptar texto: ", err.Error())
		return
	}
	//log.Println("[generateUserCredentials] secret: ", secret, " - nonce: ", nonce)
	clientIDOrigem := nomeContatoOrigem + "|" + nonce
	//log.Println("[generateUserCredentials] clientIDOrigem: ", clientIDOrigem)
	clientID = hex.EncodeToString([]byte(clientIDOrigem))
	//log.Println("[generateUserCredentials] clientID: ", clientID)
	return
}

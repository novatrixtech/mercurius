package conf

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

var mgoSession *mgo.Session
var mongoDbURI string
var mongoDbName string

/*
LoadMongoConfig Init MongoDB Connection
*/
func LoadMongoConfig() {
	mongoDbURI = Cfg.Section("").Key("mongo_uri").Value()
	mongoDbName = Cfg.Section("").Key("mongo_db").Value()		
}

/*
GetMongoSession gets connection to Mongo repo
*/
func GetMongoSession() (*mgo.Session, error) {
	if mgoSession != nil {
		mgoSession.Refresh()
		return mgoSession.Copy(), nil
	}

	LoadMongoConfig()
	mgoSession, erro := mgo.Dial(mongoDbURI)
	if erro != nil {
		log.Printf("[GetMongoSession] Erro ao tentar abrir a sessao com o Mongo: [%s]\n", erro.Error())
		return nil, erro
	}
	return mgoSession.Copy(), erro
}

/*
GetMongoCollection gets a data collection
*/
func GetMongoCollection(collectionName string) (*mgo.Collection, error) {
	mgoSession, erro := GetMongoSession()
	if erro != nil {
		log.Printf("[GetCollection] Erro ao conectar ao Mongo: [%s]\n", erro.Error())
		return nil, erro
	}

	collection := mgoSession.DB(mongoDbName).C(collectionName)
	return collection, nil
}

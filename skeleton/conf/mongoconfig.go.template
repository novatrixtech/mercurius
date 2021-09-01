package conf

import (
	"log"
	"os"
	mgo "gopkg.in/mgo.v2"
)

var mgoSession *mgo.Session
var mongoDbURI string
var mongoDbName string

/*
LoadMongoConfig Init MongoDB Connection
*/
func LoadMongoConfig() {
	if len(os.Getenv("mongo_uri")) < 1 {
		mongoDbURI = Cfg.Section("").Key("mongo_uri").Value()
		mongoDbName = Cfg.Section("").Key("mongo_db").Value()
		return
	}
	mongoDbURI = os.Getenv("mongo_uri")
	mongoDbName = os.Getenv("mongo_db")
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
	var err error
	mgoSession, err = mgo.Dial(mongoDbURI)
	if err != nil {
		log.Printf("[GetMongoSession] Error opening mongo db session: [%s]\n", err.Error())
		return nil, err
	}
	return mgoSession.Copy(), err
}

/*
GetMongoCollection gets a data collection
*/
func GetMongoCollection(collectionName string) (*mgo.Collection, error) {
	mgoSession, err := GetMongoSession()
	if err != nil {
		log.Printf("[GetCollection] Error connecting to mongo db: [%s]\n", err.Error())
		return nil, err
	}

	collection := mgoSession.DB(mongoDbName).C(collectionName)
	return collection, nil
}

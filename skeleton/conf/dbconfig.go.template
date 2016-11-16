package conf

import (
	"fmt"
	"strconv"

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	// import postgre driver
	_ "github.com/lib/pq"
)

/*
DbConnInfo - Database Connection information
*/
type DbConnInfo struct {
	DBType string
	User   string
	Pw     string
	DBName string
	Host   string
	Port   string
}

/*
Database - Database connection and configuration
*/
type Database interface {
	DSN() string
	GetDB() (*sqlx.DB, error)
}

/*
DB - Represents SQLX database connection
*/
var DB *sqlx.DB

/*
DBConnData - Value of database connection information
*/
var DBConnData *DbConnInfo

/*
GetDB - Gets Database connection
*/
func GetDB() (*sqlx.DB, error) {
	var err error
	if DB == nil {
		if DBConnData == nil {
			DBConnData = LoadDbConfig()
		}
		dsn := DSN()
		DB, err = sqlx.Open(DBConnData.DBType, dsn)
		if err != nil {
			return nil, err
		}
		err = DB.Ping()
		if err != nil {
			return nil, err
		}
		idle := Cfg.Section("").Key("idle_conn").Value()
		numIdle, err := strconv.Atoi(idle)
		if err != nil {
			return nil, err
		}
		DB.SetMaxIdleConns(numIdle)
		max := Cfg.Section("").Key("max_conn").Value()
		numMax, err := strconv.Atoi(max)
		if err != nil {
			return nil, err
		}
		DB.SetMaxOpenConns(numMax)
	}
	return DB, err
}

/*
DSN - Gets database connection string
*/
func DSN() string {
	if DBConnData.DBType == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", DBConnData.User, DBConnData.Pw, DBConnData.Host, DBConnData.Port, DBConnData.DBName)
	} else if DBConnData.DBType == "postgresql" {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DBConnData.User, DBConnData.Pw, DBConnData.Host, DBConnData.Port, DBConnData.DBName)
	} else {
		return ""
	}
}

/*
LoadDbConfig - Loads specific database connection information
*/
func LoadDbConfig() *DbConnInfo {
	return &DbConnInfo{
		DBType: Cfg.Section("").Key("db_type").Value(),
		User:   Cfg.Section("").Key("db_user").Value(),
		Pw:     Cfg.Section("").Key("db_pw").Value(),
		DBName: Cfg.Section("").Key("db_name").Value(),
		Host:   Cfg.Section("").Key("db_host").Value(),
		Port:   Cfg.Section("").Key("db_port").Value(),
	}
}
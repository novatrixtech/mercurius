package conf

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	// import MySQL driver
	_ "github.com/go-sql-driver/mysql"
	// import PostgreSQL driver
	_ "github.com/lib/pq"
	// import SQL Server driver
	_ "github.com/denisenkom/go-mssqldb"
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
		var numMax, numIdle int
		if len(os.Getenv("db_type")) < 1 {
			numIdle, err = strconv.Atoi(Cfg.Section("").Key("idle_conn").Value())
			if err != nil {
				return nil, err
			}
			numMax, err = strconv.Atoi(Cfg.Section("").Key("max_conn").Value())
			if err != nil {
				return nil, err
			}
		} else {
			numIdle, err = strconv.Atoi(os.Getenv("idle_conn"))
			if err != nil {
				return nil, err
			}
			numMax, err = strconv.Atoi(os.Getenv("max_conn"))
			if err != nil {
				return nil, err
			}
		}
		DB.SetMaxIdleConns(numIdle)
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
	} else if DBConnData.DBType == "postgres" {
		return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", DBConnData.User, DBConnData.Pw, DBConnData.Host, DBConnData.Port, DBConnData.DBName)
	} else if DBConnData.DBType == "mssql" {
		return fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s", DBConnData.Host, DBConnData.User, DBConnData.Pw, DBConnData.Port, DBConnData.DBName)
	}
	return ""
}

/*
LoadDbConfig - Loads specific database connection information
*/
func LoadDbConfig() (dbConnectionInfo *DbConnInfo) {
	if len(os.Getenv("db_type")) < 1 {
		dbConnectionInfo = &DbConnInfo{
			DBType: Cfg.Section("").Key("db_type").Value(),
			User:   Cfg.Section("").Key("db_user").Value(),
			Pw:     Cfg.Section("").Key("db_pw").Value(),
			DBName: Cfg.Section("").Key("db_name").Value(),
			Host:   Cfg.Section("").Key("db_host").Value(),
			Port:   Cfg.Section("").Key("db_port").Value(),
		}
	} else {
		dbConnectionInfo = &DbConnInfo{
			DBType: os.Getenv("db_type"),
			User:   os.Getenv("db_user"),
			Pw:     os.Getenv("db_pw"),
			DBName: os.Getenv("db_name"),
			Host:   os.Getenv("db_host"),
			Port:   os.Getenv("db_port"),
		}
	}
	return
}

//checks if db config data is coming from environment variables
func isDbConnParamsInEnvVariables() bool {
	if len(os.Getenv("db_type")) > 1 {
		return true
	}
	return false
}
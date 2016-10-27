package conf

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strconv"
)

var DB *sqlx.DB

type MySQL struct {
	DBType string
	User   string
	Pw     string
	DBName string
	Host   string
	Port   string
}

func (db *MySQL) DB() (*sqlx.DB, error) {
	var err error
	if DB == nil {
		dsn := db.DSN()
		DB, err = sqlx.Open(db.DBType, dsn)
		if err != nil {
			return nil, err
		}
		err = DB.Ping()
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

func (db *MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", db.User, db.Pw, db.Host, db.Port, db.DBName)
}

func LoadMySQLConfig() *MySQL {
	return &MySQL{
		DBType: Cfg.Section("").Key("db_type").Value(),
		User:   Cfg.Section("").Key("db_user").Value(),
		Pw:     Cfg.Section("").Key("db_pw").Value(),
		DBName: Cfg.Section("").Key("db_name").Value(),
		Host:   Cfg.Section("").Key("db_host").Value(),
		Port:   Cfg.Section("").Key("db_port").Value(),
	}
}

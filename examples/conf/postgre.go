package conf

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type PostgreSQL struct {
	DBType string
	User   string
	Pw     string
	DBName string
	Host   string
	Port   string
}

func (db *PostgreSQL) DB() (*sqlx.DB, error) {
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
		max := Cfg.Section("").Key("idle_conn").Value()
		numMax, err := strconv.Atoi(max)
		if err != nil {
			return nil, err
		}
		DB.SetMaxOpenConns(numMax)
	}
	return DB, err
}

func (db *PostgreSQL) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db.User, db.Pw, db.Host, db.Port, db.DBName)
}

func LoadPostgreSQLConfig() *PostgreSQL {
	return &PostgreSQL{
		DBType: Cfg.Section("").Key("db_type").Value(),
		User:   Cfg.Section("").Key("db_user").Value(),
		Pw:     Cfg.Section("").Key("db_pw").Value(),
		DBName: Cfg.Section("").Key("db_name").Value(),
		Host:   Cfg.Section("").Key("db_host").Value(),
		Port:   Cfg.Section("").Key("db_port").Value(),
	}
}

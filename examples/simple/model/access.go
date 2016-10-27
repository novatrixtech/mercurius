package model

import (
	"database/sql"
	"fmt"
	"github.com/felipeweb/gopher-utils"
	"github.com/go-sql-driver/mysql"
)

type Access struct {
	ID           int            `db:"id" json:"id"`
	AdvDefaultID sql.NullInt64  `db:"advertising_default_id" json:"-"`
	AdvByCarID   sql.NullInt64  `db:"advertising_by_car_id" json:"-"`
	Marca        string         `db:"marca" json:"marca"`
	Modelo       string         `db:"modelo" json:"modelo"`
	Ano          int            `db:"ano" json:"ano"`
	Combustivel  string         `db:"combustivel" json:"combustivel"`
	Propriedade  string         `db:"propriedade" json:"propriedade"`
	Latitude     float64        `db:"latitude" json:"latitude"`
	Longitude    float64        `db:"longitude" json:"longitude"`
	Date         mysql.NullTime `db:"date" json:"date"`
	Time         sql.NullString `db:"time" json:"time"`
}

func (ac *Access) FormattedDate() string {
	date := gopher_utils.DateT(ac.Date.Time, "DD/MM/YYYY")
	return fmt.Sprintf("%s %s", date, ac.Time.String)
}

package repository

import (
	"fmt"
	"strconv"
	"time"

	"github.com/novatrixtech/mercurius/examples/simple/conf"

	"github.com/go-macaron/cache"
	"github.com/jmoiron/sqlx"
	"github.com/novatrixtech/mercurius/examples/simple/model"
)

type AccessRepository struct {
	db *sqlx.DB
}

func NewAccessRepository() (*AccessRepository, error) {
	db, err := conf.GetDB()
	if err != nil {
		return nil, err
	}
	return &AccessRepository{db}, nil
}

func (repo *AccessRepository) FindAllBy(query string, cache cache.Cache) ([]model.Access, error) {
	access := []model.Access{}
	if cache.IsExist(query) {
		access = cache.Get(query).([]model.Access)
		return access, nil
	}
	err := repo.db.Select(&access, query)
	s := fmt.Sprintf("%.0f", (time.Hour * 4).Seconds())
	timeout, _ := strconv.Atoi(s)
	cache.Put(query, access, int64(timeout))
	return access, err
}

package repository

import (
	"pheet-fiber-backend/service/appinfo"

	"github.com/jmoiron/sqlx"
)

type appInfoRepository struct {
	DB *sqlx.DB
}

func NewAppInfoRepository(DB *sqlx.DB) appinfo.AppInfoRepository {
	return appInfoRepository{DB: DB}
}
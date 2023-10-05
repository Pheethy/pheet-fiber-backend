package usecase

import (
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/service/appinfo"
)

type appInfoUsecase struct {
	cfg      config.Iconfig
	infoRepo appinfo.AppInfoRepository
}

func NewAppInfoUsecase(cfg config.Iconfig, infoRepo appinfo.AppInfoRepository) appinfo.AppInfoUsecase {
	return appInfoUsecase{
		cfg: cfg,
		infoRepo: infoRepo,
	}
}

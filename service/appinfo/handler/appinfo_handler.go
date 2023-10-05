package handler

import (
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/service/appinfo"
)

type appInfoHandler struct {
	cfg config.Iconfig
	infoUs appinfo.AppInfoUsecase
}

func NewAppInfoHandler(cfg config.Iconfig, infoUs appinfo.AppInfoUsecase) appinfo.AppInfoHandler {
	return appInfoHandler{
		cfg: cfg,
		infoUs: infoUs,
	}
}
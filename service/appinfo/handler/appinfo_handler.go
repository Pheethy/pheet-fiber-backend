package handler

import (
	"net/http"
	"pheet-fiber-backend/auth/service"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/service/appinfo"

	"github.com/gofiber/fiber/v2"
)

type appInfoHandler struct {
	cfg    config.Iconfig
	infoUs appinfo.AppInfoUsecase
}

func NewAppInfoHandler(cfg config.Iconfig, infoUs appinfo.AppInfoUsecase) appinfo.AppInfoHandler {
	return appInfoHandler{
		cfg:    cfg,
		infoUs: infoUs,
	}
}

func (a appInfoHandler) GenerateAPIKey(c *fiber.Ctx) error {
	apiKey, err := service.NewAuthService(
		constants.APIKey,
		a.cfg.Jwt(),
		nil,
	)

	if err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, err.Error())
	}

	resp := map[string]interface{}{
		"key": apiKey.SignToken(),
	}

	return c.Status(http.StatusOK).JSON(resp)
}

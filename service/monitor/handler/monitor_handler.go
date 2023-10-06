package handler

import (
	"net/http"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/monitor"

	"github.com/gofiber/fiber/v2"
)

type monitorHandler struct {
	cfg config.Iconfig
}

func NewMonitorHandler(cfg config.Iconfig) monitor.IMonitorHandler {
	return &monitorHandler{
		cfg: cfg,
	}
}

func (m monitorHandler) HealthCheck(c *fiber.Ctx) error {
	resp := models.Monitor{
		Name:    m.cfg.App().Name(),
		Version: m.cfg.App().Version(),
	}

	return c.Status(http.StatusOK).JSON(resp)
}

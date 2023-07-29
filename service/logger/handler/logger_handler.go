package handler

import (
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/logger"
	"time"

	"github.com/gofiber/fiber/v2"
)

func InitLogger(c *fiber.Ctx, resp any) logger.ILogger {
	log := &models.Logger{
		Time: time.Now().Local().Format("2006-01-02 18:00:00"),
		Ip: c.IP(),
		Method: c.Method(),
		StatusCode: c.Response().StatusCode(),
	}

	log.SetQuery(c)
	log.SetBody(c)
	log.SetResp(resp)

	return log
}
package handler

import (
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/middleware"
	"pheet-fiber-backend/middleware/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewareHandler struct {
	cfg      config.Iconfig
	middleUs usecase.ImiddlewareUsecase
}

func NewMiddlewareHandler(cfg config.Iconfig, middleUs usecase.ImiddlewareUsecase) middleware.ImiddlewareHandler {
	return middlewareHandler{
		cfg:      cfg,
		middleUs: middleUs,
	}
}

func (m middlewareHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, PUT, PATCH, HEAD, DELETE",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (m middlewareHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format: "👽 ${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone: "Bangkok/Asia",
	})
}

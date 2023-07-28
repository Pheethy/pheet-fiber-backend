package handler

import (
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/middleware/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type ImiddlewareHandler interface {
	Cors() fiber.Handler
}

type middlewareHandler struct {
	cfg      config.Iconfig
	middleUs usecase.ImiddlewareUsecase
}

func NewMiddlewareHandler(cfg config.Iconfig, middleUs usecase.ImiddlewareUsecase) ImiddlewareHandler {
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

package handler

import (
	"net/http"
	_auth_service "pheet-fiber-backend/auth/service"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/middleware"
	"pheet-fiber-backend/middleware/usecase"
	"strings"

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
		Format:     "👽 ${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

func (m middlewareHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := _auth_service.ParseToken(m.cfg.Jwt(), token)
		if err != nil {
			return fiber.NewError(http.StatusUnauthorized, err.Error())
		}

		claims := result.Claims

		if !m.middleUs.FindAccessToken(claims.Id, token) {
			return fiber.NewError(http.StatusUnauthorized, "no permission to access")
		}

		c.Locals("user_id", claims.Id)
		c.Locals("role_id", claims.RoleId)
		return c.Next()
	}
}

func (m middlewareHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("user_id")
		if c.Params("user_id") != userId {
			return fiber.NewError(http.StatusBadRequest, "never gonna give you up")
		}
		return c.Next()
	}
}

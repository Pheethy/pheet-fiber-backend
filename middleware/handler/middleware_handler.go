package handler

import (
	"net/http"
	_auth_service "pheet-fiber-backend/auth/service"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/middleware"
	"pheet-fiber-backend/middleware/usecase"
	"pheet-fiber-backend/service/utils"
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

func (m middlewareHandler) Authorize(expectedRoleId ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals("role_id").(int)
		if !ok {
			return fiber.NewError(http.StatusUnprocessableEntity, "cast role_id to int failed.")
		}

		roles, err := m.middleUs.FindRole()
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}


		sum := 0
		for _, val := range expectedRoleId {
			sum += val
		}

		expectedValBinary := utils.ConvertBinary(sum , len(roles))
		userValBinary := utils.ConvertBinary(userRoleId, len(roles))
		


		for index := range userValBinary {
			if userValBinary[index]&expectedValBinary[index] == 1 {
				return c.Next()
			}
		}
		

		return fiber.NewError(http.StatusUnauthorized)
	}
}

func (m middlewareHandler) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var key = c.Get("X-API-KEY")
		if _, err := _auth_service.ParseApiKey(m.cfg.Jwt(), key); err != nil {
			return fiber.NewError(http.StatusInternalServerError, "API-KEY is invalid.")
		}
		return c.Next()
	}
}
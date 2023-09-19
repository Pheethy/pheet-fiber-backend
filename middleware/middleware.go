package middleware

import "github.com/gofiber/fiber/v2"

type ImiddlewareHandler interface {
	Cors() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
}
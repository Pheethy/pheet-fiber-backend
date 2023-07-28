package logger

import "github.com/gofiber/fiber/v2"

type ILogger interface {
	Print() ILogger
	Save()
	SetQuery(c *fiber.Ctx)
	SetBody(c *fiber.Ctx)
	SetResp(resp any)
}
package appinfo

import "github.com/gofiber/fiber/v2"

type AppInfoHandler interface {
	GenerateAPIKey(c *fiber.Ctx) error
}
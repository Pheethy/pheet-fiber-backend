package appinfo

import "github.com/gofiber/fiber/v2"

type AppInfoHandler interface {
	GenerateAPIKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
	AddCategory(c *fiber.Ctx) error
	RemoveCategory(c *fiber.Ctx) error
}
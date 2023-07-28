package monitor

import "github.com/gofiber/fiber/v2"

type IMonitorHandler interface {
	HealthCheck(c *fiber.Ctx) error
}
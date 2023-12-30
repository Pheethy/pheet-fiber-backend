package order

import "github.com/gofiber/fiber/v2"

type IOrderHandler interface {
	FetchAllOrder(c *fiber.Ctx) error
	FetchOneOrder(c *fiber.Ctx) error
}
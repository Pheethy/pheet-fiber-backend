package products

import "github.com/gofiber/fiber/v2"

/* Handlers Port */
type IProductsHandlers interface {
	FetchAllProduct(c *fiber.Ctx) error
	FetchOneProduct(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

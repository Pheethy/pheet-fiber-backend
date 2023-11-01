package product

import "github.com/gofiber/fiber/v2"

type IProductHandler interface {
	FetchOneProduct(c *fiber.Ctx) error
	FetchAllProduct(c *fiber.Ctx) error
	CreateProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}
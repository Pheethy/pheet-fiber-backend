package product

import "github.com/gofiber/fiber/v2"

type IProductHandler interface {
	FetchOneProduct(c *fiber.Ctx) error
}
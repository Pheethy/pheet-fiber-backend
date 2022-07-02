package handler

import "github.com/gofiber/fiber/v2"

type ProductHandler interface{
	GetProducts(c *fiber.Ctx)error
	GetProductById(c *fiber.Ctx)error
	GetProductByType(c *fiber.Ctx)error
	CreateProduct(c *fiber.Ctx)error
	UpdateProduct(c *fiber.Ctx)error
	DeleteProduct(c *fiber.Ctx)error
}
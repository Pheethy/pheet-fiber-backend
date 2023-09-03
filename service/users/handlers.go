package users

import "github.com/gofiber/fiber/v2"

type IUsersHandlers interface {
	SignUpCustomer(c *fiber.Ctx) error
	GetPassport(c *fiber.Ctx) error
}
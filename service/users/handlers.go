package users

import "github.com/gofiber/fiber/v2"

type IUsersHandlers interface {
	SignUpCustomer(c *fiber.Ctx) error
	GetPassport(c *fiber.Ctx) error
	FetchUserProfile(c *fiber.Ctx) error
	GenerateAdminToken(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
}
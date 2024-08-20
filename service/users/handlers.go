package users

import "github.com/gofiber/fiber/v2"

type IUsersHandlers interface {
	FetchUserProfile(c *fiber.Ctx) error
	InsertUser(c *fiber.Ctx) error
	InsertAdmin(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
	GenerateAdminToken(c *fiber.Ctx) error
}


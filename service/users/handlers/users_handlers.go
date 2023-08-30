package handlers

import (
	"net/http"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/users"

	"github.com/gofiber/fiber/v2"
)

type usersHandlers struct {
	cfg     config.Iconfig
	usersUs users.IUsersUsecase
}

func NewUsersHandler(cfg config.Iconfig, usersUs users.IUsersUsecase) users.IUsersHandlers {
	return usersHandlers{
		cfg:     cfg,
		usersUs: usersUs,
	}
}

func (u usersHandlers) SignUpCustomer(c *fiber.Ctx) error {
	var userReq = new(models.UserRegisterReq)
	if err := c.BodyParser(userReq); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if !userReq.IsEmail() {
		return fiber.NewError(http.StatusBadRequest, "email format is invalid.")
	}

	userPass, err := u.usersUs.InsertCustomer(userReq)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(userPass)
}

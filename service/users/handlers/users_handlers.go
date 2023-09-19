package handlers

import (
	"net/http"
	_auth_service "pheet-fiber-backend/auth/service"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/constants"
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

func (u usersHandlers) SignOut(c *fiber.Ctx) error {
	var ctx = c.Context()
	var userReq = new(models.UserRemoveCredential)
	if err := c.BodyParser(userReq); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if err := u.usersUs.DeleteOauth(ctx, userReq.OauthId); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"message": "sign-out success.",
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (u usersHandlers) GetPassport(c *fiber.Ctx) error {
	var ctx = c.Context()
	var userReq = new(models.UserCredential)
	if err := c.BodyParser(userReq); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	userPass, err := u.usersUs.GetPassport(ctx, userReq)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"user": userPass,
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (u usersHandlers) FetchUserProfile(c *fiber.Ctx) error {
	var ctx = c.Context()
	var id = c.Params("user_id")
	user, err := u.usersUs.FetchUserProfile(ctx, id)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"user": user,
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (u usersHandlers) GenerateAdminToken(c *fiber.Ctx) error {
	adminToken, err := _auth_service.NewAuthService(
		constants.Admin,
		u.cfg.Jwt(),
		nil,
	)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"token": adminToken.SignToken(),
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (u usersHandlers) RefreshPassport(c *fiber.Ctx) error {
	var ctx = c.Context()
	var userReq = new(models.UserRefreshCredential)
	if err := c.BodyParser(userReq); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	userPass, err := u.usersUs.RefreshPassport(ctx, userReq)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"user": userPass,
	}

	return c.Status(http.StatusOK).JSON(resp)
}
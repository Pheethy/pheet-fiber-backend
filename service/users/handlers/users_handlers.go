package handlers

import (
	"pheet-fiber-backend/auth"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/polymor"
	"pheet-fiber-backend/service/users"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type usersHandlers struct {
	cfg     config.Iconfig
	usersUs users.IUsersUsecase
}

func NewUsersHandlers(cfg config.Iconfig, usersUs users.IUsersUsecase) users.IUsersHandlers {
	return usersHandlers{
		cfg:     cfg,
		usersUs: usersUs,
	}
}

func (u usersHandlers) FetchUserProfile(c *fiber.Ctx) error {
	ctx := c.UserContext()
	userId := uuid.FromStringOrNil(c.Params("user_id"))

	userProfile, err := u.usersUs.FetchUserProfile(ctx, &userId)
	if err != nil {
		if ok := strings.Contains(err.Error(), constants.ERROR_USER_NOT_FOUND); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": constants.ERROR_USER_NOT_FOUND,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := map[string]interface{}{
		"user_profile": userProfile,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (u usersHandlers) InsertUser(c *fiber.Ctx) error {
	ctx := c.UserContext()
	userReq := new(models.User)

	if err := c.BodyParser(userReq); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	polymor.SetDefult(userReq)
	if !userReq.IsEmail() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ERROR_EMAIL_PATTERN_IS_INVALID,
		})
	}

	userPass, err := u.usersUs.InsertUser(ctx, userReq)
	if err != nil {
		if ok := strings.Contains(err.Error(), constants.ERROR_USERNAME_HAS_BEEN_USED); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		} else if ok := strings.Contains(err.Error(), constants.ERROR_EMAIL_HAS_BEEN_USED); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := map[string]interface{}{
		"message": "successful",
		"data":    userPass,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (u usersHandlers) InsertAdmin(c *fiber.Ctx) error {
	ctx := c.UserContext()
	userReq := new(models.User)

	if err := c.BodyParser(userReq); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error:": err.Error(),
		})
	}
	polymor.SetDefult(userReq)
	if !userReq.IsEmail() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ERROR_EMAIL_PATTERN_IS_INVALID,
		})
	}

	userPass, err := u.usersUs.InsertAdmin(ctx, userReq)
	if err != nil {
		if ok := strings.Contains(err.Error(), constants.ERROR_USERNAME_HAS_BEEN_USED); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		} else if ok := strings.Contains(err.Error(), constants.ERROR_EMAIL_HAS_BEEN_USED); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := map[string]interface{}{
		"message": "successful",
		"data":    userPass,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (u usersHandlers) GenerateAdminToken(c *fiber.Ctx) error {
	auth, err := auth.NewAuth(constants.Admin, u.cfg.Jwt(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	token := auth.SignToken()
	resp := map[string]interface{}{
		"token": token,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (u usersHandlers) SignIn(c *fiber.Ctx) error {
	ctx := c.UserContext()
	req := new(models.User)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userPass, err := u.usersUs.GetPassport(ctx, req)
	if err != nil {
		if ok := strings.Contains(err.Error(), constants.ERROR_USER_NOT_FOUND); ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error:": constants.ERROR_USERNAME_OR_PASSWORD_IS_INVALID,
			})
		}
		if ok := strings.Contains(err.Error(), constants.ERROR_PASSWORD_NOT_MATCH); ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error:": constants.ERROR_USERNAME_OR_PASSWORD_IS_INVALID,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := map[string]interface{}{
		"message": "successful",
		"data":    userPass,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (u usersHandlers) RefreshPassport(c *fiber.Ctx) error {
	ctx := c.UserContext()
	req := new(models.UserRefreshCredential)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userPass, err := u.usersUs.RefreshPassport(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := map[string]interface{}{
		"message": "successful",
		"data":    userPass,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (u usersHandlers) SignOut(c *fiber.Ctx) error {
	ctx := c.UserContext()
	user := new(models.UserRemoveCredential)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := u.usersUs.DeleteOAuth(ctx, user.OauthId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := map[string]interface{}{
		"message": "successful",
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

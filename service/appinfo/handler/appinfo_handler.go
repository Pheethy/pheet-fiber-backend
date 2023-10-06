package handler

import (
	"net/http"
	"pheet-fiber-backend/auth/service"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/appinfo"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type appInfoHandler struct {
	cfg    config.Iconfig
	infoUs appinfo.AppInfoUsecase
}

func NewAppInfoHandler(cfg config.Iconfig, infoUs appinfo.AppInfoUsecase) appinfo.AppInfoHandler {
	return &appInfoHandler{
		cfg:    cfg,
		infoUs: infoUs,
	}
}

func (h appInfoHandler) GenerateAPIKey(c *fiber.Ctx) error {
	apiKey, err := service.NewAuthService(
		constants.APIKey,
		h.cfg.Jwt(),
		nil,
	)

	if err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, err.Error())
	}

	resp := map[string]interface{}{
		"key": apiKey.SignToken(),
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (h appInfoHandler) FindCategory(c *fiber.Ctx) error {
	var ctx = c.Context()
	var args = new(sync.Map)
	var search = c.Query("search_word")

	if search != "" {
		args.Store("search_word", search)
	}

	cats, err := h.infoUs.FindCategory(ctx, args)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"category": cats,
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (h appInfoHandler) AddCategory(c *fiber.Ctx) error {
	var ctx = c.Context()
	var cats = make([]*models.Catagory, 0)
	if err := c.BodyParser(&cats); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if err := h.infoUs.InsertCategories(ctx, cats); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"message": "created.",
	}
	return c.Status(http.StatusOK).JSON(resp)
}

func (h appInfoHandler) RemoveCategory(c *fiber.Ctx) error {
	var ctx = c.Context()
	var id = c.Params("category_id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, "can't convert string to int")
	}

	if err := h.infoUs.DeleteCategory(ctx, intId); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"message": "deleted.",
	}

	return c.Status(http.StatusOK).JSON(resp)
}

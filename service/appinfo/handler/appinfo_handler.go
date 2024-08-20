package handler

import (
	"net/http"
	"pheet-fiber-backend/auth"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/appinfo"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type appInfoHandler struct {
	cfg       config.Iconfig
	addInfoUs appinfo.AppInfoUsecase
}

func NewAppInfoHandler(cfg config.Iconfig, addInfoUs appinfo.AppInfoUsecase) appinfo.AppInfoHandler {
	return &appInfoHandler{
		cfg:       cfg,
		addInfoUs: addInfoUs,
	}
}

func (a appInfoHandler) GenerateAPIKey(c *fiber.Ctx) error {
	apiKey, err := auth.NewAuth(constants.ApiKey, a.cfg.Jwt(), nil)
	if err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, err.Error())
	}

	resp := map[string]interface{}{
		"key": apiKey.SignToken(),
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (h appInfoHandler) FindCategory(c *fiber.Ctx) error {
	ctx := c.Context()
	args := new(sync.Map)
	search := c.Query("search_word")

	if search != "" {
		args.Store("search_word", search)
	}

	cats, err := h.addInfoUs.FindCategory(ctx, args)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"category": cats,
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (h appInfoHandler) AddCategory(c *fiber.Ctx) error {
	ctx := c.Context()
	cats := make([]*models.Categories, 0)
	if err := c.BodyParser(&cats); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if err := h.addInfoUs.InsertCategories(ctx, cats); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"message": "created.",
	}
	return c.Status(http.StatusOK).JSON(resp)
}

func (h appInfoHandler) RemoveCategory(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("category_id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, "can't convert string to int")
	}

	if err := h.addInfoUs.DeleteCategory(ctx, intId); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"message": "deleted.",
	}

	return c.Status(http.StatusOK).JSON(resp)
}

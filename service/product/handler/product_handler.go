package handler

import (
	"net/http"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/service/file"
	"pheet-fiber-backend/service/product"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type productHandler struct {
	cfg config.Iconfig
	proUs product.IProductUsecase
	fileUs file.IFileUsecase
}

func NewProductHandler(cfg config.Iconfig, proUs product.IProductUsecase, fileUs file.IFileUsecase) product.IProductHandler {
	return productHandler{
		cfg: cfg,
		proUs: proUs,
		fileUs: fileUs,
	}
}

func (h productHandler) FetchOneProduct(c *fiber.Ctx) error {
	var ctx = c.Context()
	var id = strings.TrimSpace(c.Params("product_id"))

	product, err := h.proUs.FetchOneProduct(ctx, id)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"product": product,
	}	
	return c.Status(http.StatusOK).JSON(resp)
}
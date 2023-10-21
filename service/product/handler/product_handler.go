package handler

import (
	"net/http"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/helper"
	"pheet-fiber-backend/service/file"
	"pheet-fiber-backend/service/product"
	"strconv"
	"strings"
	"sync"

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

func (h productHandler) FetchAllProduct(c *fiber.Ctx) error {
	var ctx = c.Context()
	var args = new(sync.Map)
	var paginator = helper.NewPaginator()
	var searchword = c.Query("search_word")
	var page, pageErr = strconv.Atoi(c.Query("page"))
	var perPage, perPageErr = strconv.Atoi(c.Query("per_page"))

	if searchword != "" {
		args.Store("search_word", searchword)
	}

	if pageErr == nil {
		paginator.Page = page
	}

	if perPageErr == nil {
		paginator.PerPage = perPage
	}

	products, err := h.proUs.FetchAllProduct(ctx, args, &paginator)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"products": products,
		"page": paginator.Page,
		"per_page": paginator.PerPage,
		"total_page": paginator.TotalPages,
		"total_rows": paginator.TotalEntrySizes,
	}	
	return c.Status(http.StatusOK).JSON(resp)
}
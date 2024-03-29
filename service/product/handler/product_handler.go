package handler

import (
	"fmt"
	"net/http"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/file"
	"pheet-fiber-backend/service/product"
	"strconv"
	"strings"
	"sync"

	"github.com/Pheethy/psql/helper"
	"github.com/gofiber/fiber/v2"
)

type productHandler struct {
	cfg    config.Iconfig
	proUs  product.IProductUsecase
	fileUs file.IFileUsecase
}

func NewProductHandler(cfg config.Iconfig, proUs product.IProductUsecase, fileUs file.IFileUsecase) product.IProductHandler {
	return productHandler{
		cfg:    cfg,
		proUs:  proUs,
		fileUs: fileUs,
	}
}

func (h productHandler) FetchOneProduct(c *fiber.Ctx) error {
	var ctx = c.Context()
	var id = strings.TrimSpace(c.Params("product_id"))

	product, err := h.proUs.FetchOneProduct(ctx, id)
	if err != nil {
		if ok := strings.Contains(err.Error(), "product not found"); ok {
			return fiber.NewError(http.StatusNoContent, err.Error())
		}
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
	if len(products) < 1 {
		return fiber.NewError(http.StatusNoContent, "no product")
	}

	resp := map[string]interface{}{
		"products":   products,
		"page":       paginator.Page,
		"per_page":   paginator.PerPage,
		"total_page": paginator.TotalPages,
		"total_rows": paginator.TotalEntrySizes,
	}
	return c.Status(http.StatusOK).JSON(resp)
}

func (h productHandler) CreateProduct(c *fiber.Ctx) error {
	var ctx = c.Context()
	var product = new(models.Products)

	if err := c.BodyParser(product); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	product.SetCreatedAt()
	product.SetUpdatedAt()

	/* ทำการรับ Files จาก Form */
	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Cast Form Failed.")
	}
	files := form.File["files"]

	if err := h.proUs.CraeteProduct(ctx, product, files); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	var resp = map[string]interface{}{
		"message": "created.",
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (h productHandler) UpdateProduct(c *fiber.Ctx) error {
	var ctx = c.Context()
	var newProduct = new(models.Products)
	var productId = c.Params("product_id")
	if err := c.BodyParser(newProduct); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	/* ทำการรับ Files จาก Form */
	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Cast Form Failed.")
	}
	files := form.File["files"]

	existProduct, err := h.proUs.FetchOneProduct(ctx, productId)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	if existProduct == nil {
		return fiber.NewError(http.StatusNoContent, constants.ERROR_CAN_NOT_FIND_PRODUCT)
	}

	/* Merge Data && Images Managements */
	newProduct.MergeProduct(existProduct)
	delImages, delURL := newProduct.FindDeleteImage(existProduct)
	var delReq = make([]*models.DeleteFileReq, 0)
	if len(delURL) > 0 {
		for index := range delURL {
			req := &models.DeleteFileReq{
				Destination: delURL[index],
			}
			delReq = append(delReq, req)
		}
	}

	/* Delete Images Google Cloud Platform && Database */
	if len(delReq) > 0 && len(delImages) > 0 {
		if err := h.fileUs.DeleteOnGCP(delReq); err != nil {
			return fiber.NewError(http.StatusInternalServerError, fmt.Sprintf("DeleteOnGCP failed: %v", err))
		}
		if err := h.proUs.DeleteImages(ctx, delImages); err != nil {
			return fiber.NewError(http.StatusInternalServerError, fmt.Sprintf("DeleteImage failed: %v", err))
		}
	}

	if err := h.proUs.UpdateProduct(ctx, newProduct, files); err != nil {
		return fiber.NewError(http.StatusInternalServerError, fmt.Sprintf("UpdateProduct failed: %v", err))
	}

	resp := map[string]interface{}{
		"message": "successful.",
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (h productHandler) DeleteProduct(c *fiber.Ctx) error {
	var ctx = c.Context()
	var productId = c.Params("product_id")

	existProduct, err := h.proUs.FetchOneProduct(ctx, productId)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, fmt.Sprintf("FetchOneFailed: %v", err))
	}
	if existProduct == nil {
		return fiber.NewError(http.StatusNoContent, "there is no product with this id.")
	}

	var delReq = make([]*models.DeleteFileReq, 0)
	if len(existProduct.Images) > 0 {
		for index := range existProduct.Images {
			url := existProduct.Images[index].Url
			prefix := "https://storage.googleapis.com/pheethy-dev-bucket/"
			result := strings.SplitAfter(url, prefix)
			del := &models.DeleteFileReq{
				Destination: result[1],
			}
			delReq = append(delReq, del)
		}
	}

	if len(delReq) > 0 {
		if err := h.fileUs.DeleteOnGCP(delReq); err != nil {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	if err := h.proUs.DeleteProduct(ctx, productId); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"message": "successful",
	}

	return c.Status(http.StatusOK).JSON(resp)
}

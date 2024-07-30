package handlers

import (
	"mime/multipart"
	"net/http"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/file"
	"pheet-fiber-backend/service/products"
	"strconv"
	"strings"
	"sync"

	"github.com/Pheethy/psql/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

/* Adapter */
type productsHandlers struct {
	productUs products.IProductUsecase
	fileUs    file.IFileUsecase
}

/* Constructor */
func NewProductsHandlers(productUs products.IProductUsecase, fileUs file.IFileUsecase) products.IProductsHandlers {
	return productsHandlers{
		productUs: productUs,
		fileUs:    fileUs,
	}
}

func (p productsHandlers) FetchAllProduct(c *fiber.Ctx) error {
	ctx := c.Context()
	searchWord := c.Query("search_word")
	page, pageErr := strconv.Atoi(c.Query("page"))
	perPage, perPageErr := strconv.Atoi(c.Query("per_page"))
	args := new(sync.Map)
	paginator := helper.NewPaginator()

	if searchWord != "" {
		args.Store("search_word", searchWord)
	}

	if pageErr == nil {
		paginator.Page = page
	}

	if perPageErr == nil {
		paginator.PerPage = perPage
	}

	products, err := p.productUs.FetchAllProducts(ctx, args, &paginator)
	if err != nil {
		if ok := strings.Contains(err.Error(), constants.ERROR_PRODUCT_NOT_FOUND); ok {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": constants.ERROR_PRODUCT_NOT_FOUND,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
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

func (p productsHandlers) FetchOneProduct(c *fiber.Ctx) error {
	ctx := c.Context()
	productId := uuid.FromStringOrNil(c.Params("product_id"))

	product, err := p.productUs.FetchOneProduct(ctx, &productId)
	if err != nil {
		if ok := strings.Contains(err.Error(), constants.ERROR_PRODUCT_NOT_FOUND); ok {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": constants.ERROR_PRODUCT_NOT_FOUND,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp := map[string]interface{}{
		"product": product,
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (p productsHandlers) Create(c *fiber.Ctx) error {
	ctx := c.Context()
	req := new(models.Products)

	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	/* ทำการรับ Files จาก Form */
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	files := form.File["files"]
	req.NewId()
	req.SetCreatedAt()
	req.SetUpdatedAt()

	if err := p.productUs.CraeteProduct(ctx, req, files); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp := map[string]interface{}{
		"message":    "success",
		"product_id": req.Id,
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (p productsHandlers) UpdateProduct(c *fiber.Ctx) error {
	ctx := c.Context()
	newProduct := new(models.Products)
	productId := uuid.FromStringOrNil(c.Params("product_id"))
	if err := c.BodyParser(newProduct); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	/* ทำการรับ Files จาก Form */
	form, _ := c.MultipartForm()
	files := make([]*multipart.FileHeader, 0)
	if form != nil {
		files = form.File["files"]
	}

	existProduct, err := p.productUs.FetchOneProduct(ctx, &productId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if existProduct == nil {
		return c.Status(http.StatusNoContent).JSON(fiber.Map{
			"error": constants.ERROR_PRODUCT_NOT_FOUND,
		})
	}

	/* Merge Data && Images Managements */
	newProduct.MergeProduct(existProduct)
	newProduct.SetUpdatedAt()
	delImages, delURL := newProduct.FindDeleteImage(existProduct)
	delReq := make([]*models.DeleteFileReq, 0)
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
		if err := p.fileUs.DeleteOnGCP(delReq); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"DeleteOnGCP failed:": err.Error(),
			})
		}
		if err := p.productUs.DeleteImages(ctx, delImages); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"DeleteImages failed:": err.Error(),
			})
		}
	}
	if err := p.productUs.UpdateProduct(ctx, newProduct, files); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"Update products failed:": err.Error(),
		})
	}

	resp := map[string]interface{}{
		"message": "successful",
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (p productsHandlers) DeleteProduct(c *fiber.Ctx) error {
	ctx := c.Context()
	productId := uuid.FromStringOrNil(c.Params("product_id"))

	product, err := p.productUs.FetchOneProduct(ctx, &productId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	delReq := make([]*models.DeleteFileReq, 0)
	if len(product.Images) > 0 {
		for index := range product.Images {
			url := product.Images[index].URL
			prefix := "https://storage.googleapis.com/pheethy-dev-bucket/"
			result := strings.SplitAfter(url, prefix)
			del := &models.DeleteFileReq{
				Destination: result[1],
			}
			delReq = append(delReq, del)
		}
	}

	if len(delReq) > 0 {
		if err := p.fileUs.DeleteOnGCP(delReq); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "delete on GCP failed.",
			})
		}
	}

	if err := p.productUs.DeleteProduct(ctx, product); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := map[string]interface{}{
		"message": "successful",
	}

	return c.Status(http.StatusOK).JSON(resp)
}

package handler

import (
	"main/models"
	"main/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type productHandler struct {
	proSrv service.ProductService
}

func NewProductHandler(proSrv service.ProductService) productHandler {
	return productHandler{proSrv: proSrv}
}

func (h productHandler) GetProducts(c *fiber.Ctx) error {
	products, err := h.proSrv.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	resp := map[string]interface{}{
		"products": products,
	}

	return c.JSON(resp)
}

func (h productHandler) GetProductById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	customer, err := h.proSrv.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	resp := map[string]interface{}{
		"customer": customer,
	}
	
	return c.JSON(resp)
	
}

func (h productHandler) CreateProduct(c *fiber.Ctx) error {
	var newProduct = models.Product{}
	err := c.BodyParser(&newProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	err = h.proSrv.Create(&newProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	resp := map[string]interface{}{
		"massage": "created.",
	}

	return c.JSON(resp)
}

func (h productHandler) UpdateProduct(c *fiber.Ctx) error {
	var newProduct = models.Product{}
	err := c.BodyParser(&newProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	err = h.proSrv.Update(&newProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	resp := map[string]interface{}{
		"massage": "updated.",
	}

	return c.JSON(resp)
}

func (h productHandler) DeleteProduct(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	err = h.proSrv.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	resp := map[string]interface{}{
		"massage": "deleted.",
	}

	return c.JSON(resp)
}
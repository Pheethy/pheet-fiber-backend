package handler

import (
	"main/models"
	"main/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type customerHandler struct {
	custSrv service.CustomerService
}

func NewCustomerHandler(custSrv service.CustomerService) customerHandler {
	return customerHandler{custSrv: custSrv}
}

func (h customerHandler) GetProducts(c *fiber.Ctx) error {
	products, err := h.custSrv.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	resp := map[string]interface{}{
		"products": products,
	}

	return c.JSON(resp)
}

func (h customerHandler) GetProductById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	customer, err := h.custSrv.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	resp := map[string]interface{}{
		"customer": customer,
	}
	
	return c.JSON(resp)
	
}

func (h customerHandler) CreateProduct(c *fiber.Ctx) error {
	var newProduct = models.Product{}
	err := c.BodyParser(&newProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	err = h.custSrv.Create(&newProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	resp := map[string]interface{}{
		"massage": "created.",
	}

	return c.JSON(resp)
}

func (h customerHandler) UpdateProduct(c *fiber.Ctx) error {
	var newProduct = models.Product{}
	err := c.BodyParser(&newProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	err = h.custSrv.Update(&newProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	resp := map[string]interface{}{
		"massage": "updated.",
	}

	return c.JSON(resp)
}

func (h customerHandler) DeleteProduct(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	err = h.custSrv.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	resp := map[string]interface{}{
		"massage": "deleted.",
	}

	return c.JSON(resp)
}
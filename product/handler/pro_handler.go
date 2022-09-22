package handler

import (
	"log"
	"main/auth"
	"main/models"
	"main/product"
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type productHandler struct {
	proSrv product.ProductService
}

func NewProductHandler(proSrv product.ProductService) productHandler {
	return productHandler{proSrv: proSrv}
}

func (h productHandler) Login(c *fiber.Ctx) error {
	var request = models.User{}

	err := c.BodyParser(&request)
	if err != nil {
		return c.JSON(err.Error())
	}

	if request.UserName == "" || request.Password == "" {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	user, err := h.proSrv.GetUser(request.UserName)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "incorrect username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(request.Password))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "incorrect username or password")
	}

	tokenz := auth.AccessToken(os.Getenv("SIGN"))

	resp := map[string]interface{}{
		"message": "Login-success",
		"jwt": tokenz,
	}

	return c.JSON(resp)
}

func (h productHandler) SignUp(c *fiber.Ctx) error {
	var request = models.SignUpReq{}

	err := c.BodyParser(&request)
	if err != nil {
		return c.JSON(err.Error())
	}

	if request.UserName == "" || request.Password == "" {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError)
	}
	request.Password = string(password)

	err = h.proSrv.SignUp(&request)
	if err != nil {
		return c.JSON(err.Error())
	}

	resp := map[string]interface{}{
		"massage": "signUp-sucCess",
	}

	return c.JSON(resp)
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

	product, err := h.proSrv.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	resp := map[string]interface{}{
		"product": product,
	}
	
	return c.JSON(resp)
}

func (h productHandler) GetProductByType(c *fiber.Ctx) error {
	coffType := c.Params("type")

	product, err := h.proSrv.GetProductByType(coffType)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	resp := map[string]interface{}{
		"products": product,
	}
	
	return c.JSON(resp)
}

func (h productHandler) CreateProduct(c *fiber.Ctx) error {
	var newProduct = models.Product{}
	err := c.BodyParser(&newProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	log.Println("newProductRequest:", newProduct)

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

	log.Println("RequestUpdateModels",newProduct)

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
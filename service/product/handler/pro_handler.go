package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"pheet-fiber-backend/auth"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/product"
	validate "pheet-fiber-backend/service/product/validator"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type productHandler struct {
	proSrv product.ProductUsecase
}

func NewProductHandler(proSrv product.ProductUsecase) product.ProductHandler {
	return productHandler{proSrv: proSrv}
}

func (h productHandler) Login(c *fiber.Ctx) error {
	var request = models.User{}

	err := c.BodyParser(&request)
	if err != nil {
		return c.JSON(err.Error())
	}

	if request.Username == "" || request.Password == "" {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	user, err := h.proSrv.GetUser(request.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "incorrect username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "incorrect username or password")
	}

	tokenz := auth.AccessToken(os.Getenv("SIGN"))

	resp := map[string]interface{}{
		"message": "Login-success",
		"jwt":     tokenz,
	}

	return c.JSON(resp)
}

func (h productHandler) SignUp(c *fiber.Ctx) error {
	var ctx = c.Context()
	var request = new(models.User)
	var time = models.NewTimestampFromTime(time.Now())
	err := c.BodyParser(&request)
	if err != nil {
		return c.JSON(err.Error())
	}

	if request.Username == "" || request.Password == "" {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError)
	}
	request.Password = string(password)
	request.NewUUID()
	request.SetCreatedAt(&time)
	request.SetUpdatedAt(&time)

	err = h.proSrv.SignUp(ctx, request)
	if err != nil {
		return c.JSON(err.Error())
	}

	resp := map[string]interface{}{
		"massage": "signUp-sucCess",
	}

	return c.JSON(resp)
}

func (h productHandler) GetProducts(c *fiber.Ctx) error {
	var ctx = c.Context()
	products, err := h.proSrv.GetProducts(ctx)
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
	var ctx = c.Context()
	var body = c.Body()
	var params interface{}
	err := json.Unmarshal(body, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	var newProduct = new(models.Products)
	var time = models.NewTimestampFromTime(time.Now())
	err = c.BodyParser(newProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}
	ele := validate.ValidateStruct(newProduct)
	if len(ele) > 0 {
		for _, val := range ele {
			if val.FailedField != "" {
				return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("%s: was missing on body",val.FailedField))
			}
			if val.Tag != "" {
				return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("%s: was missing on body",val.Tag))
			}
			if val.Value != "" {
				return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("%s: was missing on body",val.Value))
			}
		}
	}

	newProduct.NewUUID()
	newProduct.SetCreatedAt(&time)
	newProduct.SetUpdatedAt(&time)

	err = h.proSrv.Create(ctx, newProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	resp := map[string]interface{}{
		"massage": "product has created.",
	}

	return c.JSON(resp)
}

func (h productHandler) UpdateProduct(c *fiber.Ctx) error {
	var ctx = c.Context()
	var newProduct = new(models.Products)
	var id = uuid.FromStringOrNil(c.Params("product_id"))
	var time = models.NewTimestampFromTime(time.Now())

	err := c.BodyParser(newProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}
	newProduct.SetUpdatedAt(&time) //set update timestamp

	err = h.proSrv.Update(ctx, newProduct, &id)
	if err != nil {
		c.JSON(http.StatusInternalServerError)
	}

	resp := map[string]interface{}{
		"massage": "product has updated.",
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

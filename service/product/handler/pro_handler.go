package handler

import (
	"net/http"
	"os"
	"pheet-fiber-backend/auth"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/helper"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/product"
	"strconv"
	"strings"
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
	var time = helper.NewTimestampFromTime(time.Now())
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

func (p productHandler) Create(c *fiber.Ctx) error {
	var ctx = c.Context()
	var err error
	var imgPath string
	var proReq = new(models.Products)
	var time = helper.NewTimestampFromTime(time.Now())

	/* Binding Data && Validation*/
	if err := c.BodyParser(proReq); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := proReq.ValidationStruct(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	imgFile, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"image": "must be required",
		})
	}

	/* Image Management */
	if imgFile != nil {
		imgPath = "./uploads/products/" + imgFile.Filename
		if err := c.SaveFile(imgFile, imgPath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	/* Setting id, time and imagePath */
	proReq.NewUUID()
	proReq.SetCreatedAt(&time)
	proReq.SetUpdatedAt(&time)
	proReq.SetImagePath(imgPath)

	/* Sending Create */
	if err := p.proSrv.Create(ctx, proReq); err != nil {
		if strings.Contains(err.Error(), constants.ERROR_PRODUCTNAME_WAS_DUPLICATE) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error:": constants.ERROR_PRODUCTNAME_WAS_DUPLICATE,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp := fiber.Map{
		"message": "successful",
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h productHandler) UpdateProduct(c *fiber.Ctx) error {
	var ctx = c.Context()
	var newProduct = new(models.Products)
	var id = uuid.FromStringOrNil(c.Params("product_id"))
	var time = helper.NewTimestampFromTime(time.Now())

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

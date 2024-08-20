package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	file_mocks "pheet-fiber-backend/service/file/mocks"
	product_mocks "pheet-fiber-backend/service/products/mocks"
	"testing"
	"time"

	"github.com/Pheethy/psql/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchOneProduct(t *testing.T) {
	idPro := uuid.FromStringOrNil("68eeb8bf-1ae3-4010-ad74-47f5332b3d5b")
	idImg := uuid.FromStringOrNil("4d8c4a78-e6c1-4234-8965-e22515144ca0")
	ti := helper.NewTimestampFromTime(time.Now())
	products := &models.Products{
		TableName:    struct{}{},
		Id:           &idPro,
		Title:        "ทดสอบ",
		Description:  "ทดสอบ",
		Price:        100,
		CreatedAt:    &ti,
		UpdatedAt:    &ti,
		CategoriesId: 1,
		Categories:   &models.Categories{Id: 1, Title: "ตลาดมาก", ProductId: &idPro},
		Images: []*models.Image{
			{
				ID:        &idImg,
				FileName:  "test",
				URL:       "test",
				ProductId: &idPro,
				CreatedAt: &ti,
				UpdatedAt: &ti,
			},
		},
	}
	t.Run("success", func(t *testing.T) {
		productUs := new(product_mocks.IProductUsecase)
		fileUs := new(file_mocks.IFileUsecase)
		productUs.On("FetchOneProduct", mock.Anything, mock.AnythingOfType("*uuid.UUID")).Return(products, nil).Run(func(args mock.Arguments) {
			epCtx := args.Get(0)
			epId := args.Get(1).(*uuid.UUID)

			assert.NotNil(t, epCtx)
			assert.Equal(t, epId.String(), idPro.String())
		})
		app := fiber.New()
		productHandlers := NewProductsHandlers(productUs, fileUs)
		app.Get("/v1/product/:product_id", func(c *fiber.Ctx) error {
			return productHandlers.FetchOneProduct(c)
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", idPro.String()), nil)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, resp)
	})
	t.Run("error_product_not_found", func(t *testing.T) {
		productUs := new(product_mocks.IProductUsecase)
		fileUs := new(file_mocks.IFileUsecase)
		productUs.On("FetchOneProduct", mock.Anything, mock.AnythingOfType("*uuid.UUID")).Return(nil, errors.New(constants.ERROR_PRODUCT_NOT_FOUND)).Run(func(args mock.Arguments) {
			epCtx := args.Get(0)
			epId := args.Get(1).(*uuid.UUID)

			assert.NotNil(t, epCtx)
			assert.Equal(t, epId.String(), idPro.String())
		})
		app := fiber.New()
		productHandlers := NewProductsHandlers(productUs, fileUs)
		app.Get("/v1/product/:product_id", func(c *fiber.Ctx) error {
			return productHandlers.FetchOneProduct(c)
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", idPro.String()), nil)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})
	t.Run("error_internal_server", func(t *testing.T) {
		productUs := new(product_mocks.IProductUsecase)
		fileUs := new(file_mocks.IFileUsecase)
		productUs.On("FetchOneProduct", mock.Anything, mock.AnythingOfType("*uuid.UUID")).Return(nil, errors.New("unexpected")).Run(func(args mock.Arguments) {
			epCtx := args.Get(0)
			epId := args.Get(1).(*uuid.UUID)

			assert.NotNil(t, epCtx)
			assert.Equal(t, epId.String(), idPro.String())
		})
		app := fiber.New()
		productHandlers := NewProductsHandlers(productUs, fileUs)
		app.Get("/v1/product/:product_id", func(c *fiber.Ctx) error {
			return productHandlers.FetchOneProduct(c)
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/product/%s", idPro.String()), nil)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}

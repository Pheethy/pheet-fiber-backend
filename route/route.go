package route

import (
	"pheet-fiber-backend/service/product"
	"github.com/gofiber/fiber/v2"
)

type Route struct {
	e	fiber.Router
}

func NewRoute(e fiber.Router) *Route {
	return &Route{e: e}
}

func (r Route) RegisterProduct(handler product.ProductHandler) {
	r.e.Get("/products", handler.GetProducts)
	r.e.Get("/product/:id", handler.GetProductById)
	r.e.Get("/products/:type", handler.GetProductByType)
	r.e.Post("/product", handler.CreateProduct)
	r.e.Put("/product/:product_id", handler.UpdateProduct)
	r.e.Delete("/product/:id", handler.DeleteProduct)
}
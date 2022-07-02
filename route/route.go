package route

import (
	"main/handler"
	"github.com/gofiber/fiber/v2"
)

type Route struct {
	e	fiber.Router
}

func NewRoute(e fiber.Router) *Route {
	return &Route{e: e}
}

func (r Route) RegisterProduct(handler handler.ProductHandler) {
	r.e.Get("/products", handler.GetProducts)
	r.e.Get("/product/:id", handler.GetProductById)
	r.e.Get("/products/:type", handler.GetProductByType)
	r.e.Post("/product", handler.CreateProduct)
	r.e.Put("/product", handler.UpdateProduct)
	r.e.Delete("/product/:id", handler.DeleteProduct)
}
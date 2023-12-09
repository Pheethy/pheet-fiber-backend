package route

import (
	"pheet-fiber-backend/middleware"
	"pheet-fiber-backend/service/appinfo"
	"pheet-fiber-backend/service/file"
	"pheet-fiber-backend/service/order"
	"pheet-fiber-backend/service/product"
	"pheet-fiber-backend/service/users"

	"github.com/gofiber/fiber/v2"
)

type Route struct {
	e fiber.Router
}

func NewRoute(e fiber.Router) *Route {
	return &Route{e: e}
}

func (r Route) RegisterUsers(handler users.IUsersHandlers, m middleware.ImiddlewareHandler) {
	r.e.Post("/users/sign-in", handler.GetPassport)
	r.e.Post("/users/sign-up", handler.SignUpCustomer)
	r.e.Post("/users/sign-out", m.ApiKeyAuth(), handler.SignOut)
	r.e.Post("/users/refresh", m.ApiKeyAuth(), handler.RefreshPassport)
	r.e.Get("/users/secret", m.JwtAuth(), m.Authorize(2), handler.GenerateAdminToken)
	r.e.Get("/users/:user_id", m.JwtAuth(), m.ParamsCheck(), m.ApiKeyAuth(), handler.FetchUserProfile)
}

func (r Route) RegisterAppInfo(handler appinfo.AppInfoHandler, m middleware.ImiddlewareHandler) {
	r.e.Get("/info/apikey", m.JwtAuth(), m.Authorize(1), handler.GenerateAPIKey)
	r.e.Get("/info/category", m.ApiKeyAuth(), handler.FindCategory)
	r.e.Post("/info/category", handler.AddCategory)
	r.e.Delete("/info/category/:category_id", handler.RemoveCategory)
}

func (r Route) RegisterFile(handler file.IFileHandler, m middleware.ImiddlewareHandler) {
	r.e.Post("/file/upload",m.JwtAuth(), m.Authorize(1), m.ApiKeyAuth(), handler.UploadFile)
	r.e.Patch("/file/delete",m.JwtAuth(), m.Authorize(1), m.ApiKeyAuth(), handler.DeleteFile)
}

func (r Route) RegisterProduct(handler product.IProductHandler, m middleware.ImiddlewareHandler) {
	r.e.Get("/product/:product_id", m.ApiKeyAuth(), handler.FetchOneProduct)
	r.e.Get("/product", m.ApiKeyAuth(), handler.FetchAllProduct)
	r.e.Post("/product", m.ApiKeyAuth(), m.JwtAuth(), m.Authorize(1), handler.CreateProduct)
	r.e.Put("/product/:product_id", m.ApiKeyAuth(), m.JwtAuth(), m.Authorize(1), handler.UpdateProduct)
	r.e.Delete("/product/:product_id", m.ApiKeyAuth(), m.JwtAuth(), m.Authorize(1), handler.DeleteProduct)
}

func (r Route) RegisterOrder(handler order.IOrderHandler, m middleware.ImiddlewareHandler) {
	
}
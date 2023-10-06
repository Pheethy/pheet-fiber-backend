package route

import (
	"pheet-fiber-backend/middleware"
	"pheet-fiber-backend/service/appinfo"
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
}

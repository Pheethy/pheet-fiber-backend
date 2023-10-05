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
	r.e.Post("users/sign-in", handler.GetPassport)
	r.e.Post("users/sign-up", handler.SignUpCustomer)
	r.e.Post("users/sign-out", handler.SignOut)
	r.e.Post("users/refresh", handler.RefreshPassport)
	r.e.Get("users/secret",m.JwtAuth(), m.Authorize(2), handler.GenerateAdminToken)
	r.e.Get("users/:user_id", m.JwtAuth(), m.ParamsCheck(), handler.FetchUserProfile)
}

func (r Route) RegisterAppInfo(handler appinfo.AppInfoHandler, m middleware.ImiddlewareHandler) {

}
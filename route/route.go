package route

import (
	"pheet-fiber-backend/middleware"
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
	r.e.Post("/sign-in", handler.GetPassport)
	r.e.Post("/sign-up", handler.SignUpCustomer)
	r.e.Post("/sign-out", handler.SignOut)
	r.e.Post("/refresh", handler.RefreshPassport)
	r.e.Get("/secret",m.JwtAuth(), m.Authorize(2), handler.GenerateAdminToken)
	r.e.Get("/:user_id", m.JwtAuth(), m.ParamsCheck(), handler.FetchUserProfile)
}

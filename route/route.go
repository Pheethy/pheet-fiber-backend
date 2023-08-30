package route

import (
	"pheet-fiber-backend/service/users"

	"github.com/gofiber/fiber/v2"
)

type Route struct {
	e fiber.Router
}

func NewRoute(e fiber.Router) *Route {
	return &Route{e: e}
}

func (r Route) RegisterUsers(handler users.IUsersHandlers) {
	r.e.Post("/signup/customer", handler.SignUpCustomer)
}

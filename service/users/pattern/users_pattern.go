package pattern

import (
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/users"

	"github.com/jmoiron/sqlx"
)

type userReq struct {
	id string
	req *models.UserRegisterReq
	db *sqlx.DB
}

type customer struct {
	*userReq
}

type admin struct {
	*userReq
}

func InsertUser(db *sqlx.DB, req *models.UserRegisterReq, IsAdmin bool) users.IUsersPattern {
	if IsAdmin {
		return newAdmin()
	} else {
		return newCustomer()
	}
}

func newCustomer() users.IUsersPattern {
	return nil
}

func newAdmin() users.IUsersPattern {
	return nil
}
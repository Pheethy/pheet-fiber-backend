package repository

import (
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/service/order"

	"github.com/Pheethy/sqlx"
)

type orderRepository struct {
	db  *sqlx.DB
	cfg config.Iconfig
}

func NewOrderRepository(db *sqlx.DB, cfg config.Iconfig) order.IOrderRepository {
	return orderRepository{
		db:  db,
		cfg: cfg,
	}
}

package order

import (
	"context"
	"pheet-fiber-backend/models"
	"sync"

	"github.com/Pheethy/psql/helper"
)

type IOrderRepository interface {
	FetchAllOrder(ctx context.Context, args *sync.Map, paginator *helper.Paginator) ([]*models.Order, error)
	FetchOneOrder(ctx context.Context, orderId string) (*models.Order, error)
}
package usecase

import (
	"context"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/order"
	"pheet-fiber-backend/service/product"
	"sync"

	"github.com/Pheethy/psql/helper"
)

type orderUsecase struct {
	orderRepo order.IOrderRepository
	productRepo product.IProductRepository
}

func NewOrderUsecase(orderRepo order.IOrderRepository, productRepo product.IProductRepository) order.IOrderUsecase {
	return orderUsecase{
		orderRepo: orderRepo,
		productRepo: productRepo,
	}
}

func (o orderUsecase) FetchAllOrder(ctx context.Context, args *sync.Map, paginator *helper.Paginator) ([]*models.Order, error) {
	return o.orderRepo.FetchAllOrder(ctx, args, paginator)
}

func (o orderUsecase) FetchOneOrder(ctx context.Context, orderId string) (*models.Order, error) {
	return o.orderRepo.FetchOneOrder(ctx, orderId)
}
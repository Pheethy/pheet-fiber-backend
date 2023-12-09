package usecase

import (
	"pheet-fiber-backend/service/order"
	"pheet-fiber-backend/service/product"
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
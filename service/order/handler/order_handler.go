package handler

import "pheet-fiber-backend/service/order"

type orderHandler struct {
	orderUs order.IOrderUsecase
}

func NewOrderHandler(orderUs order.IOrderUsecase) order.IOrderHandler {
	return orderHandler{
		orderUs: orderUs,
	}
}
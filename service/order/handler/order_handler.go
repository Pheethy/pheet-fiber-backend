package handler

import (
	"net/http"
	"pheet-fiber-backend/service/order"
	"strconv"
	"sync"

	"github.com/Pheethy/psql/helper"
	"github.com/gofiber/fiber/v2"
)

type orderHandler struct {
	orderUs order.IOrderUsecase
}

func NewOrderHandler(orderUs order.IOrderUsecase) order.IOrderHandler {
	return orderHandler{
		orderUs: orderUs,
	}
}

func (o orderHandler) FetchAllOrder(c *fiber.Ctx) error {
	var ctx = c.Context()
	var args = new(sync.Map)
	var paginator = new(helper.Paginator)
	var page, pageErr = strconv.Atoi(c.Query("page"))
	var perPage, perPageErr = strconv.Atoi(c.Query("per_page"))
	var userId = c.Query("user_id")

	if userId != "" {
		args.Store("user_id", userId)
	}

	if pageErr == nil {
		paginator.Page = page
	}

	if perPageErr == nil {
		paginator.PerPage = perPage
	}

	orders, err := o.orderUs.FetchAllOrder(ctx, args, paginator)
	if err != nil {
		return nil
	}

	resp := map[string]interface{}{
		"orders":     orders,
		"page":       paginator.Page,
		"per_page":   paginator.PerPage,
		"total_page": paginator.TotalPages,
		"total_rows": paginator.TotalEntrySizes,
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (o orderHandler) FetchOneOrder(c *fiber.Ctx) error {
	var ctx = c.Context()
	var orderId = c.Params("order_id")

	order, err := o.orderUs.FetchOneOrder(ctx, orderId)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	if order == nil {
		return fiber.NewError(http.StatusNoContent)
	}

	resp := map[string]interface{}{
		"order": order,
	}

	return c.Status(http.StatusOK).JSON(resp)
}

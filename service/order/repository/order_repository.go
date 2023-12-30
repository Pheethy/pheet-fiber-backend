package repository

import (
	"context"
	"errors"
	"fmt"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/order"
	"sync"

	"github.com/Pheethy/psql/helper"
	"github.com/Pheethy/psql/orm"
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

func (o orderRepository) FetchAllOrder(ctx context.Context, args *sync.Map, paginator *helper.Paginator) ([]*models.Order, error) {
	sql := fmt.Sprintf(`
	SELECT
		%s,
		%s,
		%s
	FROM
		orders
	JOIN
		transfer_slip
	ON
		orders.id = transfer_slip.order_id
	JOIN
		products_orders
	ON
		orders.id = products_orders.order_id
	`,
		orm.GetSelector(models.Order{}),
		orm.GetSelector(models.TransferSlip{}),
		orm.GetSelector(models.ProductOrder{}),
	)

	stmt, err := o.db.PreparexContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	options := orm.NewMapperOption()
	orders, err := o.orms(ctx, rows, options)
	if err != nil {
		return nil, err
	}

	if len(orders) < 1 {
		return nil, errors.New("no content")
	}

	return orders, nil
}

func (o orderRepository) FetchOneOrder(ctx context.Context, orderId string) (*models.Order, error) {
	sql := fmt.Sprintf(`
	SELECT
		%s,
		%s,
		products_orders.product_id "products.id",
		products_orders.products_title "products.title",
		products_orders.products_description "products.description",
		products_orders.products_price "products.price",
		products_orders.products_created_at "products.created_at",
		products_orders.products_updated_at "products.updated_at",
		%s
	FROM
		orders
	JOIN
		(
			SELECT
				products_orders.*,
				products.id "products_id",
				products.title "products_title",
				products.description "products_description",
				products.price "products_price",
				products.created_at "products_created_at",
				products.updated_at "products_updated_at"
			FROM
				products_orders
			JOIN
				products
			ON
				products_orders.product_id = products.id
		) AS products_orders
	ON
		orders.id = products_orders.order_id
	JOIN
		transfer_slip
	ON
		orders.id = transfer_slip.order_id
	WHERE
		orders.id = $1
	`,
		orm.GetSelector(models.Order{}),
		orm.GetSelector(models.ProductOrder{}),
		orm.GetSelector(models.TransferSlip{}),
	)

	stmt, err := o.db.PreparexContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	options := orm.NewMapperOption().SetOverridePKField(
		orm.NewMapperOptionPKField(new(models.Products), []string{models.FIELD_PRODUCTS_ID}),
	)

	order, err := o.orm(ctx, rows, options)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o orderRepository) orms(ctx context.Context, rows *sqlx.Rows, options orm.MapperOption) ([]*models.Order, error) {
	mapper, err := orm.OrmContext(ctx, new(models.Order), rows, options)
	if err != nil {
		return nil, err
	}
	orders := mapper.GetData().([]*models.Order)
	if len(orders) == 0 {
		return nil, errors.New("order not found")
	}
	return orders, nil
}

func (o orderRepository) orm(ctx context.Context, rows *sqlx.Rows, options orm.MapperOption) (*models.Order, error) {
	mapper, err := orm.OrmContext(ctx, new(models.Order), rows, options)
	if err != nil {
		return nil, err
	}

	orders := mapper.GetData().([]*models.Order)
	if len(orders) == 0 {
		return nil, errors.New("order not found")
	}

	return orders[0], nil
}

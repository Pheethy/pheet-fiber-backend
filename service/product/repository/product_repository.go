package repository

import (
	"context"
	"fmt"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/orm"
	"pheet-fiber-backend/service/product"

	"github.com/BlackMocca/sqlx"
)

type productRepository struct {
	db  *sqlx.DB
	cfg config.Iconfig
}

func NewProductRepository(db *sqlx.DB, cfg config.Iconfig) product.IProductRepository {
	return productRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r productRepository) FetchOneProduct(ctx context.Context, id string) (*models.Products, error) {
	sql := fmt.Sprintf(`
		SELECT
			%s,
			%s
		FROM
			products
		JOIN
		 	images
		ON
			products.id = images.product_id
		WHERE
			products.id = $1
	`,
		orm.GetSelector(models.Products{}),
		orm.GetSelector(models.Image{}),
	)

	stmt, err := r.db.PreparexContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapper, err := orm.Orm(new(models.Products), rows, orm.NewMapperOption())
	if err != nil {
		panic(err)
	}

	products := mapper.GetData().([]*models.Products)

	if len(products) > 0 {
		for index := range products {
			categories, err := r.FetchCategoriesByProductId(ctx, products[index].ID)
			if err != nil {
				return nil, err
			}
			products[index].Categories = categories
		}
	}

	return products[0], nil
}

func (r productRepository) FetchCategoriesByProductId(ctx context.Context, productId string) (*models.Categories, error) {
	sql := fmt.Sprintf(`
		SELECT
			%s
		FROM
			categories
		JOIN
			products_categories
		ON
			categories.id = products_categories.category_id
		WHERE
			products_categories.product_id = $1
	`,
		orm.GetSelector(models.Categories{}),
	)

	stmt, err := r.db.PreparexContext(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(productId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	mapper, err := orm.Orm(new(models.Categories), rows, orm.NewMapperOption())
	if err != nil {
		return nil, fmt.Errorf("orm failed: %v", err)
	}

	categories := mapper.GetData().([]*models.Categories)

	return categories[0], nil
}

package repository

import (
	"context"
	"fmt"
	"log"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/products"
	"strings"
	"sync"

	"github.com/Pheethy/psql/helper"
	"github.com/Pheethy/psql/orm"
	"github.com/Pheethy/sqlx"
	"github.com/gofrs/uuid"
)

/* Adapter */
type productsRepository struct {
	psqlDB *sqlx.DB
}

/* Constructor */
func NewProductsRepository(psqlDB *sqlx.DB) products.IProductsRepositoryDB {
	return productsRepository{
		psqlDB: psqlDB,
	}
}

func (p productsRepository) whereCond(args *sync.Map) ([]string, []interface{}) {
	conds := []string{}
	var valArgs []interface{}

	if v, ok := args.Load("search_word"); ok {
		if v != nil {
			cond := "LOWER(products.title) LIKE CONCAT('%%',?::text,'%%')"
			conds = append(conds, cond)
			searchWord := strings.ToLower(v.(string))
			searchWord = strings.ReplaceAll(searchWord, " ", "")
			valArgs = append(valArgs, searchWord)
			valArgs = append(valArgs, searchWord)
		}
	}

	return conds, valArgs
}

func (p productsRepository) FetchAllProducts(ctx context.Context, args *sync.Map, paginator *helper.Paginator) ([]*models.Products, error) {
	conds, valArgs := p.whereCond(args)
	var where string
	var paginateSQL string
	if len(conds) > 0 {
		where += "WHERE " + strings.Join(conds, " AND ")
	}
	if paginator != nil {
		limit := int(paginator.PerPage)
		skipItem := (int(paginator.Page) - 1) * int(paginator.PerPage)
		paginateSQL = fmt.Sprintf(`
			LIMIT %d
			OFFSET %d
			`,
			limit,
			skipItem,
		)
	}

	sql := fmt.Sprintf(`
		SELECT
			%s,
			%s,
			%s,
			products.total_row
		FROM
			(
				SELECT
					*,
					COUNT(*) OVER() as "total_row"
				FROM
					products
				%s
				%s
			) as products
		JOIN
			(
				SELECT
					categories.*,
					products_categories.product_id "product_id"
				FROM
					categories
				JOIN
					products_categories
				ON
					products_categories.category_id = categories.id
			) AS categories
		ON
			products.id = categories.product_id
		JOIN
			images
		ON
			products.id = images.product_id
		%s
		ORDER BY
			products.created_at ASC
	`,
		orm.GetSelector(models.Products{}),
		orm.GetSelector(models.Categories{}),
		orm.GetSelector(models.Image{}),
		where,
		paginateSQL,
		where,
	)

	sql = sqlx.Rebind(sqlx.DOLLAR, sql)
	stmt, err := p.psqlDB.PreparexContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, valArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products, err := p.orms(ctx, rows, paginator)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p productsRepository) FetchOneProduct(ctx context.Context, productId *uuid.UUID) (*models.Products, error) {
	sql := fmt.Sprintf(`
		SELECT
			%s,
			%s,
			%s
		FROM
			products
		JOIN
			(
				SELECT
					categories.*,
					products_categories.product_id "product_id"
				FROM
					categories
				JOIN
					products_categories
				ON
					categories.id = products_categories.category_id
			) AS categories
		ON
			products.id = categories.product_id
		JOIN
			images
		ON
			products.id = images.product_id
		WHERE
			products.id = $1
	`,
		orm.GetSelector(models.Products{}),
		orm.GetSelector(models.Categories{}),
		orm.GetSelector(models.Image{}),
	)

	stmt, err := p.psqlDB.PreparexContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, productId)
	if err != nil {
		return nil, err
	}

	product, err := p.orm(ctx, rows)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p productsRepository) CraeteProduct(ctx context.Context, req *models.Products) error {
	tx, err := p.psqlDB.Beginx()
	if err != nil {
		return fmt.Errorf("begin failed: %v", err)
	}

	if err := p.upsertProduct(ctx, tx, req); err != nil {
		return fmt.Errorf("create product failed: %v", err)
	}

	if err := p.createProductsCategories(ctx, tx, req); err != nil {
		return fmt.Errorf("create products_categories failed: %v", err)
	}

	if err := p.upsertImages(ctx, tx, req); err != nil {
		return fmt.Errorf("create images failed: %v", err)
	}

	return tx.Commit()
}

func (p productsRepository) UpdateProduct(ctx context.Context, req *models.Products) error {
	tx, err := p.psqlDB.Beginx()
	if err != nil {
		return fmt.Errorf("begin failed: %v", err)
	}

	if err := p.updateProducts(ctx, tx, req); err != nil {
		return fmt.Errorf("update product failed: %v", err)
	}

	if err := p.updateProductsCategories(ctx, tx, req); err != nil {
		return fmt.Errorf("update products_categories failed: %v", err)
	}

	if err := p.upsertImages(ctx, tx, req); err != nil {
		return fmt.Errorf("create images failed: %v", err)
	}

	return tx.Commit()
}

func (p productsRepository) upsertProduct(ctx context.Context, tx *sqlx.Tx, product *models.Products) error {
	sql := `
		INSERT INTO products (id, title, description, price, created_at, updated_at)
		VALUES(
			$1::uuid,
			$2::text,
			$3::text,
			$4::float,
			$5::timestamp,
			$6::timestamp
		)
		ON CONFLICT (id)
		DO UPDATE SET
			title=$7::text,
			description=$8::text,
			price=$9::float,
			updated_at=$10::timestamp
	`
	stmt, err := p.psqlDB.PreparexContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx,
		/* create */
		product.Id,
		product.Title,
		product.Description,
		product.Price,
		product.CreatedAt,
		product.UpdatedAt,
		/* update */
		product.Title,
		product.Description,
		product.Price,
		product.UpdatedAt,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("exec failed: %v", err)
	}

	return nil
}

func (r productsRepository) createProductsCategories(ctx context.Context, tx *sqlx.Tx, product *models.Products) error {
	sql := `
		INSERT INTO products_categories (
			product_id,
			category_id
		)
		VALUES (
			$1::uuid,
			$2::integer
		)
	`

	stmt, err := tx.PreparexContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("prepare failed: %v", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx,
		product.Id,
		product.CategoriesId,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("exec failed: %v", err)
	}

	return nil
}

func (r productsRepository) updateProductsCategories(ctx context.Context, tx *sqlx.Tx, product *models.Products) error {
	sql := `
		UPDATE
			products_categories
		SET
			category_id=$1::int
		WHERE
			product_id=$2::uuid
	`

	stmt, err := tx.PreparexContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("prepare failed: %v", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx,
		product.CategoriesId,
		product.Id,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("exec failed: %v", err)
	}

	return nil
}

func (p productsRepository) updateProducts(ctx context.Context, tx *sqlx.Tx, product *models.Products) error {
	sql := `
		UPDATE
			products
		SET
			title=$1::text,
			description=$2::text,
			price=$3::float,
			updated_at=$4::timestamp
		WHERE
			id=$5::uuid
	`
	stmt, err := tx.PreparexContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx,
		product.Title,
		product.Description,
		product.Price,
		product.UpdatedAt,
		product.Id,
	); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r productsRepository) upsertImages(ctx context.Context, tx *sqlx.Tx, product *models.Products) error {
	sql := `
		INSERT INTO images (
	        id,
	        filename,
	        url,
	        product_id,
	        created_at,
	        updated_at
		) VALUES (
	        $1::uuid,
	        $2::text,
	        $3::text,
	        $4::uuid,
	        $5::timestamp,
	        $6::timestamp
		)
		ON CONFLICT (id)
		DO UPDATE SET
	        filename=$7::text,
	        url=$8::text,
	        product_id=$9::uuid,
	        updated_at=$10::timestamp
	`
	stmt, err := tx.PreparexContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("prepare failed: %v", err)
	}

	for index := range product.Images {
		if _, err := stmt.ExecContext(ctx,
			// create
			product.Images[index].ID,
			product.Images[index].FileName,
			product.Images[index].URL,
			product.Id,
			product.Images[index].CreatedAt,
			product.Images[index].UpdatedAt,
			// update
			product.Images[index].FileName,
			product.Images[index].URL,
			product.Id,
			product.Images[index].UpdatedAt,
		); err != nil {
			log.Printf("exec failed: %v\n", err)
			tx.Rollback()
			return fmt.Errorf("exec failed: %v", err)
		}
	}

	return nil
}

func (p productsRepository) DeleteProduct(ctx context.Context, product *models.Products) error {
	tx, err := p.psqlDB.Beginx()
	if err != nil {
		return nil
	}

	if err := p.deleteProductsCategories(ctx, tx, product.Id); err != nil {
		return err
	}

	if err := p.deleteImages(ctx, tx, product.Id); err != nil {
		return err
	}

	if err := p.deleteProduct(ctx, tx, product.Id); err != nil {
		return err
	}

	return tx.Commit()
}

func (p productsRepository) deleteProduct(ctx context.Context, tx *sqlx.Tx, productId *uuid.UUID) error {
	sql := `
		DELETE
		FROM
			products
		WHERE
			products.id = $1::uuid;
	`
	stmt, err := tx.PreparexContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("prepare failed: %v", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, productId); err != nil {
		return tx.Rollback()
	}

	return nil
}

func (p productsRepository) deleteImages(ctx context.Context, tx *sqlx.Tx, productId *uuid.UUID) error {
	sql := `
		DELETE
		FROM
			images
		WHERE
			product_id=$1::uuid;
	`
	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, productId); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (p productsRepository) deleteProductsCategories(ctx context.Context, tx *sqlx.Tx, productId *uuid.UUID) error {
	sql := `
		DELETE
		FROM
			products_categories
		WHERE
			product_id=$1::uuid;
	`
	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, productId); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (p productsRepository) DeleteImages(ctx context.Context, ids []*uuid.UUID) error {
	tx, err := p.psqlDB.Beginx()
	if err != nil {
		return err
	}
	sql := `
		DELETE FROM images
		WHERE id=$1::uuid;
	`
	stmt, err := tx.PreparexContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if len(ids) > 0 {
		for _, id := range ids {
			if _, err := stmt.ExecContext(ctx,
				id,
			); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (p productsRepository) orms(ctx context.Context, rows *sqlx.Rows, paginator *helper.Paginator) ([]*models.Products, error) {
	mapper, err := orm.OrmContext(ctx, new(models.Products), rows, orm.NewMapperOption())
	if err != nil {
		return nil, err
	}

	products := mapper.GetData().([]*models.Products)
	if paginator != nil {
		paginator.SetPaginatorByAllRows(mapper.GetPaginateTotal())
	}

	if len(products) == 0 {
		return nil, fmt.Errorf(constants.ERROR_PRODUCT_NOT_FOUND)
	}

	return products, nil
}

func (p productsRepository) orm(ctx context.Context, rows *sqlx.Rows) (*models.Products, error) {
	mapper, err := orm.OrmContext(ctx, new(models.Products), rows, orm.NewMapperOption())
	if err != nil {
		return nil, err
	}
	products := mapper.GetData().([]*models.Products)

	if len(products) == 0 {
		return nil, fmt.Errorf(constants.ERROR_PRODUCT_NOT_FOUND)
	}

	return products[0], nil
}

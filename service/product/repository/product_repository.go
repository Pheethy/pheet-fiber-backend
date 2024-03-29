package repository

import (
	"context"
	"errors"
	"fmt"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/product"
	"strings"
	"sync"

	"github.com/Pheethy/psql/helper"
	"github.com/Pheethy/psql/orm"
	"github.com/Pheethy/sqlx"
	"github.com/gofrs/uuid"
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

func (r productRepository) whereCond(args *sync.Map) ([]string, []interface{}) {
	var conds = []string{}
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

func (r productRepository) FetchAllProduct(ctx context.Context, args *sync.Map, paginate *helper.Paginator) ([]*models.Products, error) {
	conds, vals := r.whereCond(args)
	var where string
	var paginateSQL string
	if len(conds) > 0 {
		where += "WHERE " + strings.Join(conds, " AND ")
	}
	if paginate != nil {
		var limit = int(paginate.PerPage)
		var skipItem = (int(paginate.Page) - 1) * int(paginate.PerPage)
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
		 	images
		ON
			products.id = images.product_id
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
		%s
		ORDER BY
			products.created_at ASC
	`,
		orm.GetSelector(models.Products{}),
		orm.GetSelector(models.Image{}),
		orm.GetSelector(models.Categories{}),
		where,
		paginateSQL,
		where,
	)
	sql = sqlx.Rebind(sqlx.DOLLAR, sql)
	stmt, err := r.db.PreparexContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, vals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products, err := r.orms(ctx, rows, paginate)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r productRepository) FetchOneProduct(ctx context.Context, id string) (*models.Products, error) {
	sql := fmt.Sprintf(`
		SELECT
			%s,
			%s,
			%s
		FROM
			products
		JOIN
		 	images
		ON
			products.id = images.product_id
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
		WHERE
			products.id = $1
	`,
		orm.GetSelector(models.Products{}),
		orm.GetSelector(models.Image{}),
		orm.GetSelector(models.Categories{}),
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

	product, err := r.orm(ctx, rows)
	if err != nil {
		return nil, err
	}

	return product, nil
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

func (r productRepository) CraeteProduct(ctx context.Context, req *models.Products) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin failed: %v", err)
	}

	if err := r.createProduct(ctx, req); err != nil {
		return fmt.Errorf("create product failed: %v", err)
	}

	if err := r.createProductsCategories(ctx, tx, req); err != nil {
		return fmt.Errorf("create products_categories failed: %v", err)
	}

	if err := r.upsertImages(ctx, tx, req); err != nil {
		return fmt.Errorf("create images failed: %v", err)
	}

	return tx.Commit()
}

func (r productRepository) createProduct(ctx context.Context, product *models.Products) error {
	sql := `
		INSERT INTO products (title, description, price, created_at, updated_at)
		VALUES(
			$1::text,
			$2::text,
			$3::float,
			$4::timestamp,
			$5::timestamp
		)
		RETURNING "id";
	`
	stmt, err := r.db.PreparexContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx,
		product.Title,
		product.Description,
		product.Price,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID); err != nil {
		if strings.Contains(err.Error(), constants.ERROR_PQ_UNIQUE_PRODUCTNAME) {
			return errors.New(constants.ERROR_PRODUCTNAME_WAS_DUPLICATE)
		}
		return err
	}

	return nil
}
func (r productRepository) createProductsCategories(ctx context.Context, tx *sqlx.Tx, product *models.Products) error {
	sql := `
		INSERT INTO products_categories (
			product_id,
			category_id
		)
		VALUES (
			$1::text,
			$2::int
		)
	`

	stmt, err := tx.PreparexContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("prepare failed: %v", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx,
		product.ID,
		product.CategoriesId,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("exec failed: %v", err)
	}

	return nil
}

func (r productRepository) upsertImages(ctx context.Context, tx *sqlx.Tx, product *models.Products) error {
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
			$4::text,
			$5::timestamp,
			$6::timestamp
		)
		ON CONFLICT (id)
		DO UPDATE SET
			filename=$7::text,
			url=$8::text,
			product_id=$9::text,
			updated_at=$10::timestamp
	`
	stmt, err := tx.PreparexContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("prepare failed: %v", err)
	}

	for index := range product.Images {
		if _, err := stmt.ExecContext(ctx,
			//create
			product.Images[index].ID,
			product.Images[index].FilenName,
			product.Images[index].Url,
			product.ID,
			product.Images[index].CreatedAt,
			product.Images[index].UpdatedAt,
			//update
			product.Images[index].FilenName,
			product.Images[index].Url,
			product.ID,
			product.Images[index].UpdatedAt,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("exec failed: %v", err)
		}
	}

	return nil
}

func (r productRepository) UpdateProduct(ctx context.Context, product *models.Products) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	if err := r.updateProducts(ctx, tx, product); err != nil {
		return fmt.Errorf("update product failed: %v", err)
	}

	if err := r.upsertImages(ctx, tx, product); err != nil {
		return fmt.Errorf("update image failed: %v", err)
	}

	if err := r.updateProductsCategories(ctx, tx, product); err != nil {
		return fmt.Errorf("update products_categories failed: %v", err)
	}

	return tx.Commit()
}

func (r productRepository) updateProducts(ctx context.Context, tx *sqlx.Tx, product *models.Products) error {
	sql := `
		UPDATE
			products
		SET
			title=$1::text,
			description=$2::text,
			price=$3::float,
			updated_at=$4::timestamp
		WHERE
			id=$5::text
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
		product.ID,
	); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r productRepository) updateImage(ctx context.Context, tx *sqlx.Tx, images []*models.Image) error {
	sql := `
		UPDATE
			images
		SET
			filename=$1::text,
			url=$2::text,
			updated_at=$3::timestamp
		WHERE
			id=$4::uuid
	`
	stmt, err := tx.PreparexContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if len(images) > 0 {
		for index := range images {
			if _, err := stmt.ExecContext(ctx,
				images[index].FilenName,
				images[index].Url,
				images[index].UpdatedAt,
				images[index].ID,
			); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return nil
}

func (r productRepository) updateProductsCategories(ctx context.Context, tx *sqlx.Tx, product *models.Products) error {
	sql := `
		UPDATE
			products_categories
		SET
			category_id=$1::int
		WHERE
			product_id=$2::text
	`

	stmt, err := tx.PreparexContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx,
		product.CategoriesId,
		product.ID,
	); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r productRepository) DeleteProduct(ctx context.Context, productId string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	if err := r.deleteProductsCategories(ctx, tx, productId); err != nil {
		return fmt.Errorf("delete products_categories failed: %v", err)
	}
	if err := r.deleteImages(ctx, tx, productId); err != nil {
		return fmt.Errorf("delete images failed: %v", err)
	}
	if err := r.deleteProduct(ctx, tx, productId); err != nil {
		return fmt.Errorf("delete product failed: %v", err)
	}

	return tx.Commit()
}

func (r productRepository) deleteProduct(ctx context.Context, tx *sqlx.Tx, productId string) error {
	sql := `
		DELETE
		FROM
			products
		WHERE
			id=$1::text;
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

func (r productRepository) deleteImages(ctx context.Context, tx *sqlx.Tx, productId string) error {
	sql := `
		DELETE
		FROM
			images
		WHERE
			product_id=$1::text;
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

func (r productRepository) deleteProductsCategories(ctx context.Context, tx *sqlx.Tx, productId string) error {
	sql := `
		DELETE
		FROM
			products_categories
		WHERE
			product_id=$1::text;
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

func (r productRepository) DeleteImages(ctx context.Context, ids []*uuid.UUID) error {
	tx, err := r.db.Beginx()
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

func (r productRepository) orm(ctx context.Context, rows *sqlx.Rows) (*models.Products, error) {
	mapper, err := orm.OrmContext(ctx, new(models.Products), rows, orm.NewMapperOption())
	if err != nil {
		panic(err)
	}

	products := mapper.GetData().([]*models.Products)

	if len(products) == 0 {
		return nil, errors.New("product not found")
	}

	return products[0], nil
}

func (r productRepository) orms(ctx context.Context, rows *sqlx.Rows, paginator *helper.Paginator) ([]*models.Products, error) {
	mapper, err := orm.OrmContext(ctx, new(models.Products), rows, orm.NewMapperOption())
	if err != nil {
		panic(err)
	}

	products := mapper.GetData().([]*models.Products)
	if paginator != nil {
		paginator.SetPaginatorByAllRows(mapper.GetPaginateTotal())
	}

	if len(products) == 0 {
		return nil, fmt.Errorf("product not found: %v", err)
	}

	return products, nil
}

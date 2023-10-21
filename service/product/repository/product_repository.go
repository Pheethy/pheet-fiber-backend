package repository

import (
	"context"
	"fmt"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/helper"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/orm"
	"pheet-fiber-backend/service/product"
	"strings"
	"sync"

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

		}
	}

	return conds, valArgs
}

func (r productRepository) FetchAllProduct(ctx context.Context, args *sync.Map, paginate *helper.Paginator) ([]*models.Products, error) {
	conds, vals := r.whereCond(args)
	var where string
	var paginateSQL string
	if len(conds) > 0 {
		where += " WHERE " + strings.Join(conds, " AND ")
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
			%s
		FROM
			products
		JOIN
		 	images
		ON
			products.id = images.product_id
		%s
		ORDER BY
			products.created_at ASC
		%s
	`,
		orm.GetSelector(models.Products{}),
		orm.GetSelector(models.Image{}),
		where,
		paginateSQL,
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

func (r productRepository) orm(ctx context.Context, rows *sqlx.Rows) (*models.Products, error) {
	mapper, err := orm.OrmContext(ctx, new(models.Products), rows, orm.NewMapperOption())
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

	if len(products) == 0 {
		return nil, fmt.Errorf("product not found: %v", err)
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

	if len(products) > 0 {
		for index := range products {
			categories, err := r.FetchCategoriesByProductId(ctx, products[index].ID)
			if err != nil {
				return nil, err
			}
			products[index].Categories = categories
		}
	}

	// /* worker pools */
	// var jobsCh = make(chan *models.Products, len(products))
	// var errCh = make(chan error, len(products))

	// /* ทำการนำ products ใส่ไปใน jobs channel */
	// for _, r := range products {
	// 	jobsCh <- r
	// }
	// close(jobsCh)

	// /* สร้าง worker */
	// var worker int = 10
	// for i := 0; i < worker; i++ {
	// 	//working zone
	// 	go r.fillCategories(ctx, jobsCh, errCh)
	// }

	// /* สร้าง loop สำหรับการรับ error */
	// for a := 0; a < len(products); a++ {
	// 	//handler err โดยการรับค่า err จาก Channel errCh
	// 	if err := <-errCh; err != nil {
	// 		return nil, fmt.Errorf("err: %v", err)
	// 	}
	// }

	if len(products) == 0 {
		return nil, fmt.Errorf("product not found: %v", err)
	}

	return products, nil
}

func (r productRepository) fillCategories(ctx context.Context, jobs <-chan *models.Products, errs chan<- error) {
	for job := range jobs {
		category, err := r.FetchCategoriesByProductId(ctx, job.ID)
		if err != nil {
			errs <- err
		}
		job.Categories = category
		/* กรณีไม่มี error ก็ต้องทำการ return ค่า nil ออกไป errCh เพราะเราประกาศรับค่าไว้ */
		errs <- nil
	}
}

package repository

import (
	"context"
	"fmt"
	"log"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/appinfo"
	"strings"
	"sync"

	"github.com/BlackMocca/sqlx"
)

type appInfoRepository struct {
	db *sqlx.DB
}

func NewAppInfoRepository(db *sqlx.DB) appinfo.AppInfoRepository {
	return &appInfoRepository{db: db}
}

func (r appInfoRepository) whereCond(args *sync.Map) ([]string, []interface{}) {
	var cond = make([]string, 0)
	var vals = make([]interface{}, 0)

	if v, ok := args.Load("search_word"); ok {
		search := `(LOWER("title") LIKE $1)`
		cond = append(cond, search)
		vals = append(vals, "%"+strings.ToLower(v.(string))+"%")
	}

	return cond, vals
}

func (r appInfoRepository) FindCategory(ctx context.Context, args *sync.Map) ([]*models.Categories, error) {
	conds, vals := r.whereCond(args)
	var where string
	if len(conds) > 0 {
		where = "WHERE" + strings.Join(conds, "AND ")
	}

	sql := fmt.Sprintf(`
		SELECT
			"id",
			"title"
		FROM
			"categories"
		%s
	`,
		where,
	)

	var cats = make([]*models.Categories, 0)
	if err := r.db.SelectContext(ctx, &cats, sql, vals...); err != nil {
		return nil, fmt.Errorf("Insert catagories failed: %v", err)
	}

	return cats, nil
}

func (r appInfoRepository) InsertCategories(ctx context.Context, cats []*models.Categories) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	sql := `
		INSERT INTO "categories" (
			"title"
		) VALUES (
			$1::text
		)
		RETURNING "id";
	`
	stmt, err := tx.PreparexContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if len(cats) > 0 && cats != nil {
		for index := range cats {
			if err := stmt.QueryRowContext(ctx, sql, cats[index].Title).Scan(&cats[index].Id); err != nil {
				log.Fatal(err)
				tx.Rollback()
			}
		}
	}

	return tx.Commit()
}

func (r appInfoRepository) DeleteCategory(ctx context.Context, id int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	sql := `
		DELETE
		FROM
			"categories"
		WHERE
			"id" = $1
	`

	stmt, err := tx.PreparexContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, id); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete categories failed: %v", err)
	}

	return tx.Commit()
}
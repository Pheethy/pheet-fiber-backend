package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/product"

	"github.com/jmoiron/sqlx"
)

/* Adapter entity conform Interface Pod*/
type productRepositoryDB struct {
	psqlDB *sqlx.DB
}

// constructor //
func NewProductRepository(db *sqlx.DB) product.ProductRepository {
	return productRepositoryDB{psqlDB: db}
}

func (r productRepositoryDB) FetchAll(ctx context.Context) ([]*models.Products, error) {
	sql := `
	SELECT
		id, name, detail, type, price, cover, created_at, updated_at
	FROM
		products
	`
	var products []*models.Products
	err := r.psqlDB.SelectContext(ctx, &products, sql)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r productRepositoryDB) FetchByType(coffType string) ([]*models.Products, error) {
	sql := fmt.Sprintf(`
	SELECT
		id, name, type, price, description, image
	FROM
		product
	WHERE
		type = '%s'
	`, coffType)
	var products []*models.Products
	err := r.psqlDB.Select(&products, sql)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r productRepositoryDB) FetchById(id int) (*models.Products, error) {
	sql := `
	SELECT
		id, name, type, price, description, image
	FROM
		product
	WHERE
		id=?
	`
	var product models.Products
	err := r.psqlDB.Get(&product, sql, id)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r productRepositoryDB) FetchUser(username string) (*models.User, error) {
	var user = new(models.User)
	sql := `
	SELECT
		id, username, password
	FROM
		users
	WHERE
		username = $1
	`
	err := r.psqlDB.Get(user, sql, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r productRepositoryDB) Create(ctx context.Context, product *models.Products) error {
	tx, err := r.psqlDB.Beginx()
	if err != nil {
		return err
	}
	sql := `
		INSERT INTO products (id, name, detail, type, price, cover, created_at, updated_at)
		VALUES(
			$1::uuid,
			$2::text,
			$3::text,
			$4::product_type,
			$5::int,
			$6::text,
			$7::timestamp,
			$8::timestamp
		)
	`
	stmt, err := tx.Preparex(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx,
		product.Id,
		product.Name,
		product.Detail,
		product.Type,
		product.Price,
		product.Image,
		product.CreatedAt,
		product.UpdatedAt,
	); err != nil {
		log.Println("err db:", err)
	}
	
	return tx.Commit()
}

func (r productRepositoryDB) SignUp(ctx context.Context, user *models.User) error {
	tx, err := r.psqlDB.Beginx()
	if err != nil {
		return err
	}
	sql := `
	INSERT INTO
		users (id, username, password, created_at, updated_at)
	VALUES(
		$1::uuid,
		$2::text,
		$3::text,
		$4::timestamp,
		$5::timestamp
	)
	`
	stmt, err := tx.Preparex(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx,
		user.Id,
		user.Username,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	); err != nil {
		log.Println("err db:", err)
	}
	
	return tx.Commit()
}

func (r productRepositoryDB) Update(product *models.Products) error {
	// sql := `
	// 	UPDATE 
	// 		product
	// 	SET
	// 		name = ?,
	// 		type = ?,
	// 		price = ?,
	// 		description = ?,
	// 		image = ?
	// 	WHERE
	// 		id = ?
	// `

	// result, err := r.psqlDB.Exec(sql, product.Name, product.Type, product.Price, product.Description, product.Image, product.Id)
	// if err != nil {
	// 	panic(err)
	// }

	// affected, err := result.RowsAffected()
	// if err != nil {
	// 	return err
	// }

	// if affected < 1 {
	// 	return errors.New("Update fail")
	// }

	return nil
}

func (r productRepositoryDB) Delete(id int) error {
	sql := `
	DELETE FROM 
		product
	WHERE
		id = ?
	`
	result, err := r.psqlDB.Exec(sql, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected < 1 {
		return errors.New("Delete fail")
	}

	return nil
}

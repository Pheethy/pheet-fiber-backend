package repository

import (
	"context"
	"errors"
	"fmt"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/product"

	"github.com/jmoiron/sqlx"
)

/* Adapter entity conform Interface Pod*/
type productRepositoryDB struct {
	msqlDB *sqlx.DB
}

// constructor //
func NewProductRepository(db *sqlx.DB) product.ProductRepository {
	return productRepositoryDB{msqlDB: db}
}

func (r productRepositoryDB) FetchAll(ctx context.Context) ([]*models.Product, error) {
	sql := `
	SELECT
		id, name, type, price, description, image
	FROM
		product
	`
	var products []*models.Product
	err := r.msqlDB.SelectContext(ctx, &products, sql)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r productRepositoryDB) FetchByType(coffType string) ([]*models.Product, error) {
	sql := fmt.Sprintf(`
	SELECT
		id, name, type, price, description, image
	FROM
		product
	WHERE
		type = '%s'
	`, coffType)
	var products []*models.Product
	err := r.msqlDB.Select(&products, sql)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r productRepositoryDB) FetchById(id int) (*models.Product, error) {
	sql := `
	SELECT
		id, name, type, price, description, image
	FROM
		product
	WHERE
		id=?
	`
	var product models.Product
	err := r.msqlDB.Get(&product, sql, id)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r productRepositoryDB) Create(product *models.Product) error {
	sql := `
	INSERT INTO
		product (id, name, type, price, description, image)
	VALUES
		(?, ?, ?, ?, ?, ?)
	`
	result, err := r.msqlDB.Exec(sql, product.Id, product.Name, product.Type, product.Price, product.Description, product.Image)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected < 1 {
		return errors.New("Create fail")
	}

	return nil
}

func (r productRepositoryDB) SignUp(user *models.SignUpReq) error {
	sql := `
	INSERT INTO
		user (username, password)
	VALUES
		(?, ?)
	`
	result, err := r.msqlDB.Exec(sql, user.UserName, user.Password)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected < 1 {
		return errors.New("SignUp Error")
	}

	return nil
}

func (r productRepositoryDB) FetchUser(username string) (*models.User, error) {
	sql := `
	SELECT
		id, username, password
	FROM
		user
	WHERE
		username = ?
	`
	var user models.User
	err := r.msqlDB.Get(&user, sql, username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r productRepositoryDB) Update(product *models.Product) error {
	sql := `
		UPDATE 
			product
		SET
			name = ?,
			type = ?,
			price = ?,
			description = ?,
			image = ?
		WHERE
			id = ?
	`

	result, err := r.msqlDB.Exec(sql, product.Name, product.Type, product.Price, product.Description, product.Image, product.Id)
	if err != nil {
		panic(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected < 1 {
		return errors.New("Update fail")
	}

	return nil
}

func (r productRepositoryDB) Delete(id int) error {
	sql := `
	DELETE FROM 
		product
	WHERE
		id = ?
	`
	result, err := r.msqlDB.Exec(sql, id)
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

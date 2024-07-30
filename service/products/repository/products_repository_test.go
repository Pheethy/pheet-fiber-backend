package repository

import (
	"context"
	"database/sql"
	"pheet-api-flavorparser/models"
	"testing"
	"time"

	"github.com/Pheethy/psql"
	"github.com/Pheethy/psql/helper"
	"github.com/Pheethy/sqlx"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v2"
)

func setupMockDB() (*sql.DB, *sqlx.DB, sqlmock.Sqlmock, error) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, err
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	client := new(psql.Client)
	client.SetDB(sqlxDB)

	return db, sqlxDB, sqlMock, err
}

func Test_Fetch_One_Product(t *testing.T) {
	db, sqlxDB, sqlMock, err := setupMockDB()
	assert.NoError(t, err)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		idPro := uuid.FromStringOrNil("68eeb8bf-1ae3-4010-ad74-47f5332b3d5b")
		idImg := uuid.FromStringOrNil("4d8c4a78-e6c1-4234-8965-e22515144ca0")
		ti := helper.NewTimestampFromTime(time.Now())
		products := &models.Products{
			TableName:    struct{}{},
			Id:           &idPro,
			Title:        "ทดสอบ",
			Description:  "ทดสอบ",
			Price:        100,
			CreatedAt:    &ti,
			UpdatedAt:    &ti,
			CategoriesId: 1,
			Categories:   &models.Categories{Id: 1, Title: "ตลาดมาก", ProductId: &idPro},
			Images: []*models.Images{
				{
					Id:        &idImg,
					FileName:  "test",
					URL:       "test",
					ProductId: &idPro,
					CreatedAt: &ti,
					UpdatedAt: &ti,
				},
			},
		}

		rows := sqlmock.NewRows([]string{
			"products.id", "products.title", "products.description", "products.price", "products.created_at", "products.updated_at",
		})
		rows.AddRow(
			products.Id,
			products.Title,
			products.Description,
			products.Price,
			products.CreatedAt,
			products.UpdatedAt,
		)

		catRows := sqlmock.NewRows([]string{
			"categories.id",
			"categories.title",
			"categories.product_id",
		})
		catRows.AddRow(
			products.Categories.Id,
			products.Categories.Title,
			products.Categories.ProductId,
		)

		imgRows := sqlmock.NewRows([]string{
			"images.id",
			"images.filename",
			"images.url",
			"images.product_id",
			"images.created_at",
			"images.updated_at",
		})
		for index := range products.Images {
			imgRows.AddRow(
				products.Images[index].Id,
				products.Images[index].FileName,
				products.Images[index].URL,
				products.Images[index].ProductId,
				products.Images[index].CreatedAt,
				products.Images[index].UpdatedAt,
			)
		}

		sqlProduct := `
      SELECT
        (.+)
      FROM
        (.+)
      WHERE
        (.+)
    `

		sqlCategories := `
      SELECT
        (.+)
      FROM
        (.+)
      WHERE
        (.+)
    `

		sqlImg := `
      SELECT
        (.+)
      FROM
        (.+)
      WHERE
        (.+)
    `

		sqlMock.ExpectPrepare(sqlProduct).ExpectQuery().WithArgs().WillReturnRows(rows)
		sqlMock.ExpectPrepare(sqlCategories).ExpectQuery().WithArgs().WillReturnRows(catRows)
		sqlMock.ExpectPrepare(sqlImg).ExpectQuery().WithArgs().WillReturnRows(imgRows)

		repo := NewProductsRepository(sqlxDB)

		epProduct, err := repo.FetchOneProduct(context.Background(), &idPro)

		assert.NoError(t, err)
		assert.NotNil(t, epProduct)

		assert.Equal(t, epProduct.Id, products.Id)
		assert.Equal(t, epProduct.Title, products.Title)
		assert.Equal(t, epProduct.Description, products.Description)
		assert.Equal(t, epProduct.Price, products.Price)
		assert.Equal(t, epProduct.CreatedAt.String(), products.CreatedAt.String())
		assert.Equal(t, epProduct.UpdatedAt.String(), products.UpdatedAt.String())
	})
	t.Run("content_not_found", func(t *testing.T) {
		idPro := uuid.FromStringOrNil("68eeb8bf-1ae3-4010-ad74-47f5332b3d5b")
		rows := sqlmock.NewRows([]string{
			"products.id", "products.title", "products.description", "products.price", "products.created_at", "products.updated_at",
		})

		catRows := sqlmock.NewRows([]string{
			"categories.id",
			"categories.title",
			"categories.product_id",
		})

		imgRows := sqlmock.NewRows([]string{
			"images.id",
			"images.filename",
			"images.url",
			"images.product_id",
			"images.created_at",
			"images.updated_at",
		})

		sqlProduct := `
      SELECT
        (.+)
      FROM
        (.+)
      WHERE
        (.+)
      `

		sqlCategories := `
      SELECT
        (.+)
      FROM
        (.+)
      WHERE
        (.+)
      `

		sqlImg := `
      SELECT
        (.+)
      FROM
        (.+)
      WHERE
        (.+)
      `

		sqlMock.ExpectPrepare(sqlProduct).ExpectQuery().WithArgs().WillReturnRows(rows)
		sqlMock.ExpectPrepare(sqlCategories).ExpectQuery().WithArgs().WillReturnRows(catRows)
		sqlMock.ExpectPrepare(sqlImg).ExpectQuery().WithArgs().WillReturnRows(imgRows)

		repo := NewProductsRepository(sqlxDB)

		epProduct, err := repo.FetchOneProduct(context.Background(), &idPro)

		assert.Error(t, err)
		assert.Nil(t, epProduct)
	})
}

func TestCraeteProduct(t *testing.T) {
	db, sqlxDB, sqlMock, err := setupMockDB()
	assert.NoError(t, err)
	defer db.Close()
	repo := NewProductsRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		id := uuid.FromStringOrNil("713f522b-947e-4674-b146-46dc2e4d80e1")
		idImg := uuid.FromStringOrNil("faa938bc-0b46-4510-b264-f7c105369386")
		ti := helper.NewTimestampFromTime(time.Now())
		mockProduct := &models.Products{
			Id:          &id,
			Title:       "Hey! bro.",
			Description: "Testional",
			Price:       150,
			CreatedAt:   &ti,
			UpdatedAt:   &ti,
			Categories: &models.Categories{
				Id:        1,
				Title:     "computer",
				ProductId: &id,
			},
			Images: []*models.Images{
				{
					Id:        &idImg,
					FileName:  "pheet.png",
					URL:       "test/path",
					ProductId: &id,
					CreatedAt: &ti,
					UpdatedAt: &ti,
				},
			},
		}
		sqlUpsertProduct := `
      INSERT INTO products
    `
		sqlUpsertProductCategories := `
      INSERT INTO products_categories
    `
		sqlUpsertImages := `
      INSERT INTO images
    `

		sqlMock.ExpectBegin()
		sqlMock.ExpectPrepare(sqlUpsertProduct).ExpectExec().WithArgs(
			/* create */
			mockProduct.Id,
			mockProduct.Title,
			mockProduct.Description,
			mockProduct.Price,
			mockProduct.CreatedAt,
			mockProduct.UpdatedAt,
			/* update */
			mockProduct.Title,
			mockProduct.Description,
			mockProduct.Price,
			mockProduct.UpdatedAt,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		sqlMock.ExpectPrepare(sqlUpsertProductCategories).ExpectExec().WithArgs(
			mockProduct.Id,
			mockProduct.CategoriesId,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		images := sqlMock.ExpectPrepare(sqlUpsertImages)
		for index := range mockProduct.Images {
			images.ExpectExec().WithArgs(
				/* create */
				mockProduct.Images[index].Id,
				mockProduct.Images[index].FileName,
				mockProduct.Images[index].URL,
				mockProduct.Images[index].ProductId,
				mockProduct.Images[index].CreatedAt,
				mockProduct.Images[index].UpdatedAt,
				/* update */
				mockProduct.Images[index].FileName,
				mockProduct.Images[index].URL,
				mockProduct.Images[index].ProductId,
				mockProduct.Images[index].UpdatedAt,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		}

		sqlMock.ExpectCommit()

		err := repo.CraeteProduct(context.Background(), mockProduct)
		assert.NoError(t, err)
	})
}

func TestUpdateProduct(t *testing.T) {
	db, sqlxDB, sqlMock, err := setupMockDB()
	assert.NoError(t, err)
	defer db.Close()
	repo := NewProductsRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		id := uuid.FromStringOrNil("713f522b-947e-4674-b146-46dc2e4d80e1")
		idImg := uuid.FromStringOrNil("faa938bc-0b46-4510-b264-f7c105369386")
		ti := helper.NewTimestampFromTime(time.Now())
		mockProduct := &models.Products{
			Id:          &id,
			Title:       "Hey! bro.",
			Description: "Testional",
			Price:       150,
			CreatedAt:   &ti,
			UpdatedAt:   &ti,
			Categories: &models.Categories{
				Id:        1,
				Title:     "computer",
				ProductId: &id,
			},
			Images: []*models.Images{
				{
					Id:        &idImg,
					FileName:  "pheet.png",
					URL:       "test/path",
					ProductId: &id,
					CreatedAt: &ti,
					UpdatedAt: &ti,
				},
			},
		}

		sqlUpdateProduct := `
    UPDATE
      products
    `
		sqlUpsertImages := `
      INSERT INTO images
    `
		sqlUpdateProductsCategories := `
    UPDATE
			products_categories
    `

		sqlMock.ExpectBegin()

		sqlMock.ExpectPrepare(sqlUpdateProduct).ExpectExec().WithArgs(
			mockProduct.Title,
			mockProduct.Description,
			mockProduct.Price,
			mockProduct.UpdatedAt,
			mockProduct.Id,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		sqlMock.ExpectPrepare(sqlUpdateProductsCategories).ExpectExec().WithArgs(
			mockProduct.CategoriesId,
			mockProduct.Id,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		images := sqlMock.ExpectPrepare(sqlUpsertImages)
		for index := range mockProduct.Images {
			images.ExpectExec().WithArgs(
				/* create */
				mockProduct.Images[index].Id,
				mockProduct.Images[index].FileName,
				mockProduct.Images[index].URL,
				mockProduct.Images[index].ProductId,
				mockProduct.Images[index].CreatedAt,
				mockProduct.Images[index].UpdatedAt,
				/* update */
				mockProduct.Images[index].FileName,
				mockProduct.Images[index].URL,
				mockProduct.Images[index].ProductId,
				mockProduct.Images[index].UpdatedAt,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		}

		sqlMock.ExpectCommit()
		err := repo.UpdateProduct(context.Background(), mockProduct)
		assert.NoError(t, err)
	})
}

func TestDeleteProduct(t *testing.T) {
	db, sqlxDB, sqlMock, err := setupMockDB()
	assert.NoError(t, err)
	defer db.Close()
	repo := NewProductsRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		id := uuid.FromStringOrNil("713f522b-947e-4674-b146-46dc2e4d80e1")
		idImg := uuid.FromStringOrNil("faa938bc-0b46-4510-b264-f7c105369386")
		ti := helper.NewTimestampFromTime(time.Now())
		mockProduct := &models.Products{
			Id:          &id,
			Title:       "Hey! bro.",
			Description: "Testional",
			Price:       150,
			CreatedAt:   &ti,
			UpdatedAt:   &ti,
			Categories: &models.Categories{
				Id:        1,
				Title:     "computer",
				ProductId: &id,
			},
			Images: []*models.Images{
				{
					Id:        &idImg,
					FileName:  "pheet.png",
					URL:       "test/path",
					ProductId: &id,
					CreatedAt: &ti,
					UpdatedAt: &ti,
				},
			},
		}
		sqlMock.ExpectBegin()
		sqlDeleteProductsCategories := `DELETE FROM products_categories`
		sqlDeleteImages := `DELETE FROM images`
		sqlDeleteProducts := `DELETE FROM products`
		sqlMock.ExpectPrepare(sqlDeleteProductsCategories).ExpectExec().WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
		sqlMock.ExpectPrepare(sqlDeleteImages).ExpectExec().WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
		sqlMock.ExpectPrepare(sqlDeleteProducts).ExpectExec().WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
		sqlMock.ExpectCommit()

		repo.DeleteProduct(context.Background(), mockProduct)
	})
}

package products

import (
	"context"
	"mime/multipart"
	"pheet-fiber-backend/models"
	"sync"

	"github.com/Pheethy/psql/helper"
	"github.com/gofrs/uuid"
)

/* Usecase Port */
type IProductUsecase interface {
	FetchAllProducts(ctx context.Context, args *sync.Map, paginator *helper.Paginator) ([]*models.Products, error)
	FetchOneProduct(ctx context.Context, productId *uuid.UUID) (*models.Products, error)
	CraeteProduct(ctx context.Context, req *models.Products, files []*multipart.FileHeader) error
	UpdateProduct(ctx context.Context, product *models.Products, files []*multipart.FileHeader) error
	DeleteProduct(ctx context.Context, product *models.Products) error
	DeleteImages(ctx context.Context, ids []*uuid.UUID) error
}

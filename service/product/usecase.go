package product

import (
	"context"
	"mime/multipart"
	"github.com/Pheethy/psql/helper"
	"pheet-fiber-backend/models"
	"sync"

	"github.com/gofrs/uuid"
)

type IProductUsecase interface {
	FetchOneProduct(ctx context.Context, id string) (*models.Products, error)
	FetchAllProduct(ctx context.Context, args *sync.Map, paginate *helper.Paginator) ([]*models.Products, error)
	CraeteProduct(ctx context.Context, req *models.Products, files []*multipart.FileHeader) error
	UpdateProduct(ctx context.Context, product *models.Products, files []*multipart.FileHeader) error
	DeleteProduct(ctx context.Context, productId string) error
	DeleteImages(ctx context.Context, ids []*uuid.UUID) error
}
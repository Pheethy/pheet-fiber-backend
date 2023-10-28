package product

import (
	"context"
	"mime/multipart"
	"pheet-fiber-backend/helper"
	"pheet-fiber-backend/models"
	"sync"
)

type IProductUsecase interface {
	FetchOneProduct(ctx context.Context, id string) (*models.Products, error)
	FetchAllProduct(ctx context.Context, args *sync.Map, paginate *helper.Paginator) ([]*models.Products, error)
	CraeteProduct(ctx context.Context, req *models.Products, files []*multipart.FileHeader) error
	UpdateProduct(ctx context.Context, product *models.Products) error
}
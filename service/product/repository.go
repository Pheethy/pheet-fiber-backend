package product

import (
	"context"
	"pheet-fiber-backend/helper"
	"pheet-fiber-backend/models"
	"sync"
)

type IProductRepository interface {
	FetchOneProduct(ctx context.Context, id string) (*models.Products, error)
	FetchAllProduct(ctx context.Context, args *sync.Map, paginate *helper.Paginator) ([]*models.Products, error)
}
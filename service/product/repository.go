package product

import (
	"context"
	"pheet-fiber-backend/models"
)

type IProductRepository interface {
	FetchOneProduct(ctx context.Context, id string) (*models.Products, error)
}
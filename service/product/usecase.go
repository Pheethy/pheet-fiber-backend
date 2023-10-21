package product

import (
	"context"
	"pheet-fiber-backend/models"
)

type IProductUsecase interface {
	FetchOneProduct(ctx context.Context, id string) (*models.Products, error)
}
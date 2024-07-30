package products

import (
	"context"
	"pheet-fiber-backend/models"
	"sync"

	"github.com/Pheethy/psql/helper"
	"github.com/gofrs/uuid"
)

/* Repository Port */
type IProductsRepositoryDB interface {
	FetchAllProducts(ctx context.Context, args *sync.Map, paginator *helper.Paginator) ([]*models.Products, error)
	FetchOneProduct(ctx context.Context, productId *uuid.UUID) (*models.Products, error)
	CraeteProduct(ctx context.Context, req *models.Products) error
	UpdateProduct(ctx context.Context, req *models.Products) error
	DeleteProduct(ctx context.Context, product *models.Products) error
	DeleteImages(ctx context.Context, ids []*uuid.UUID) error
}

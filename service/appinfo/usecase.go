package appinfo

import (
	"context"
	"pheet-fiber-backend/models"
	"sync"
)

type AppInfoUsecase interface {
	FindCategory(ctx context.Context, args *sync.Map) ([]*models.Catagory, error)
	InsertCategories(ctx context.Context, cats []*models.Catagory) error
	DeleteCategory(ctx context.Context, id int) error
}
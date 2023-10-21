package appinfo

import (
	"context"
	"pheet-fiber-backend/models"
	"sync"
)

type AppInfoUsecase interface {
	FindCategory(ctx context.Context, args *sync.Map) ([]*models.Categories, error)
	InsertCategories(ctx context.Context, cats []*models.Categories) error
	DeleteCategory(ctx context.Context, id int) error
}
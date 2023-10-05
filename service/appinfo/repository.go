package appinfo

import (
	"context"
	"pheet-fiber-backend/models"
	"sync"
)

/* pod interface */
type AppInfoRepository interface {
	FindCategory(ctx context.Context, args *sync.Map) ([]*models.Catagory, error)
	InsertCategories(ctx context.Context, cats []*models.Catagory) error
}
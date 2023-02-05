package product

import (
	"context"
	"pheet-fiber-backend/models"
)

// *สร้าง pod interface กำหนดว่ามีกี่* //
type ProductRepository interface {
	FetchAll(ctx context.Context)([]*models.Products, error)
	FetchById(id int)(*models.Products, error)
	FetchByType(coffType string)([]*models.Products, error)
	FetchUser(username string)(*models.User, error)
	Create(ctx context.Context, product *models.Products) error
	SignUp(ctx context.Context, user *models.User) error
	Update(product *models.Products) error
	Delete(id int)error
}
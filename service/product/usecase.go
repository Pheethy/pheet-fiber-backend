package product

import (
	"context"
	"pheet-fiber-backend/models"

	"github.com/gofrs/uuid"
)

//* สร้าง pod interface กำหนดว่ามี service อะไรให้ใช้บ้าง*//
type ProductUsecase interface{
	GetProducts(ctx context.Context)([]*models.Products,error)
	GetProduct(ictx context.Context, id int)(*models.Products, error)
	GetProductByType(coffType string)([]*models.Products, error)
	GetUser(username string)(*models.User, error)
	Create(ctx context.Context, product *models.Products) error
	SignUp(ctx context.Context, user *models.User) error
	Update(ctx context.Context, product *models.Products, id *uuid.UUID) error
	Delete(id int)error
}
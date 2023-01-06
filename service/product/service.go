package product

import (
	"context"
	"pheet-fiber-backend/models"
)

//* สร้าง pod interface กำหนดว่ามี service อะไรให้ใช้บ้าง*//
type ProductUsecase interface{
	GetProducts(ctx context.Context)([]*models.Product,error)
	GetProduct(id int)(*models.Product, error)
	GetProductByType(coffType string)([]*models.Product, error)
	GetUser(username string)(*models.User, error)
	Create(product *models.Product)error
	SignUp(user *models.SignUpReq)error
	Update(product *models.Product) error
	Delete(id int)error
}
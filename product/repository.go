package product

import "main/models"

// *สร้าง pod interface กำหนดว่ามีกี่* //
type ProductRepository interface {
	FetchAll()([]*models.Product, error)
	FetchById(id int)(*models.Product, error)
	FetchByType(coffType string)([]*models.Product, error)
	FetchUser(username string)(*models.User, error)
	Create(product *models.Product)error
	SignUp(user *models.SignUpReq)error
	Update(product *models.Product) error
	Delete(id int)error
}
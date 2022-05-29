package repository

import "main/models"

// *สร้าง pod interface กำหนดว่ามีกี่* //
type ProductRepository interface {
	FetchAll()([]*models.Product, error)
	FetchById(id int)(*models.Product, error)
	Create(product *models.Product)error
	Update(product *models.Product)error
	Delete(id int)error
}
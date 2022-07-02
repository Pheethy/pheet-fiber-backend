package service

import "main/models"

//* สร้าง pod interface กำหนดว่ามี service อะไรให้ใช้บ้าง*//
type ProductService interface{
	GetProducts()([]*models.Product,error)
	GetProduct(id int)(*models.Product, error)
	GetProductByType(coffType string)([]*models.Product, error)
	Create(product *models.Product)error
	Update(product *models.Product)error
	Delete(id int)error
}
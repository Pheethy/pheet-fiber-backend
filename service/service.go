package service

import (
	"main/models"
	"main/repository"
)

// *Adapter* //
type productService struct{
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository)productService{
	return productService{productRepo: productRepo}
}

func (r productService)GetProducts()([]*models.Product,error){
	return r.productRepo.FetchAll()	
}

func (r productService)GetProduct(id int)(*models.Product, error){
	return r.productRepo.FetchById(id)
}

func (r productService)GetUser(username string)(*models.User, error){
	return r.productRepo.FetchUser(username)
}

func (r productService)GetProductByType(coffType string)([]*models.Product, error){
	return r.productRepo.FetchByType(coffType)
}

func (r productService)Create(product *models.Product)error{
	return r.productRepo.Create(product)
}

func (r productService)SignUp(user *models.SignUpReq)error{
	return r.productRepo.SignUp(user)
}

func(r productService)Update(product *models.Product)error{
	return r.productRepo.Update(product)
}

func(r productService)Delete(id int)error{
	return r.productRepo.Delete(id)
}
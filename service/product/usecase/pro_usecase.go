package service

import (
	"context"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/product"
)

// *Adapter* //
type productUsecase struct{
	productRepo product.ProductRepository
}

func NewProductUsecase(productRepo product.ProductRepository) product.ProductUsecase { 
	return productUsecase{productRepo: productRepo}
}

func (r productUsecase)GetProducts(ctx context.Context)([]*models.Product,error){
	return r.productRepo.FetchAll(ctx)	
}

func (r productUsecase)GetProduct(id int)(*models.Product, error){
	return r.productRepo.FetchById(id)
}

func (r productUsecase)GetUser(username string)(*models.User, error){
	return r.productRepo.FetchUser(username)
}

func (r productUsecase)GetProductByType(coffType string)([]*models.Product, error){
	return r.productRepo.FetchByType(coffType)
}

func (r productUsecase)Create(product *models.Product)error{
	return r.productRepo.Create(product)
}

func (r productUsecase)SignUp(user *models.SignUpReq)error{
	return r.productRepo.SignUp(user)
}

func(r productUsecase)Update(product *models.Product) error{
	return r.productRepo.Update(product)
}

func(r productUsecase)Delete(id int)error{
	return r.productRepo.Delete(id)
}
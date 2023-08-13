package service

import (
	"context"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/product"

	"github.com/gofrs/uuid"
)

// *Adapter* //
type productUsecase struct{
	productRepo product.ProductRepository
}

func NewProductUsecase(productRepo product.ProductRepository) product.ProductUsecase { 
	return productUsecase{productRepo: productRepo}
}

func (r productUsecase)GetProducts(ctx context.Context)([]*models.Products,error){
	return r.productRepo.FetchAll(ctx)	
}

func (r productUsecase)GetProduct(ctx context.Context, id int)(*models.Products, error){
	return r.productRepo.FetchById(ctx, id)
}

func (r productUsecase)GetUser(username string)(*models.User, error){
	return r.productRepo.FetchUser(username)
}

func (r productUsecase)GetProductByType(coffType string)([]*models.Products, error){
	return r.productRepo.FetchByType(coffType)
}

func (r productUsecase)Create(ctx context.Context, product *models.Products) error{
	return r.productRepo.Create(ctx, product)
}

func (r productUsecase)SignUp(ctx context.Context, user *models.User) error {
	return r.productRepo.SignUp(ctx, user)
}

func(r productUsecase)Update(ctx context.Context, product *models.Products, id *uuid.UUID) error{
	return r.productRepo.Update(ctx, product, id)
}

func(r productUsecase)Delete(id int)error{
	return r.productRepo.Delete(id)
}
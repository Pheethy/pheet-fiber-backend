package usecase

import (
	"context"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/file"
	"pheet-fiber-backend/service/product"
)

type productUsecase struct {
	proRepo product.IProductRepository
	fileUs file.IFileUsecase
}

func NewProductUsecase(proRepo product.IProductRepository, fileUs file.IFileUsecase) product.IProductUsecase {
	return productUsecase{
		proRepo: proRepo,
		fileUs: fileUs,
	}
}

func (u productUsecase) FetchOneProduct(ctx context.Context, id string) (*models.Products, error) {
	return u.proRepo.FetchOneProduct(ctx, id)
}
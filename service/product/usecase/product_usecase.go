package usecase

import (
	"context"
	"pheet-fiber-backend/helper"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/file"
	"pheet-fiber-backend/service/product"
	"sync"
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

func (u productUsecase) FetchAllProduct(ctx context.Context, args *sync.Map, paginate *helper.Paginator) ([]*models.Products, error) {
	return u.proRepo.FetchAllProduct(ctx, args, paginate)
}
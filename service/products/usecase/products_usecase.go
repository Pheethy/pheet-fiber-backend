package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"path/filepath"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/file"
	"pheet-fiber-backend/service/products"
	"pheet-fiber-backend/service/utils"
	"strings"
	"sync"

	"github.com/Pheethy/psql/helper"
	"github.com/gofrs/uuid"
)

/* Adapter */
type productsUsecase struct {
	productRepo products.IProductsRepositoryDB
	fileUs      file.IFileUsecase
	cfg         config.Iconfig
}

/* Constructor */
func NewProductsUsecase(productRepo products.IProductsRepositoryDB, fileUs file.IFileUsecase, cfg config.Iconfig) products.IProductUsecase {
	return productsUsecase{
		productRepo: productRepo,
		fileUs:      fileUs,
		cfg:         cfg,
	}
}

func (p productsUsecase) FetchAllProducts(ctx context.Context, args *sync.Map, paginator *helper.Paginator) ([]*models.Products, error) {
	return p.productRepo.FetchAllProducts(ctx, args, paginator)
}

func (p productsUsecase) FetchOneProduct(ctx context.Context, productId *uuid.UUID) (*models.Products, error) {
	return p.productRepo.FetchOneProduct(ctx, productId)
}

func (p productsUsecase) CraeteProduct(ctx context.Context, req *models.Products, files []*multipart.FileHeader) error {
	if len(files) > 0 {
		reqFile := make([]*models.FileReq, 0)
		for _, file := range files {
			ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
			if ok := p.validateFileType(ext); !ok {
				return errors.New("file type is invalid")
			}

			if file.Size > int64(p.cfg.App().FileLimit()) {
				return fmt.Errorf("file size must less than %d MiB", int(math.Ceil(float64(p.cfg.App().FileLimit())/math.Pow(1024, 2))))
			}

			filename := utils.RandFileName(ext)
			reqFile = append(reqFile, &models.FileReq{
				File:        file,
				Destination: constants.PRODUCT_IMAGE_DESTINETION + "/" + filename,
				Extension:   ext,
				FileName:    file.Filename,
			})
		}

		/* upload images to google cloud platfrom */
		newFileInfo, err := p.fileUs.UploadToGCP(ctx, reqFile)
		if err != nil {
			return fmt.Errorf("upload product image failed: %v", err.Error())
		}

		images := make([]*models.Image, 0)
		for index := range newFileInfo {
			image := &models.Image{
				FileName:  newFileInfo[index].FileName,
				URL:       newFileInfo[index].Url,
				ProductId: req.Id,
			}
			image.NewId()
			image.SetCreatedAt()
			image.SetUpdatedAt()
			images = append(images, image)
		}

		req.Images = images
	}
	return p.productRepo.CraeteProduct(ctx, req)
}

func (p productsUsecase) UpdateProduct(ctx context.Context, product *models.Products, files []*multipart.FileHeader) error {
	if len(files) > 0 {
		reqFile := make([]*models.FileReq, 0)
		for _, file := range files {
			ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
			if ok := p.validateFileType(ext); !ok {
				return errors.New("file type is invalid")
			}

			if file.Size > int64(p.cfg.App().FileLimit()) {
				return fmt.Errorf("file size must less than %d MiB", int(math.Ceil(float64(p.cfg.App().FileLimit())/math.Pow(1024, 2))))
			}

			filename := utils.RandFileName(ext)
			reqFile = append(reqFile, &models.FileReq{
				File:        file,
				Destination: constants.PRODUCT_IMAGE_DESTINETION + "/" + filename,
				Extension:   ext,
				FileName:    file.Filename,
			})
		}

		/* upload images to google cloud platfrom */
		newFileInfo, err := p.fileUs.UploadToGCP(ctx, reqFile)
		if err != nil {
			return fmt.Errorf("upload product image failed: %v", err.Error())
		}

		images := make([]*models.Image, 0)
		for index := range newFileInfo {
			image := &models.Image{
				FileName:  newFileInfo[index].FileName,
				URL:       newFileInfo[index].Url,
				ProductId: product.Id,
			}
			image.NewId()
			image.SetCreatedAt()
			image.SetUpdatedAt()
			images = append(images, image)
		}

		product.Images = images
	}
	return p.productRepo.UpdateProduct(ctx, product)
}

func (p productsUsecase) DeleteProduct(ctx context.Context, product *models.Products) error {
	return p.productRepo.DeleteProduct(ctx, product)
}

func (p productsUsecase) DeleteImages(ctx context.Context, ids []*uuid.UUID) error {
	return p.productRepo.DeleteImages(ctx, ids)
}

func (p productsUsecase) validateFileType(ext string) bool {
	if ext == "" {
		return false
	}

	expMap := []string{"png", "jpg", "jpeg"}
	for index := range expMap {
		if expMap[index] == ext {
			return true
		}
	}
	return false
}

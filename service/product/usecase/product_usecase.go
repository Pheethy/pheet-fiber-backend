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
	"pheet-fiber-backend/helper"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/file"
	"pheet-fiber-backend/service/product"
	"pheet-fiber-backend/service/utils"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
)

type productUsecase struct {
	proRepo product.IProductRepository
	fileUs  file.IFileUsecase
	cfg     config.Iconfig
}

func NewProductUsecase(proRepo product.IProductRepository, fileUs file.IFileUsecase, cfg config.Iconfig) product.IProductUsecase {
	return productUsecase{
		proRepo: proRepo,
		fileUs:  fileUs,
		cfg:     cfg,
	}
}

func (u productUsecase) FetchOneProduct(ctx context.Context, id string) (*models.Products, error) {
	return u.proRepo.FetchOneProduct(ctx, id)
}

func (u productUsecase) FetchAllProduct(ctx context.Context, args *sync.Map, paginate *helper.Paginator) ([]*models.Products, error) {
	return u.proRepo.FetchAllProduct(ctx, args, paginate)
}

func (u productUsecase) CraeteProduct(ctx context.Context, req *models.Products, files []*multipart.FileHeader) error {
	if len(files) > 0 {
		var reqFile = make([]*models.FileReq, 0)
		for _, file := range files {
			ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
			if ok := u.validateFileType(ext); !ok {
				return errors.New("file type is invalid")
			}

			if file.Size > int64(u.cfg.App().FileLimit()) {
				return fmt.Errorf("file size must less than %d MiB", int(math.Ceil(float64(u.cfg.App().FileLimit())/math.Pow(1024, 2))))
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
		newFileInfo, err := u.fileUs.UploadToGCP(reqFile)
		if err != nil {
			return fmt.Errorf("upload product image failed: %v", err.Error())
		}

		var images = make([]*models.Image, 0)
		for index := range newFileInfo {
			image := &models.Image{
				FilenName: newFileInfo[index].FileName,
				Url:       newFileInfo[index].Url,
				ProductId: req.ID,
			}
			image.NewId()
			image.SetCreatedAt()
			image.SetUpdatedAt()
			images = append(images, image)
		}

		req.Images = images
	}
	return u.proRepo.CraeteProduct(ctx, req)
}

func (u productUsecase) UpdateProduct(ctx context.Context, product *models.Products, files []*multipart.FileHeader) error {
	if len(files) > 0 {
		var reqFile = make([]*models.FileReq, 0)
		for _, file := range files {
			ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
			if ok := u.validateFileType(ext); !ok {
				return errors.New("file type is invalid")
			}

			if file.Size > int64(u.cfg.App().FileLimit()) {
				return fmt.Errorf("file size must less than %d MiB", int(math.Ceil(float64(u.cfg.App().FileLimit())/math.Pow(1024, 2))))
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
		newFileInfo, err := u.fileUs.UploadToGCP(reqFile)
		if err != nil {
			return fmt.Errorf("upload product image failed: %v", err.Error())
		}

		var images = make([]*models.Image, 0)
		for index := range newFileInfo {
			image := &models.Image{
				FilenName: newFileInfo[index].FileName,
				Url:       newFileInfo[index].Url,
				ProductId: product.ID,
			}
			image.NewId()
			image.SetCreatedAt()
			image.SetUpdatedAt()
			images = append(images, image)
		}

		product.Images = images
	}
	return u.proRepo.UpdateProduct(ctx, product)
}

func (u productUsecase) DeleteProduct(ctx context.Context, productId string) error {
	return u.proRepo.DeleteProduct(ctx, productId)
}

func (u productUsecase) DeleteImages(ctx context.Context, ids []*uuid.UUID) error {
	return u.proRepo.DeleteImages(ctx, ids)
}

func (u productUsecase) validateFileType(ext string) bool {
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

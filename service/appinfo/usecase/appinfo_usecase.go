package usecase

import (
	"context"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/appinfo"
	"sync"
)

type appInfoUsecase struct {
	cfg      config.Iconfig
	infoRepo appinfo.AppInfoRepository
}

func NewAppInfoUsecase(cfg config.Iconfig, infoRepo appinfo.AppInfoRepository) appinfo.AppInfoUsecase {
	return appInfoUsecase{
		cfg: cfg,
		infoRepo: infoRepo,
	}
}


func (u appInfoUsecase) FindCategory(ctx context.Context, args *sync.Map) ([]*models.Catagory, error) {
	return u.infoRepo.FindCategory(ctx, args)
}

func (u appInfoUsecase) InsertCategories(ctx context.Context, cats []*models.Catagory) error {
	return u.infoRepo.InsertCategories(ctx, cats)
}

func (u appInfoUsecase) DeleteCategory(ctx context.Context, id int) error {
	return u.infoRepo.DeleteCategory(ctx, id)
}
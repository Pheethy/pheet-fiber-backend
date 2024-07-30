package file

import (
	"context"
	"pheet-fiber-backend/models"
)

type IFileUsecase interface {
	UploadToGCP(ctx context.Context, fileReq []*models.FileReq) ([]*models.FileResp, error)
	DeleteOnGCP(req []*models.DeleteFileReq) error
}

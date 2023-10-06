package file

import "pheet-fiber-backend/models"

type IFileUsecase interface {
	UploadToGCP(fileReq []*models.FileReq) ([]*models.FileResp, error)
}
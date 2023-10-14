package handler

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/file"
	"pheet-fiber-backend/service/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type fileHandler struct {
	cfg config.Iconfig
	fileUs file.IFileUsecase
}

func NewFileHandler(cfg config.Iconfig, fileUs file.IFileUsecase) file.IFileHandler {
	return &fileHandler{
		cfg: cfg,
		fileUs: fileUs,
	}
}

func (f *fileHandler) UploadFile(c *fiber.Ctx) error {
	var req = make([]*models.FileReq, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Cast Form Failed.")
	}
	/* ทำการรับ Files จาก Form */
	files := form.File["files"]
	dest := c.FormValue("destination")


	for _, file := range files {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if ok := f.validateFileType(ext); !ok {
			return fiber.NewError(http.StatusBadRequest, "file type is invalid.")
		}

		if file.Size > int64(f.cfg.App().FileLimit()) {
			return fiber.NewError(http.StatusBadRequest, fmt.Sprintf("file size must less than %d MiB", int(math.Ceil(float64(f.cfg.App().FileLimit())/math.Pow(1024, 2)))))
		}

		filename := utils.RandFileName(ext)
		req = append(req, &models.FileReq{
			File: file,
			Destination: dest + "/" + filename,
			Extension: ext,
			FileName: file.Filename,
		})
	}
	
	newFileInfo, err := f.fileUs.UploadToGCP(req)
	if err != nil {
		log.Println("incase")
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"message": "uploaded.",
		"resp": newFileInfo,
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (f *fileHandler) DeleteFile(c *fiber.Ctx) error {
	var req = make([]*models.DeleteFileReq, 0)
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if err := f.fileUs.DeleteOnGCP(req); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"message": "successful.",
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (f *fileHandler) validateFileType(ext string) bool {
	if ext == "" {
		return false
	}
	
	expMap := []string{"png","jpg","jpeg"}
	for index := range expMap {
		if expMap[index] == ext {
			return true
		}
	}
	return false
}
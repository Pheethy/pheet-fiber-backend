package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/file"
	"time"

	"cloud.google.com/go/storage"
)

type fileUsecase struct {
	cfg config.Iconfig
}

func NewFileUsecase(cfg config.Iconfig) file.IFileUsecase {
	return &fileUsecase{cfg: cfg}
}

func (f fileUsecase) UploadToGCP(fileReq []*models.FileReq) ([]*models.FileResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Err new GCP client: %v", err)
	}
	defer client.Close()

	/* ทำ pool worker */
	/* สร้าง Channel Jobs */
	jobsCh := make(chan *models.FileReq, len(fileReq))
	/* สร้าง Channel Result */
	resultCh := make(chan *models.FileResp, len(fileReq))
	/* สร้าง Channel Error */
	errCh := make(chan error, len(fileReq))
	/* สร้าง Entity สำหรับ Response Request นี้ */
	resp := make([]*models.FileResp, 0)

	/* ทำการนำ file request ใส่ไปใน jobs channel */
	for _, r := range fileReq {
		jobsCh <- r
	}
	close(jobsCh)

	/* ประกาศจำนวน worker */
	var workers int = 5
	/* สร้าง loop สำหรับการทำงาน upload file */
	for i := 0; i < workers; i++ {
		//working zone
		go f.streamFileUpload(ctx, client, jobsCh, resultCh, errCh)
	}

	/* สร้าง loop สำหรับการ Response */
	for a := 0; a < len(fileReq); a++ {
		//handler err โดยการรับค่า err จาก Channel errCh
		if err := <-errCh; err != nil {
			return nil, fmt.Errorf("Response err: %v", err)
		}
		//ทำการนำ result จาก resultCh ใส่ใน resp
		result := <-resultCh
		resp = append(resp, result)
	}

	return resp, nil
}

// เอา example มาจาก https://cloud.google.com/storage/docs/uploading-objects-from-memory
func (f fileUsecase) streamFileUpload(ctx context.Context, client *storage.Client, jobs <-chan *models.FileReq, result chan<- *models.FileResp, errs chan<- error) {
	/* concept upload แปลง file -> []byte -> buffer -> upload */

	/* recap เราเอา fileReq loop เข่้า jobsCh ที่ละตัว ก็ต้องเอา ออกทีละตัว มาใส่ jobs  ทำให้ต้อง range เอา job อีกรอบ */
	for job := range jobs {
		/* แปลง File *multipart.FileHeader -> multipart.File*/
		container, err := job.File.Open()
		if err != nil {
			errs <- err
			return
		}

		/* ทำการแปลง multipart.File -> []byte */
		byt, err := ioutil.ReadAll(container)
		if err != nil {
			errs <- err
			return
		}
		/* ทำการแปลง []byte -> *byte.Buffer */
		buff := bytes.NewBuffer(byt)

		// Upload an object with storage.Writer.
		wc := client.Bucket(f.cfg.App().GCPBucket()).Object(job.Destination).NewWriter(ctx)
		wc.ChunkSize = 0 // note retries are not supported for chunk size 0.

		if _, err = io.Copy(wc, buff); err != nil {
			errs <- fmt.Errorf("io.Copy: %w", err)
			return
		}
		// Data can continue to be added to the file until the writer is closed.
		if err := wc.Close(); err != nil {
			errs <- fmt.Errorf("Writer.Close: %w", err)
			return
		}
		fmt.Printf("👽 %v uploaded to %v.\n", job.FileName, job.Destination)

		newFile := &models.FilePub{
			File: &models.FileResp{
				FileName: job.FileName,
				Url: fmt.Sprintf("https://storage.googleapis.com/%s/%s", f.cfg.App().GCPBucket(), job.Destination),
			},
			Bucket: f.cfg.App().GCPBucket(),
			Destination: job.Destination,
		}

		if err := newFile.MakePublic(ctx, client); err != nil {
			errs <- err
			return
		}

		/* กรณีไม่มี error ก็ต้องทำการ return ค่า nil ออกไป errCh เพราะเราประกาศรับค่าไว้ */
		errs <- nil
		result <- newFile.File
	}
}


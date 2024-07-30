package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/file"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/opentracing/opentracing-go"
)

type fileUsecase struct {
	cfg config.Iconfig
}

func NewFileUsecase(cfg config.Iconfig) file.IFileUsecase {
	return &fileUsecase{cfg: cfg}
}

func (f fileUsecase) UploadToGCP(ctx context.Context, fileReq []*models.FileReq) ([]*models.FileResp, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UploadToGCP")
	defer span.Finish()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("err new GCP client: %v", err)
	}
	defer client.Close()

	jobsCh := make(chan *models.FileReq, len(fileReq))
	resultCh := make(chan *models.FileResp, len(fileReq))
	errCh := make(chan error, len(fileReq))
	resp := make([]*models.FileResp, 0)

	for _, r := range fileReq {
		jobsCh <- r
	}
	close(jobsCh)

	workers := 5
	wg := new(sync.WaitGroup)
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			f.streamFileUpload(ctx, client, jobsCh, resultCh, errCh)
		}()
	}
	wg.Wait()
	close(errCh)
	close(resultCh)

	for a := 0; a < len(fileReq); a++ {
		if err := <-errCh; err != nil {
			return nil, errors.New(err.Error())
		}
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
		fmt.Printf("🍫 %v uploaded to %v.\n", job.FileName, job.Destination)

		newFile := &models.FilePub{
			File: &models.FileResp{
				FileName: job.FileName,
				Url:      fmt.Sprintf("https://storage.googleapis.com/%s/%s", f.cfg.App().GCPBucket(), job.Destination),
			},
			Bucket:      f.cfg.App().GCPBucket(),
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

func (f fileUsecase) DeleteOnGCP(req []*models.DeleteFileReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("err new GCP client: %v", err)
	}
	defer client.Close()

	jobsCh := make(chan *models.DeleteFileReq, len(req))
	errCh := make(chan error, len(req))

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	workers := 5
	wg := new(sync.WaitGroup)
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			f.deleteFile(ctx, client, jobsCh, errCh)
		}()
	}

	wg.Wait()
	close(errCh)

	for i := 0; i < len(req); i++ {
		if err := <-errCh; err != nil {
			return errors.New(err.Error())
		}
	}

	return nil
}

// deleteFile removes specified object.
func (f fileUsecase) deleteFile(ctx context.Context, client *storage.Client, jobs <-chan *models.DeleteFileReq, errs chan<- error) {
	for job := range jobs {
		o := client.Bucket(f.cfg.App().GCPBucket()).Object(job.Destination)

		attrs, err := o.Attrs(ctx)
		if err != nil {
			if ok := strings.Contains(err.Error(), "object doesn't exist"); ok {
				errs <- fmt.Errorf("object.Attrs: %w", errors.New("can't found image"))
				return
			}
			errs <- fmt.Errorf("object.Attrs: %w", err)
			return
		}
		o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

		if err := o.Delete(ctx); err != nil {
			errs <- fmt.Errorf("Object(%q).Delete: %w", job.Destination, err)
			return
		}
		fmt.Printf("Blob %v deleted.\n", job.Destination)

		/* กรณีไม่มี error ก็ต้องทำการ return ค่า nil ออกไป errCh เพราะเราประกาศรับค่าไว้ */
		errs <- nil
	}
}

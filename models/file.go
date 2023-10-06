package models

import (
	"context"
	"fmt"
	"mime/multipart"

	"cloud.google.com/go/storage"
)

type FileReq struct {
	File        *multipart.FileHeader `form:"file"`
	Destination string                `form:"destination"` /* ที่อยู่ File (ปลายทาง)*/
	Extension   string                `form:"extension"`   /* นามสกุล File */
	FileName    string                `form:"file_name"`   /* ชื่อ File */
}

type FileResp struct {
	FileName string `json:"file_name"` /* ชื่อ File */
	Url      string `json:"url"`       /* url ของรูปภาพ */
}

type FilePub struct {
	Bucket string
	Destination string
	File *FileResp
}

type DeleteFileReq struct {
	Destination string `json:"destination"` /* ที่อยู่สำหรับลบ File */
}

// makePublic gives all users read access to an object.
func (f *FilePub) MakePublic(ctx context.Context, client *storage.Client) error {
	acl := client.Bucket(f.Bucket).Object(f.Destination).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return fmt.Errorf("ACLHandle.Set: %w", err)
	}
	fmt.Printf("Blob %v is now publicly accessible.\n", f.Destination)
	return nil
}

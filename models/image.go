package models

import (
	"pheet-fiber-backend/helper"

	"github.com/gofrs/uuid"
)

type Image struct {
	TableName struct{}          `json:"-" db:"images" pk:"ID"`
	ID        *uuid.UUID        `json:"id" db:"id" type:"uuid"`
	FilenName string            `json:"filename" db:"filename" type:"string"`
	Url       string            `json:"url" db:"url" type:"string"`
	ProductId string            `json:"product_id" db:"product_id" type:"string"`
	CreatedAt *helper.Timestamp `db:"created_at" json:"created_at" type:"timestamp"`
	UpdatedAt *helper.Timestamp `db:"updated_at" json:"updated_at" type:"timestamp"`
}

package models

import (
	"github.com/Pheethy/psql/helper"
	"time"

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

func (i *Image) NewId() {
	uuid, _ := uuid.NewV4()
	i.ID = &uuid
}

func (i *Image) SetCreatedAt() {
	time := helper.NewTimestampFromTime(time.Now())
	i.CreatedAt = &time
}

func (i *Image) SetUpdatedAt() {
	time := helper.NewTimestampFromTime(time.Now())
	i.UpdatedAt = &time
}
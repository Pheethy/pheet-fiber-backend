package models

import (
	"pheet-fiber-backend/helper"

	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
)

// *Entity เพื่อจะส่งข้อมูลออกไป *//
type Products struct {
	Id        *uuid.UUID `db:"id" json:"id" type:"uuid"`
	Name      string     `db:"name" json:"name" type:"string" validate:"required"`
	Detail    string     `db:"detail" json:"detail" type:"string"`
	Type      string     `db:"type" json:"type" type:"string"`
	Price     int64      `db:"price" json:"price" type:"int64"`
	Image     string     `db:"cover" json:"image" type:"string"`
	CreatedAt *helper.Timestamp `db:"created_at" json:"created_at" type:"timestamp"`
	UpdatedAt *helper.Timestamp `db:"updated_at" json:"updated_at" type:"timestamp"`
}

type Element struct {
	FailedField string
	Tag         string
	Value       string
}

func (p *Products) NewUUID() {
	id, _ := uuid.NewV4()
	p.Id = &id
}

func (p *Products) SetCreatedAt(time *helper.Timestamp) {
	p.CreatedAt = time
}

func (p *Products) SetUpdatedAt(time *helper.Timestamp) {
	p.UpdatedAt = time
}

func (p *Products) ValidationStruct() error {
	var validate = validator.New()
	err := validate.Struct(p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Products) SetImagePath(path string) {
	p.Image = path
}

func (p *Products) MergeProduct(existPro *Products) {
	switch {
	case p.Name == "":
		p.Name = existPro.Name
		fallthrough
	case p.Detail == "":
		p.Detail = existPro.Detail
		fallthrough
	case p.Price == 0:
		p.Price = existPro.Price
		fallthrough
	case p.Type == "":
		p.Type = existPro.Type
	}
}
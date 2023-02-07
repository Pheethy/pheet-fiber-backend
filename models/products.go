package models

import (
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
	CreatedAt *Timestamp `db:"created_at" json:"created_at" type:"timestamp"`
	UpdatedAt *Timestamp `db:"updated_at" json:"updated_at" type:"timestamp"`
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

func (p *Products) SetCreatedAt(time *Timestamp) {
	p.CreatedAt = time
}

func (p *Products) SetUpdatedAt(time *Timestamp) {
	p.UpdatedAt = time
}

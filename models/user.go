package models

import (
	"pheet-fiber-backend/helper"

	"github.com/gofrs/uuid"
)

type User struct {
	Id *uuid.UUID `db:"id" json:"id" type:"uuid"`
	Username string `db:"username" json:"username" type:"string"`
	Password string `db:"password" json:"password" type:"string"`
	CreatedAt *helper.Timestamp `db:"created_at" json:"created_at" type:"timestamp"`
	UpdatedAt *helper.Timestamp `db:"updated_at" json:"updated_at" type:"timestamp"`
}

func (p *User) NewUUID(){
	id, _ := uuid.NewV4()
	p.Id = &id
}

func (p *User) SetCreatedAt(time *helper.Timestamp) {
	p.CreatedAt = time
}

func (p *User) SetUpdatedAt(time *helper.Timestamp) {
	p.UpdatedAt = time
}
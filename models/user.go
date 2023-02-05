package models

import "github.com/gofrs/uuid"

type User struct {
	Id *uuid.UUID `db:"id" json:"id" type:"uuid"`
	Username string `db:"username" json:"username" type:"string"`
	Password string `db:"password" json:"password" type:"string"`
	CreatedAt *Timestamp `db:"created_at" json:"created_at" type:"timestamp"`
	UpdatedAt *Timestamp `db:"updated_at" json:"updated_at" type:"timestamp"`
}

func (p *User) NewUUID(){
	id, _ := uuid.NewV4()
	p.Id = &id
}

func (p *User) SetCreatedAt(time *Timestamp) {
	p.CreatedAt = time
}

func (p *User) SetUpdatedAt(time *Timestamp) {
	p.UpdatedAt = time
}
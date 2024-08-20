package models

import (
	"time"

	"github.com/Pheethy/psql/helper"
	"github.com/gofrs/uuid"
)

type OAuth struct {
	TableName    struct{}          `json:"-" db:"oauth" pk:"Id"`
	Id           *uuid.UUID        `json:"id" db:"id" type:"uuid"`
	UserId       *uuid.UUID        `json:"user_id" db:"user_id" type:"uuid"`
	AccessToken  string            `json:"access_token" db:"access_token" type:"string"`
	RefreshToken string            `json:"refresh_token" db:"refresh_token" type:"string"`
	CreatedAt    *helper.Timestamp `json:"created_at" db:"created_at" type:"timestamp"`
	UpdatedAt    *helper.Timestamp `json:"updated_at" db:"updated_at" type:"timestamp"`
}

func (u *OAuth) NewId() {
	id, _ := uuid.NewV4()
	u.Id = &id
}

func (u *OAuth) SetCreatedAt() {
	ti := helper.NewTimestampFromTime(time.Now())
	u.CreatedAt = &ti
}

func (u *OAuth) SetUpdatedAt() {
	ti := helper.NewTimestampFromTime(time.Now())
	u.UpdatedAt = &ti
}

func (u *OAuth) SetToken(userPass *UserPassport) {
	u.UserId = userPass.User.Id
	u.AccessToken = userPass.Token.AccessToken
	u.RefreshToken = userPass.Token.RefreshToken
}


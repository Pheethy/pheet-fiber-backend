package models

import (
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/Pheethy/psql/helper"
	"github.com/gofrs/uuid"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	TableName struct{}          `json:"-" db:"users" pk:"Id"`
	Id        *uuid.UUID        `json:"id" db:"id" type:"uuid"`
	Username  string            `json:"username" form:"username" db:"username" type:"string"`
	Password  string            `json:"password" form:"password" db:"password" type:"string"`
	Email     string            `json:"email" form:"email" db:"email" type:"string"`
	RoleId    int64             `json:"role_id" db:"role_id" type:"int64"`
	RoleTitle string            `json:"role_title" db:"role_title" type:"string"`
	CreatedAt *helper.Timestamp `json:"created_at" db:"created_at" type:"timestamp"`
	UpdatedAt *helper.Timestamp `json:"updated_at" db:"updated_at" type:"timestamp"`
}

type Users []*User

func (u *User) NewId() {
	id, _ := uuid.NewV4()
	u.Id = &id
}

func (u *User) SetCreatedAt() {
	ti := helper.NewTimestampFromTime(time.Now())
	u.CreatedAt = &ti
}

func (u *User) SetUpdatedAt() {
	ti := helper.NewTimestampFromTime(time.Now())
	u.UpdatedAt = &ti
}

func (u *User) BcryptHashing() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return fmt.Errorf("hashing password failed: %s", err.Error())
	}
	u.Password = string(hash)
	return nil
}

func (u *User) ComparePassword(i *User) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(i.Password)); err != nil {
		return false
	}
	return true
}

func (u *User) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w\-.]+@([\w\-]+\.)+[\w\-]{2,4}$`, u.Email)
	if err != nil {
		return false
	}
	return match
}

type UserPassport struct {
	User  *User      `json:"user"`
	Token *UserToken `json:"token"`
}

type UserToken struct {
	Id           *uuid.UUID `json:"id" db:"id"`
	AccessToken  string     `json:"access_token" db:"access_token"`
	RefreshToken string     `json:"refresh_token" db:"refresh_token"`
}

type UserProfile struct {
	Id        *uuid.UUID        `json:"id"`
	Username  string            `json:"username" form:"username" db:"username" type:"string"`
	Email     string            `json:"email" form:"email" db:"email" type:"string"`
	RoleId    int64             `json:"role_id" db:"role_id" type:"int64"`
	RoleTitle string            `json:"role_title" db:"role_title" type:"string"`
	CreatedAt *helper.Timestamp `json:"created_at" db:"created_at" type:"timestamp"`
	UpdatedAt *helper.Timestamp `json:"updated_at" db:"updated_at" type:"timestamp"`
}

func NewUserProfileWithParams(params map[string]interface{}) *UserProfile {
	ptr := new(UserProfile)

	for key, val := range params {
		switch key {
		case "id":
			id := uuid.FromStringOrNil(val.(string))
			ptr.Id = &id
		case "username":
			ptr.Username = cast.ToString(val)
		case "email":
			ptr.Email = cast.ToString(val)
		case "role_id":
			ptr.RoleId = cast.ToInt64(val)
		case "role_title":
			ptr.RoleTitle = cast.ToString(val)
		case "created_at":
			if val != nil {
				if reflect.TypeOf(val).Kind() == reflect.String {
					timestamp := helper.NewTimestampFromString(val.(string))
					ptr.CreatedAt = &timestamp
				} else if reflect.TypeOf(val).String() == "time.Time" {
					timestamp := helper.NewTimestampFromTime(val.(time.Time))
					ptr.CreatedAt = &timestamp
				}
			}
		case "updated_at":
			if val != nil {
				if reflect.TypeOf(val).Kind() == reflect.String {
					timestamp := helper.NewTimestampFromString(val.(string))
					ptr.UpdatedAt = &timestamp
				} else if reflect.TypeOf(val).String() == "time.Time" {
					timestamp := helper.NewTimestampFromTime(val.(time.Time))
					ptr.UpdatedAt = &timestamp
				}
			}
		}
	}

	return ptr
}

type UserClaims struct {
	Id     *uuid.UUID `json:"id" db:"id"`
	RoleId int64      `json:"role_id" db:"role_id"`
}

type UserRefreshCredential struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

type UserRemoveCredential struct {
	OauthId *uuid.UUID `json:"oauth_id" db:"id" form:"oauth_id"`
}


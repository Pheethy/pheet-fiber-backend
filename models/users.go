package models

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Id       string `json:"id" db:"id" type:"string"`
	Email    string `json:"email" db:"email"`
	UserName string `json:"username" db:"username"`
	RoleId   int    `json:"role_id" db:"role_id"`
}

type UserRegisterReq struct {
	Email    string `json:"email" db:"email" form:"email"`
	Username string `json:"username" db:"username" form:"username"`
	Password string `json:"password" db:"password" form:"password"`
}

type UserCredential struct {
	Email    string `json:"email" db:"email" form:"email"`
	Password string `json:"password" db:"password" form:"password"`
}

type UserCredentialCheck struct {
	Id       string `db:"id"`
	Email    string `db:"email"`
	Username string `db:"username"`
	Password string `db:"password"`
	RoleId   int    `db:"role_id"`
}

func (u *UserRegisterReq) BcryptHashing() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return fmt.Errorf("hashing password failed: %w", err)
	}

	u.Password = string(hash)
	return nil
}

func (u *UserRegisterReq) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w\-.]+@([\w\-]+\.)+[\w\-]{2,4}$`, u.Email)
	if err != nil {
		return false
	}
	return match
}

type UserPassport struct {
	User  *Users     `json:"user"`
	Token *UserToken `json:"token"`
}

type UserToken struct {
	Id           string `json:"id" db:"id"`
	AccessToken  string `json:"access_token" db:"access_token"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
}

type UserClaims struct {
	Id string `json:"id" db:"id"`
	RoleId int `json:"role_id" db:"role_id"`
}

type UserRefreshCredential struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

type UserRemoveCredential struct {
	OauthId string `json:"oauth_id" db:"id" form:"oauth_id"`
}
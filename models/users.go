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
	Email    string `json:"email" db:"email"`
	UserName string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func (u *UserRegisterReq) BcryptHashing() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return fmt.Errorf("Hashing Password Failed: %w", err)
	}

	u.Password = string(hash)
	return nil
}

func (u *UserRegisterReq) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w-I.]+@(Ew-J+1.)+[\w-](2,4)$`, u.Password)
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

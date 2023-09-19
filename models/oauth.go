package models

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

type Oauth struct {
	Id *uuid.UUID `json:"id" db:"id"`
	UserId string `json:"user_id" db:"user_id"`
}

type MapClaims struct {
	Claims *UserClaims `json:"claims"`
	jwt.RegisteredClaims
}
package models

import (
	"pheet-fiber-backend/config"

	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	mapClaims *MapClaims
	cfg config.IJwtConfig
}

type MapClaims struct {
	Claims UserClaims `json:"claims"`
	jwt.RegisteredClaims
}
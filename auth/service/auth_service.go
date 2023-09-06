package service

import (
	"fmt"
	"math"
	"pheet-fiber-backend/auth"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type serviceAuth struct {
	cfg config.IJwtConfig
	mapClaims *MapClaims
}

type MapClaims struct {
	Claims *models.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

func NewAuthService(tokenType constants.TokenType, cfg config.IJwtConfig, claims *models.UserClaims) (auth.ServiceAuth, error) {
	switch tokenType {
	case constants.Access:
		return newAccessToken(cfg, claims), nil
	case constants.Refresh:
		return newRefreshToken(cfg, claims), nil
	default:
		return nil, fmt.Errorf("unkwon token type")
	}
}

func (a *serviceAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func jwtTimeDurationCal(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func newAccessToken(cfg config.IJwtConfig, claims *models.UserClaims) auth.ServiceAuth {
	return &serviceAuth{
		cfg: cfg,
		mapClaims: &MapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "pheet-shop-api",
				Subject:   "access-token",
				Audience:  []string{"cutomer", "admin"},
				ExpiresAt: jwtTimeDurationCal(int64(cfg.AccessExpiresAt())),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *models.UserClaims) auth.ServiceAuth {
	return &serviceAuth{
		cfg: cfg,
		mapClaims: &MapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "pheet-shop-api",
				Subject:   "refresh-token",
				Audience:  []string{"cutomer", "admin"},
				ExpiresAt: jwtTimeDurationCal(int64(cfg.RefreshExpiresAt())),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}
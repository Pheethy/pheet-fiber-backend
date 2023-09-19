package service

import (
	"errors"
	"fmt"
	"log"
	"math"
	"pheet-fiber-backend/auth"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type serviceAuth struct {
	cfg       config.IJwtConfig
	mapClaims *models.MapClaims
}

type adminAuth struct {
	*serviceAuth
}

func NewAuthService(tokenType constants.TokenType, cfg config.IJwtConfig, claims *models.UserClaims) (auth.ServiceAuth, error) {
	switch tokenType {
	case constants.Access:
		return newAccessToken(cfg, claims), nil
	case constants.Refresh:
		return newRefreshToken(cfg, claims), nil
	case constants.Admin:
		return newAdminToken(cfg), nil
	default:
		return nil, fmt.Errorf("unkwon token type")
	}
}

func (a *serviceAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func (a *adminAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*models.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid.")
		}
		return cfg.SecretKey(), nil
	})

	//handler error
	if err != nil {
		// checkFormToken && checkExpired
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid.")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token was expired.")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*models.MapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*models.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid.")
		}
		return cfg.AdminKey(), nil
	})

	//handler error
	if err != nil {
		// checkFormToken && checkExpired
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid.")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token was expired.")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*models.MapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func jwtTimeDurationCal(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func RepeatToken(cfg config.IJwtConfig, claims *models.UserClaims, exp int64) string {
	obj := &serviceAuth{
		cfg: cfg,
		mapClaims: &models.MapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "pheet-shop-api",
				Subject:   "refresh-token",
				Audience:  []string{"cutomer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}

	return obj.SignToken()
}

func newAccessToken(cfg config.IJwtConfig, claims *models.UserClaims) auth.ServiceAuth {
	log.Println("exAcc", cfg.AccessExpiresAt())
	return &serviceAuth{
		cfg: cfg,
		mapClaims: &models.MapClaims{
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
	log.Println("exRefresh", cfg.RefreshExpiresAt())
	return &serviceAuth{
		cfg: cfg,
		mapClaims: &models.MapClaims{
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

func newAdminToken(cfg config.IJwtConfig) auth.ServiceAuth {
	return &adminAuth{
		&serviceAuth{
			cfg: cfg,
			mapClaims: &models.MapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "pheet-shop-api",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}

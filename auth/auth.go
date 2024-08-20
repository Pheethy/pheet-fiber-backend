package auth

import (
	"errors"
	"fmt"
	"math"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IAuth interface {
	SignToken() string
}

type authAPIKey struct {
	mapClaims *mapClaims /* payload jwt */
	cfg       config.IJwtConfig
}

type authAdminKey struct {
	mapClaims *mapClaims /* payload jwt */
	cfg       config.IJwtConfig
}

type adpAuth struct {
	mapClaims *mapClaims /* payload jwt */
	cfg       config.IJwtConfig
}

type mapClaims struct {
	Claims *models.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

func NewAuth(tokenType constants.TokenType, cfg config.IJwtConfig, claims *models.UserClaims) (IAuth, error) {
	switch tokenType {
	case constants.Access:
		return newAccessToken(cfg, claims), nil
	case constants.Refresh:
		return newRefreshToken(cfg, claims), nil
	case constants.Admin:
		return newAdminToken(cfg), nil
	case constants.ApiKey:
		return newAPIKeyToken(cfg), nil

	default:
		return nil, errors.New(constants.ERROR_UNKNOW_TOKEN_TYPE)
	}
}

func (a *adpAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func (a *authAPIKey) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.ApiKey())
	return ss
}

func (a *authAdminKey) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.AdminKey())
	return ss
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*mapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &mapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(constants.ERROR_SIGNING_METHOD_IS_INVALID)
		}
		return cfg.SecretKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New(constants.ERROR_JWT_FORMAT_IS_INVALID)
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New(constants.ERROR_JWT_WAS_EXPIRED)
		}
		return nil, fmt.Errorf("parse token failed: %s", err.Error())
	}

	if claims, ok := token.Claims.(*mapClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New(constants.ERROR_CLAIMS_TYPE_IS_INVALID)
	}
}

func ParseTokenAdmin(cfg config.IJwtConfig, tokenString string) (*mapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &mapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(constants.ERROR_SIGNING_METHOD_IS_INVALID)
		}
		return cfg.AdminKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New(constants.ERROR_JWT_FORMAT_IS_INVALID)
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New(constants.ERROR_JWT_WAS_EXPIRED)
		}
		return nil, fmt.Errorf("parse token failed: %s", err.Error())
	}

	if claims, ok := token.Claims.(*mapClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New(constants.ERROR_CLAIMS_TYPE_IS_INVALID)
	}
}

func ParseAPIKey(cfg config.IJwtConfig, tokenString string) (*mapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &mapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(constants.ERROR_SIGNING_METHOD_IS_INVALID)
		}
		return cfg.ApiKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New(constants.ERROR_JWT_FORMAT_IS_INVALID)
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New(constants.ERROR_JWT_WAS_EXPIRED)
		}
		return nil, fmt.Errorf("parse token failed: %s", err.Error())
	}

	if claims, ok := token.Claims.(*mapClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New(constants.ERROR_CLAIMS_TYPE_IS_INVALID)
	}
}

func RepeatClaims(cfg config.IJwtConfig, claims *models.UserClaims, exp int) string {
	obj := &adpAuth{
		mapClaims: &mapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "flavorparser-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
		cfg: cfg,
	}

	return obj.SignToken()
}

func newAccessToken(cfg config.IJwtConfig, claims *models.UserClaims) IAuth {
	return &adpAuth{
		mapClaims: &mapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "flavorparser-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
		cfg: cfg,
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *models.UserClaims) IAuth {
	return &adpAuth{
		mapClaims: &mapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "flavorparser-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
		cfg: cfg,
	}
}

func newAdminToken(cfg config.IJwtConfig) IAuth {
	return &authAdminKey{
		mapClaims: &mapClaims{
			Claims: nil,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "flavorparser-api",
				Subject:   "admin-token",
				Audience:  []string{"admin"},
				ExpiresAt: jwtTimeDurationCal(300),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
		cfg: cfg,
	}
}

func newAPIKeyToken(cfg config.IJwtConfig) IAuth {
	return &authAPIKey{
		mapClaims: &mapClaims{
			Claims: nil,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "flavorparser-api",
				Subject:   "api-key",
				Audience:  []string{"admin", "customer"},
				ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(2, 0, 0)),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
		cfg: cfg,
	}
}

func jwtTimeDurationCal(sec int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(sec) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(sec int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(int64(sec), 0))
}

package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func AccessToken(signature string) string {
	/* สร้าง Standard Claims พร้อมกำหนดเวลาหมดอายุ*/
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	})

	ss, err := token.SignedString([]byte(signature))
	if err != nil {
		panic(err)
	}
	return ss
}

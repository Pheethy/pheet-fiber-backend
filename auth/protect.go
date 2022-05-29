package auth

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func Protect(signature []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		s := c.GetReqHeaders()
		tokenz := strings.TrimPrefix(s["Authorization"], "Bearer ")
		
		_, err := jwt.Parse(tokenz, func(token *jwt.Token) (interface{}, error) {
			//เช็ค Method ว่าเป็น MethodHMAC ไหม
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			//เช็ค Method ผ่าน return signature ให้
			return []byte(signature), nil
		})
		if err != nil {
			return fiber.ErrUnauthorized
		}

		return c.Next()
	}
}
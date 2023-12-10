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
		// Check if Authorization header exists and has at least one value
		authHeaders, ok := s["Authorization"]
		if !ok || len(authHeaders) == 0 {
			return fiber.ErrUnauthorized
		}
		tokenz := strings.TrimPrefix(authHeaders[0], "Bearer ")
		
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
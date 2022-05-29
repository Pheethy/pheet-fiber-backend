package auth

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

)

func AccessToken(signature string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		/* สร้าง Standard Claims พร้อมกำหนดเวลาหมดอายุ*/
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		})

		ss, err := token.SignedString([]byte(signature))
		if err != nil {
			c.JSON(err)
		}

		resp := map[string]interface{}{
			"token": ss,
		}

		return c.JSON(resp)
	}

}

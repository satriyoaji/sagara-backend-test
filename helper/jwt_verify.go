package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func JwtVerify(c *fiber.Ctx) (*jwt.Token, error) {
	cookie := c.Cookies("jwt")

	return jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(GoDotEnvVariable("SECRET_KEY")), nil
	})
}
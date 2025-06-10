package middleware

import (
	"app/conf"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func Jwt() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(conf.Server.Secret)},
	})
}

var JwtExpireTime = time.Hour * 7

func GenerateJwt(username string) (string, error) {
	now := time.Now().Add(JwtExpireTime)
	claims := jwt.MapClaims{
		"username": username,
		"exp":      jwt.NewNumericDate(now.Add(JwtExpireTime)),
		"iat":      jwt.NewNumericDate(now),
		"nbf":      jwt.NewNumericDate(now),
		"iss":      conf.AppName,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.Server.Secret))
}

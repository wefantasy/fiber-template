package middleware

import (
	"app/util/httputil"
	"github.com/gofiber/fiber/v2"
)

func ErrorParse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return httputil.JsonErrorParse(c, err)
		}
		return err
	}
}

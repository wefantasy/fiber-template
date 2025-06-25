package middleware

import (
	"app/log"
	"app/util/httputil"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"strings"
)

func Recover() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			log.Errorf("Panic occurred during request, requestId: %s, ip: %s, method: %s, path: %s, reason: %v",
				strings.TrimSpace(httputil.GetRequestId(c)), c.IP(), c.Method(), c.Path(), e)
		},
	})
}

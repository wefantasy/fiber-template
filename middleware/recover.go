package middleware

import (
	"app/util/httputil"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"strings"
)

func Recover() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			zap.L().Error("Panic occurred during request",
				zap.String("requestId", strings.TrimSpace(httputil.GetRequestId(c))),
				zap.String("ip", c.IP()),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.Any("reason", e),
			)
		},
	})
}

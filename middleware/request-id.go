package middleware

import (
	"app/util/httputil"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func RequestId() fiber.Handler {
	return requestid.New(requestid.Config{
		ContextKey: fiber.HeaderXRequestID,
		Generator:  httputil.GenerateRequestId,
	})
}

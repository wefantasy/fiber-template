package middleware

import (
	"app/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func RequestId() fiber.Handler {
	return requestid.New(requestid.Config{
		ContextKey: fiber.HeaderXRequestID,
		Generator:  util.GenerateRequestId,
	})
}

package middleware

import (
	"github.com/gofiber/contrib/circuitbreaker"
	"github.com/gofiber/fiber/v2"
	"time"
)

func CircuitBreaker() fiber.Handler {
	cb := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: 3,
		Timeout:          5 * time.Second,
		SuccessThreshold: 2,
	})
	return circuitbreaker.Middleware(cb)
}

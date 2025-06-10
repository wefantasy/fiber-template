package middleware

import (
	"app/log"
	"github.com/gofiber/fiber/v2"
)

func HookRoute(r fiber.Route) error {
	log.Info(r.Method + " " + r.Path)
	return nil
}

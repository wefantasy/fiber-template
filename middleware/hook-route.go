package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func HookRoute(r fiber.Route) error {
	log.Info(r.Method + " " + r.Path + " " + r.Name)
	return nil
}

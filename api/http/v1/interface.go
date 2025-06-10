package v1

import (
	"github.com/gofiber/fiber/v2"
)

type BaseContro interface {
	Name() string
	RegisterRoute(fiber.Router)
}

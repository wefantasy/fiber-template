package middleware

import (
	"app/docs"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

func Swagger() fiber.Handler {
	swaggerJSON, err := docs.SwaggerFS.ReadFile("swagger.json")
	if err != nil {
		panic(err)
	}
	return swagger.New(swagger.Config{
		BasePath:    "/",
		FileContent: swaggerJSON,
		Path:        "docs",
	})
}

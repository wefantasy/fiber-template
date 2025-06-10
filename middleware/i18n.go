package middleware

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
)

func FiberI18n() fiber.Handler {
	return fiberi18n.New(&fiberi18n.Config{
		AcceptLanguages: []language.Tag{language.Chinese, language.English},
		DefaultLanguage: language.Chinese,
		RootPath:        "./locales",
	})
}

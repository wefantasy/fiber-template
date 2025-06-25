package i18n

import (
	"app/conf"
	"app/log"
	"embed"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"io/fs"
)

//go:embed *.yaml
var localeFS embed.FS

var Bundle *i18n.Bundle

var Languages []language.Tag

func Initialize() {
	for _, lang := range conf.Languages {
		l, err := language.Parse(lang)
		if err != nil {
			log.Panic(err)
		}
		Languages = append(Languages, l)
	}
	Bundle = i18n.NewBundle(Languages[0])
	Bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	files, err := fs.Glob(localeFS, "*.yaml")
	if err != nil {
		log.Error(err)
	}
	for _, file := range files {
		log.Infof("load locale config: %s", file)
		_, err := Bundle.LoadMessageFileFS(localeFS, file)
		if err != nil {
			log.Error(err)
		}
	}
}

func Localize(id string) string {
	return LocalizeWithLang(Languages[0], id)
}

func LocalizeWithCtx(c *fiber.Ctx, id string) string {
	var matcher = language.NewMatcher(Languages)
	lang := c.Cookies("lang")
	accept := c.Get("Accept-Language")
	tag, _ := language.MatchStrings(matcher, lang, accept)
	return LocalizeWithLang(tag, id)
}

func LocalizeWithLang(lang language.Tag, id string) string {
	localizer := i18n.NewLocalizer(Bundle, lang.String())
	msg, err := localizer.LocalizeMessage(&i18n.Message{
		ID: id,
	})
	if err != nil {
		log.Debug(err)
		return ""
	}
	return msg
}

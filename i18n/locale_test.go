package i18n

import (
	"app/conf"
	"golang.org/x/text/language"
	"testing"
)

func Test_I18n(t *testing.T) {
	conf.Initialize()
	Initialize()
	t.Log(Localize("PasswordCryptFailed"))
	t.Log(LocalizeWithLang(language.English, "PasswordCryptFailed"))
}

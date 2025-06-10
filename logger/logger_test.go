package logger

import (
	"app/conf"
	"testing"

	"github.com/gofiber/fiber/v2/log"
)

func init() {
	conf.Initialize()
	Initialize()
}

func Test_Logger(t *testing.T) {
	log.Info("Test_Logger", 123)
	log.Error("Test_Error")
}

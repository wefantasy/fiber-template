package log

import (
	"app/conf"
	"go.uber.org/zap"
	"testing"
)

func init() {
	conf.Initialize()
	Initialize()
}

func Test_Logger(t *testing.T) {
	Info("Test_Logger", 123)
	Error("Test_Error")
	zap.S().Info("Test_Logger", 123)
	zap.L().Info("Test_Error")
}

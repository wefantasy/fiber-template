package log

import (
	"app/conf"
	"go.uber.org/zap"
	"testing"
)

func Test_Logger(t *testing.T) {
	conf.Initialize()
	Initialize()
	Info("Test_Logger", 123)
	Error("Test_Error")
	zap.S().Info("Test_Logger", 123)
	zap.L().Info("Test_Error")
}

package log

import (
	"app/conf"
	"app/util"
	"go.uber.org/zap"
	"testing"
)

func TestLogger(t *testing.T) {
	conf.Initialize()
	Initialize()
	Error("123")
	zap.S().Error("234")
	T(util.NewRootContext()).Error("345")
	F(nil).Error("456")
}

package conf

import (
	"testing"
)

func TestConfig(t *testing.T) {
	Initialize()
	t.Logf("%+v", Conf)
}

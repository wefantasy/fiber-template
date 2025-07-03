package conf

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	Initialize()
	t.Logf("config: %+v", Conf)
}

func TestEnv(t *testing.T) {
	os.Setenv("APPNAME", "testEnvName")
	Initialize()
	t.Logf("config: %+v", Conf)
}

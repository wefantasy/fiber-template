package conf

import (
	"testing"
)

func Test_Insert(t *testing.T) {
	Initialize()
	t.Log(Conf.Languages)
}

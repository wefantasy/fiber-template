package code

import (
	"testing"
)

func Test_Error(t *testing.T) {
	t.Log(ParamError.ToError())
}

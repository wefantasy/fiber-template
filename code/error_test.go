package code

import (
	"fmt"
	"testing"
)

func Test_Error(t *testing.T) {
	err := fmt.Errorf("")
	t.Log(IsSuccess(err))
}

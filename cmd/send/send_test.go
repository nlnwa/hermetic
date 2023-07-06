package send

import (
	"testing"
)

func TestValidateCmdName(t *testing.T) {
	if err := send(); err != nil {
		panic("send failed")
	}
}

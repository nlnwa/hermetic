package verify

import "testing"

func TestValidateCmdName(t *testing.T) {
	if err := verify(); err != nil {
		panic("verify failed")
	}
}

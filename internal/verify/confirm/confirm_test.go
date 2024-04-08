package confirmImplmentation

import (
	"testing"
)

func TestVerify(t *testing.T) {
	t.Run("Verify function test", func(t *testing.T) {
		err := Verify(("test-topic"))
		if err != nil {
			t.Errorf("Verify() returned error: %v", err)
		}
	})
}

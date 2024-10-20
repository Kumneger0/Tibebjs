package runtime

import "testing"

func TestInitialize(t *testing.T) {
	ctx := TestInitialize(t)
	if ctx == nil {
		t.Errorf("Failed to initialize the JS runtime")
	}
}

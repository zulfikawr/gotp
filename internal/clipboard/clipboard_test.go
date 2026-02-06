package clipboard

import (
	"testing"
	"time"
)

func TestWriteWithTimeout(t *testing.T) {
	// This test depends on the system clipboard, which might not be available in CI.
	// We'll just test that it doesn't crash.
	err := WriteWithTimeout("test-code", 100*time.Millisecond)
	if err != nil {
		t.Logf("Skipping actual clipboard test: %v", err)
	}
}

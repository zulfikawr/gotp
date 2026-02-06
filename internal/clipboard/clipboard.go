package clipboard

import (
	"time"

	"github.com/atotto/clipboard"
)

// WriteWithTimeout copies the text to the clipboard and clears it after the specified timeout.
func WriteWithTimeout(text string, timeout time.Duration) error {
	if err := clipboard.WriteAll(text); err != nil {
		return err
	}

	if timeout > 0 {
		go func() {
			time.Sleep(timeout)
			// Only clear if the current content is still what we wrote
			current, err := clipboard.ReadAll()
			if err == nil && current == text {
				_ = clipboard.WriteAll("")
			}
		}()
	}

	return nil
}

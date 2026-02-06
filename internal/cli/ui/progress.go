package ui

import (
	"fmt"
	"strings"
)

// ProgressBar returns a string representing a progress bar with Gruvbox colors.
func ProgressBar(current, total int, width int) string {
	if total <= 0 {
		return ""
	}
	percent := float64(current) / float64(total)
	filledLength := int(float64(width) * percent)
	if filledLength > width {
		filledLength = width
	}
	if filledLength < 0 {
		filledLength = 0
	}

	color := SuccessBright
	if percent < 0.2 {
		color = DangerBright
	} else if percent < 0.5 {
		color = WarningBright
	}

	bar := color + strings.Repeat("━", filledLength) + Reset + TextMuted + strings.Repeat("─", width-filledLength) + Reset
	return bar
}

// PrintCodeDisplay displays the TOTP code with a themed layout.
func PrintCodeDisplay(name, code string, remaining, total int) {
	fmt.Fprintf(Out, "%s%s%s\n", PrimaryBright+Bold, name, Reset)
	fmt.Fprintf(Out, "%s└──  %s%s%s  |  %s %s%d%ss remaining...\n", TextMuted, WarningBright+Bold, code, Reset, ProgressBar(remaining, total, 10), TextMuted, remaining, Reset)
}

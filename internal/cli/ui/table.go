package ui

import (
	"fmt"
	"strings"
)

// PrintTable prints a simple formatted table to ui.Out with Gruvbox theme.
func PrintTable(headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print header
	for i, h := range headers {
		fmt.Fprintf(Out, "%s%-*s%s  ", PrimaryBright+Bold, widths[i], h, Reset)
	}
	fmt.Fprintln(Out)

	// Print separator
	for _, w := range widths {
		fmt.Fprintf(Out, "%s%s%s  ", TextMuted, strings.Repeat("─", w), Reset)
	}
	fmt.Fprintln(Out)

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				fmt.Fprintf(Out, "%s%-*s%s  ", TextPrimary, widths[i], cell, Reset)
			}
		}
		fmt.Fprintln(Out)
	}
}

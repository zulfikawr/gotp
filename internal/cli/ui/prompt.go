package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

var (
	// In is the reader used for inputs. Can be overridden in tests.
	In io.Reader = os.Stdin
	// Out is the writer used for prompts. Can be overridden in tests.
	Out io.Writer = os.Stdout
	// PasswordReader is the function used to read passwords.
	PasswordReader func(int) ([]byte, error) = term.ReadPassword
	// IsTerminal is the function used to check if a file descriptor is a terminal.
	// Can be overridden in tests to simulate a TTY.
	IsTerminal func(int) bool = term.IsTerminal

	// scanner is used to read from In.
	scanner *bufio.Scanner

	// UseColor determines if ANSI color codes should be used.
	UseColor = true
)

// Gruvbox Theme Constants
const (
	BgHardCode   = "\033[38;2;29;32;33m"
	BgMediumCode = "\033[38;2;40;40;40m"
	BgSoftCode   = "\033[38;2;50;48;47m"
	SurfaceCode  = "\033[38;2;60;56;54m"

	TextPrimaryCode   = "\033[38;2;235;219;178m"
	TextContrastCode  = "\033[38;2;251;241;199m"
	TextSecondaryCode = "\033[38;2;189;174;147m"
	TextMutedCode     = "\033[38;2;146;131;116m"

	DangerCode    = "\033[38;2;204;36;29m"
	AttentionCode = "\033[38;2;214;93;14m"
	WarningCode   = "\033[38;2;215;153;33m"
	SuccessCode   = "\033[38;2;152;151;26m"
	InfoCode      = "\033[38;2;104;157;106m"
	PrimaryCode   = "\033[38;2;69;133;136m"
	AccentCode    = "\033[38;2;177;98;134m"

	DangerBrightCode    = "\033[38;2;251;73;52m"
	AttentionBrightCode = "\033[38;2;254;128;25m"
	WarningBrightCode   = "\033[38;2;250;189;47m"
	SuccessBrightCode   = "\033[38;2;184;187;38m"
	InfoBrightCode      = "\033[38;2;142;192;124m"
	PrimaryBrightCode   = "\033[38;2;131;165;152m"
	AccentBrightCode    = "\033[38;2;211;134;155m"

	ResetCode = "\033[0m"
	BoldCode  = "\033[1m"
)

var (
	BgHard   = BgHardCode
	BgMedium = BgMediumCode
	BgSoft   = BgSoftCode
	Surface  = SurfaceCode

	TextPrimary   = TextPrimaryCode
	TextContrast  = TextContrastCode
	TextSecondary = TextSecondaryCode
	TextMuted     = TextMutedCode

	Danger    = DangerCode
	Attention = AttentionCode
	Warning   = WarningCode
	Success   = SuccessCode
	Info      = InfoCode
	Primary   = PrimaryCode
	Accent    = AccentCode

	DangerBright    = DangerBrightCode
	AttentionBright = AttentionBrightCode
	WarningBright   = WarningBrightCode
	SuccessBright   = SuccessBrightCode
	InfoBright      = InfoBrightCode
	PrimaryBright   = PrimaryBrightCode
	AccentBright    = AccentBrightCode

	Reset = ResetCode
	Bold  = BoldCode
)

// SetColor enables or disables ANSI color codes.
func SetColor(enabled bool) {
	UseColor = enabled
	if !enabled {
		BgHard, BgMedium, BgSoft, Surface = "", "", "", ""
		TextPrimary, TextContrast, TextSecondary, TextMuted = "", "", "", ""
		Danger, Attention, Warning, Success, Info, Primary, Accent = "", "", "", "", "", "", ""
		DangerBright, AttentionBright, WarningBright, SuccessBright, InfoBright, PrimaryBright, AccentBright = "", "", "", "", "", "", ""
		Reset = ""
		Bold = ""
	} else {
		BgHard, BgMedium, BgSoft, Surface = BgHardCode, BgMediumCode, BgSoftCode, SurfaceCode
		TextPrimary, TextContrast, TextSecondary, TextMuted = TextPrimaryCode, TextContrastCode, TextSecondaryCode, TextMutedCode
		Danger, Attention, Warning, Success, Info, Primary, Accent = DangerCode, AttentionCode, WarningCode, SuccessCode, InfoCode, PrimaryCode, AccentCode
		DangerBright, AttentionBright, WarningBright, SuccessBright, InfoBright, PrimaryBright, AccentBright = DangerBrightCode, AttentionBrightCode, WarningBrightCode, SuccessBrightCode, InfoBrightCode, PrimaryBrightCode, AccentBrightCode
		Reset = ResetCode
		Bold = BoldCode
	}
}

func getScanner() *bufio.Scanner {
	if scanner == nil {
		scanner = bufio.NewScanner(In)
	}
	return scanner
}

// ResetScanner resets the scanner, useful for tests when In changes.
func ResetScanner() {
	scanner = nil
}

// PromptPassword strictly prompts the user for a password from a terminal.
// It will fail if the input is not a TTY to prevent password leakage.
func PromptPassword(prompt string) ([]byte, error) {
	fmt.Fprintf(Out, "%s%s%s", Primary, prompt, Reset)

	if f, ok := In.(*os.File); ok && IsTerminal(int(f.Fd())) {
		password, err := PasswordReader(int(f.Fd()))
		fmt.Fprintln(Out)
		return password, err
	}

	// For tests or non-interactive shells, handle based on PasswordReader
	// if PasswordReader has been mocked (e.g., in tests)
	if PasswordReader != nil && In != os.Stdin {
		password, err := PasswordReader(0)
		fmt.Fprintln(Out)
		return password, err
	}

	fmt.Fprintln(Out)
	return nil, fmt.Errorf("secure password prompt requires a terminal")
}

// PromptConfirm prompts the user for a yes/no confirmation.
func PromptConfirm(prompt string, defaultYes bool) bool {
	suffix := "[y/N]"
	if defaultYes {
		suffix = "[Y/n]"
	}
	fmt.Fprintf(Out, "%s%s%s %s%s%s: ", Primary, prompt, Reset, TextMuted, suffix, Reset)

	s := getScanner()
	if s.Scan() {
		input := strings.ToLower(strings.TrimSpace(s.Text()))
		if input == "" {
			return defaultYes
		}
		return input == "y" || input == "yes"
	}
	return defaultYes
}

// PromptString prompts the user for a string input.
func PromptString(prompt string, defaultValue string) string {
	if defaultValue != "" {
		fmt.Fprintf(Out, "%s%s%s %s[%s]%s: ", Primary, prompt, Reset, TextMuted, defaultValue, Reset)
	} else {
		fmt.Fprintf(Out, "%s%s%s: ", Primary, prompt, Reset)
	}

	s := getScanner()
	if s.Scan() {
		input := strings.TrimSpace(s.Text())
		if input == "" {
			return defaultValue
		}
		return input
	}
	return defaultValue
}

// PromptRequired prompts for a string that cannot be empty.
func PromptRequired(prompt string) string {
	for {
		fmt.Fprintf(Out, "%s%s%s: ", Primary, prompt, Reset)
		s := getScanner()
		if s.Scan() {
			input := strings.TrimSpace(s.Text())
			if input != "" {
				return input
			}
			fmt.Fprintf(Out, "%sError: This field is required.%s\n", DangerBright, Reset)
		} else {
			return ""
		}
	}
}

// PromptValidate prompts for a string and validates it with a function.
func PromptValidate(prompt string, validate func(string) error) string {
	for {
		fmt.Fprintf(Out, "%s%s%s: ", Primary, prompt, Reset)
		s := getScanner()
		if s.Scan() {
			input := strings.TrimSpace(s.Text())
			if err := validate(input); err == nil {
				return input
			} else {
				fmt.Fprintf(Out, "%sError: %v%s\n", DangerBright, err, Reset)
			}
		} else {
			return ""
		}
	}
}

// Dimmed returns the text wrapped in dimmed ANSI codes.
func Dimmed(text string) string {
	return TextMuted + text + Reset
}

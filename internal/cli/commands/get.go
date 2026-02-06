package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zulfikawr/gotp/internal/cli/ui"
	"github.com/zulfikawr/gotp/internal/clipboard"
	"github.com/zulfikawr/gotp/internal/config"
	"github.com/zulfikawr/gotp/internal/totp"
	"github.com/zulfikawr/gotp/internal/vault"
	"github.com/zulfikawr/gotp/pkg/base32"
)

func NewGetCmd() *cobra.Command {
	var copyToClipboard bool
	var timeout int
	var watch bool

	cmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Get TOTP code for an account",
		Long:  `Generate and display the current Time-based One-Time Password (TOTP) code for a stored account. Includes a live-updating watch mode and clipboard integration.`,
		Args:  cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			vaultPath := config.GetVaultPath()
			isJSON, _ := cmd.Flags().GetBool("json")

			// Check if vault exists first
			if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
				fmt.Fprintf(ui.Out, "%sError: Vault file not found at %s%s\n", ui.DangerBright, vaultPath, ui.Reset)
				fmt.Fprintf(ui.Out, "%sTip: Run '%s%sgotp %sinit%s' to create a new secure vault.%s\n", ui.TextMuted, ui.Reset, ui.SuccessBright, ui.WarningBright, ui.TextMuted, ui.Reset)
				return nil
			}

			v, _, err := vault.LoadVaultInteractive(vaultPath, ui.PromptPassword)
			if err != nil {
				fmt.Fprintf(ui.Out, "%sError: %v%s\n", ui.DangerBright, err, ui.Reset)
				return nil
			}

			var target *vault.Account
			for i := range v.Accounts {
				if strings.EqualFold(v.Accounts[i].Name, name) {
					target = &v.Accounts[i]
					break
				}
			}

			if target == nil {
				fmt.Fprintf(ui.Out, "%sError: Account %q not found%s\n", ui.DangerBright, name, ui.Reset)
				return nil
			}

			secretBytes, err := base32.Decode(string(target.Secret))
			if err != nil {
				fmt.Fprintf(ui.Out, "%sError: Failed to decode secret: %v%s\n", ui.DangerBright, err, ui.Reset)
				return nil
			}

			if watch {
				if isJSON {
					fmt.Fprintf(ui.Out, "%sError: Watch mode is not compatible with JSON output%s\n", ui.DangerBright, ui.Reset)
					return nil
				}

				// Set up signal handling to restore cursor on Ctrl+C
				sigChan := make(chan os.Signal, 1)
				signal.Notify(sigChan, os.Interrupt)
				defer signal.Stop(sigChan)

				// Function to clean up and restore cursor
				cleanup := func() {
					fmt.Fprint(ui.Out, "\033[J")    // Clear from cursor to end of screen
					fmt.Fprint(ui.Out, "\033[?25h") // Show cursor
				}
				defer cleanup()

				// Hide cursor
				fmt.Fprint(ui.Out, "\033[?25l")

				// Channel to signal watch loop to exit
				done := make(chan bool, 1)

				// Goroutine to handle signals
				go func() {
					<-sigChan
					done <- true
				}()

				// Watch loop
				for {
					select {
					case <-done:
						// Signal received, move to start of display and clear
						fmt.Fprintf(ui.Out, "\r\033[J")
						return nil
					default:
						now := time.Now()
						code, _ := totp.GenerateTOTP(totp.TOTPParams{
							Secret:    secretBytes,
							Timestamp: now,
							Period:    target.Period,
							Digits:    target.Digits,
							Algorithm: target.Algorithm,
						})
						remaining := totp.RemainingSeconds(now, target.Period)

						// Move to start of line, clear to end of screen, then print
						fmt.Fprintf(ui.Out, "\r\033[J")
						ui.PrintCodeDisplay(target.Name, code, remaining, target.Period)
						// Move back up to the start of the display for the next frame
						fmt.Fprintf(ui.Out, "\033[2A")
						time.Sleep(500 * time.Millisecond)
					}
				}
			}

			now := time.Now()
			code, err := totp.GenerateTOTP(totp.TOTPParams{
				Secret:    secretBytes,
				Timestamp: now,
				Period:    target.Period,
				Digits:    target.Digits,
				Algorithm: target.Algorithm,
			})
			if err != nil {
				fmt.Fprintf(ui.Out, "%sError: Failed to generate code: %v%s\n", ui.DangerBright, err, ui.Reset)
				return nil
			}

			if isJSON {
				res := map[string]interface{}{
					"account": target.Name,
					"code":    code,
				}
				data, _ := json.Marshal(res)
				fmt.Fprintln(ui.Out, string(data))
			} else {
				remaining := totp.RemainingSeconds(now, target.Period)
				ui.PrintCodeDisplay(target.Name, code, remaining, target.Period)
			}

			if copyToClipboard {
				if err := clipboard.WriteWithTimeout(code, time.Duration(timeout)*time.Second); err != nil {
					fmt.Fprintf(ui.Out, "%sWarning: failed to copy to clipboard: %v%s\n", ui.WarningBright, err, ui.Reset)
				} else if !isJSON {
					fmt.Fprintf(ui.Out, "%s✓ Code copied to clipboard (clears in %ds)%s\n", ui.SuccessBright, timeout, ui.Reset)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&copyToClipboard, "copy", "c", false, "Copy code to clipboard")
	cmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "Clipboard clear timeout in seconds")
	cmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch mode (continuous update)")
	return cmd
}

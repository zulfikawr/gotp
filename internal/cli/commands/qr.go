package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zulfikawr/gotp/internal/cli/ui"
	"github.com/zulfikawr/gotp/internal/config"
	"github.com/zulfikawr/gotp/internal/qr"
	"github.com/zulfikawr/gotp/internal/vault"
)

func NewQrCmd() *cobra.Command {
	var output string
	var size int
	var terminal bool
	var compact bool

	cmd := &cobra.Command{
		Use:   "qr <account>",
		Short: "Generate or parse QR codes",
		Long: `Generate QR codes for TOTP accounts or parse QR code images.

Generate: Creates a QR code image from an account's otpauth:// URI.
Parse: Extracts an otpauth:// URI from a QR code image file.

Examples:
  gotp qr "My Account" --output qr.png
  gotp qr --parse qr.png
  gotp qr "My Account" --terminal
  gotp qr "My Account" --terminal --compact`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath := config.GetVaultPath()

			// Parse mode
			if cmd.Flags().Changed("parse") {
				parseFile, _ := cmd.Flags().GetString("parse")
				if parseFile == "" {
					return fmt.Errorf("parse file path required")
				}

				fmt.Fprintf(ui.Out, "%sParsing QR code: %s%s\n", ui.InfoBright, parseFile, ui.Reset)

				uri, err := qr.ParseImageFile(parseFile)
				if err != nil {
					fmt.Fprintf(ui.Out, "%s✗ Failed to parse QR code: %v%s\n", ui.DangerBright, err, ui.Reset)
					return err
				}

				// Validate it's an otpauth URI
				if err := qr.ValidateOTPAuthURI(uri); err != nil {
					fmt.Fprintf(ui.Out, "%s✗ Invalid URI format: %v%s\n", ui.DangerBright, err, ui.Reset)
					return err
				}

				fmt.Fprintf(ui.Out, "%s✓ Successfully parsed QR code%s\n", ui.SuccessBright, ui.Reset)
				fmt.Fprintf(ui.Out, "%sURI: %s%s\n", ui.PrimaryBright, uri, ui.Reset)

				// Try to parse as account
				acc, err := vault.FromURI(uri)
				if err == nil {
					fmt.Fprintf(ui.Out, "%sAccount: %s (Issuer: %s, Username: %s)%s\n",
						ui.PrimaryBright, acc.Name, acc.Issuer, acc.Username, ui.Reset)
				}

				return nil
			}

			// Generate mode requires an account name
			if len(args) == 0 {
				fmt.Fprintf(ui.Out, "%s✗ Account name required for QR generation%s\n", ui.DangerBright, ui.Reset)
				return fmt.Errorf("account name required for QR generation")
			}

			accountName := args[0]

			fmt.Fprintf(ui.Out, "%sLoading vault...%s\n", ui.InfoBright, ui.Reset)

			// Load vault
			v, _, err := vault.LoadVaultInteractive(vaultPath, ui.PromptPassword)
			if err != nil {
				fmt.Fprintf(ui.Out, "%s✗ Failed to load vault: %v%s\n", ui.DangerBright, err, ui.Reset)
				return err
			}

			// Find account
			var targetAccount *vault.Account
			for _, acc := range v.Accounts {
				if strings.EqualFold(acc.Name, accountName) {
					targetAccount = &acc
					break
				}
			}

			if targetAccount == nil {
				fmt.Fprintf(ui.Out, "%s✗ Account not found: %s%s\n", ui.DangerBright, accountName, ui.Reset)
				return fmt.Errorf("account not found: %s", accountName)
			}

			// Generate URI
			uri := targetAccount.ToURI()

			// Terminal mode
			if terminal {
				fmt.Fprintf(ui.Out, "%sGenerating QR code for: %s%s\n", ui.PrimaryBright, targetAccount.Name, ui.Reset)
				fmt.Fprintf(ui.Out, "%sURI: %s%s\n\n", ui.InfoBright, uri, ui.Reset)

				if err := qr.GenerateQRCodeToTerminal(uri); err != nil {
					fmt.Fprintf(ui.Out, "%s✗ Failed to generate terminal QR code: %v%s\n", ui.DangerBright, err, ui.Reset)
					return fmt.Errorf("failed to generate terminal QR code: %w", err)
				}
				return nil
			}

			// Generate image file
			if output == "" {
				output = fmt.Sprintf("%s_qr.png", strings.ReplaceAll(targetAccount.Name, " ", "_"))
			}

			fmt.Fprintf(ui.Out, "%sGenerating QR code...%s\n", ui.InfoBright, ui.Reset)

			if err := qr.GenerateQRCodeToFile(uri, output, size); err != nil {
				fmt.Fprintf(ui.Out, "%s✗ Failed to generate QR code: %v%s\n", ui.DangerBright, err, ui.Reset)
				return fmt.Errorf("failed to generate QR code: %w", err)
			}

			fmt.Fprintf(ui.Out, "%s✓ QR code generated: %s%s\n", ui.SuccessBright, output, ui.Reset)
			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: <account>_qr.png)")
	cmd.Flags().IntVarP(&size, "size", "s", 256, "QR code size in pixels")
	cmd.Flags().BoolVar(&terminal, "terminal", false, "Display QR code in terminal")
	cmd.Flags().BoolVar(&compact, "compact", false, "Use compact display (terminal only)")
	cmd.Flags().String("parse", "", "Parse a QR code image file")

	return cmd
}

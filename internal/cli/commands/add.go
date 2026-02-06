package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/zulfikawr/gotp/internal/cli/ui"
	"github.com/zulfikawr/gotp/internal/config"
	"github.com/zulfikawr/gotp/internal/totp"
	"github.com/zulfikawr/gotp/internal/vault"
	"github.com/zulfikawr/gotp/pkg/base32"
)

func NewAddCmd() *cobra.Command {
	var secret, issuer, username, algo, uri string
	var digits, period int
	var tags []string

	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add a new TOTP account",
		Long:  `Add a new TOTP account to your secure vault. You can either provide the details manually via flags or interactive mode, or use an otpauth:// URI.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath := config.GetVaultPath()

			// Check if vault exists first
			if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
				fmt.Fprintf(ui.Out, "%sError: Vault file not found at %s%s\n", ui.DangerBright, vaultPath, ui.Reset)
				fmt.Fprintf(ui.Out, "%sTip: Run '%s%sgotp %sinit%s' to create a new secure vault.%s\n", ui.TextMuted, ui.Reset, ui.SuccessBright, ui.WarningBright, ui.TextMuted, ui.Reset)
				return nil // Exit gracefully
			}

			v, key, err := vault.LoadVaultInteractive(vaultPath, ui.PromptPassword)
			if err != nil {
				fmt.Fprintf(ui.Out, "%sError: %v%s\n", ui.DangerBright, err, ui.Reset)
				return nil
			}

			var acc *vault.Account
			if uri != "" {
				acc, err = vault.FromURI(uri)
				if err != nil {
					return err
				}
			} else {
				var name string
				if len(args) > 0 {
					name = args[0]
				}
				if name == "" {
					fmt.Fprintln(ui.Out, ui.Dimmed("No account name provided. Entering interactive mode..."))
					name = ui.PromptRequired("Account Name (e.g., GitHub)")
				}

				if secret == "" {
					secret = ui.PromptValidate("Secret (Base32)", func(s string) error {
						s = strings.ToUpper(strings.ReplaceAll(s, " ", ""))
						_, err := base32.Decode(s)
						return err
					})
				}
				secret = strings.ToUpper(strings.ReplaceAll(secret, " ", ""))

				if issuer == "" && !cmd.Flags().Changed("issuer") {
					issuer = ui.PromptString("Issuer (Optional)", "")
				}
				if username == "" && !cmd.Flags().Changed("username") {
					username = ui.PromptString("Username/Email (Optional)", "")
				}

				acc = vault.NewAccount(name, []byte(secret))
				acc.Issuer = issuer
				acc.Username = username
				if cmd.Flags().Changed("digits") {
					acc.Digits = digits
				}
				if cmd.Flags().Changed("period") {
					acc.Period = period
				}
				if cmd.Flags().Changed("algorithm") {
					acc.Algorithm = totp.HashAlgorithm(strings.ToUpper(algo))
				}
			}

			acc.ID = uuid.New().String()
			acc.Tags = tags
			v.Accounts = append(v.Accounts, *acc)

			if err := vault.CreateBackup(vaultPath, 3); err != nil {
				fmt.Fprintf(ui.Out, "%sWarning: failed to create backup: %v%s\n", ui.WarningBright, err, ui.Reset)
			}

			if err := vault.SaveVaultWithKey(vaultPath, v, key); err != nil {
				return err
			}

			fmt.Fprintf(ui.Out, "%s✓ Added account: %s%s\n", ui.SuccessBright, acc.Name, ui.Reset)
			return nil
		},
	}

	cmd.Flags().StringVar(&uri, "uri", "", "Add from otpauth:// URI")
	cmd.Flags().StringVarP(&secret, "secret", "s", "", "Base32 secret")
	cmd.Flags().StringVarP(&issuer, "issuer", "i", "", "Service issuer")
	cmd.Flags().StringVarP(&username, "username", "u", "", "Username/Email")
	cmd.Flags().StringVarP(&algo, "algorithm", "a", "SHA1", "Hash algorithm")
	cmd.Flags().IntVarP(&digits, "digits", "d", 6, "Code digits")
	cmd.Flags().IntVarP(&period, "period", "p", 30, "Time period")
	cmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, "Comma-separated tags")

	return cmd
}

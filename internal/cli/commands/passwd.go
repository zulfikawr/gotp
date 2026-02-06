package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zulfikawr/gotp/internal/cli/ui"
	"github.com/zulfikawr/gotp/internal/config"
	"github.com/zulfikawr/gotp/internal/crypto"
	"github.com/zulfikawr/gotp/internal/vault"
)

func NewPasswdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "passwd",
		Short: "Change master password",
		Long:  `Securely update the master password for your secure vault. This will re-encrypt all stored accounts using the new password.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath := config.GetVaultPath()

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

			newPassword, err := ui.PromptPassword("Enter new master password: ")
			if err != nil {
				return err
			}
			confirm, err := ui.PromptPassword("Confirm new master password: ")
			if err != nil {
				return err
			}

			if !crypto.SecureCompare(newPassword, confirm) {
				fmt.Fprintf(ui.Out, "%sError: Passwords do not match%s\n", ui.DangerBright, ui.Reset)
				return nil
			}

			if err := vault.CreateBackup(vaultPath, 3); err != nil {
				fmt.Fprintf(ui.Out, "%sWarning: failed to create backup: %v%s\n", ui.WarningBright, err, ui.Reset)
			}

			if err := vault.SaveVault(vaultPath, v, newPassword); err != nil {
				fmt.Fprintf(ui.Out, "%sError: Failed to save vault: %v%s\n", ui.DangerBright, err, ui.Reset)
				return nil
			}

			// Clear session on password change
			_ = vault.ClearSession()

			fmt.Fprintf(ui.Out, "%s✓ Master password changed successfully%s\n", ui.SuccessBright, ui.Reset)
			return nil
		},
	}

	return cmd
}

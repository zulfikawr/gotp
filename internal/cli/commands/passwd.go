package commands

import (
	"fmt"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath := config.GetVaultPath()

			v, _, err := vault.LoadVaultInteractive(vaultPath, ui.PromptPassword)
			if err != nil {
				return err
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
				return fmt.Errorf("passwords do not match")
			}

			if err := vault.CreateBackup(vaultPath, 3); err != nil {
				fmt.Fprintf(ui.Out, "%sWarning: failed to create backup: %v%s\n", ui.WarningBright, err, ui.Reset)
			}

			if err := vault.SaveVault(vaultPath, v, newPassword); err != nil {
				return err
			}

			// Clear session on password change
			_ = vault.ClearSession()

			fmt.Fprintf(ui.Out, "%s✓ Master password changed successfully%s\n", ui.SuccessBright, ui.Reset)
			return nil
		},
	}

	return cmd
}

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

func NewInitCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new vault",
		Long:  `Create a new secure vault for storing your TOTP accounts. Requires a master password that will be used for encryption and authentication.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath := config.GetVaultPath()

			if _, err := os.Stat(vaultPath); err == nil && !force {
				fmt.Fprintf(ui.Out, "%sError: Vault already exists at %s%s\n", ui.DangerBright, vaultPath, ui.Reset)
				fmt.Fprintf(ui.Out, "%sTip: Use the '%s--force%s' flag to overwrite the existing vault.%s\n", ui.TextMuted, ui.InfoBright, ui.TextMuted, ui.Reset)
				return nil
			}

			password, err := ui.PromptPassword("Enter master password: ")
			if err != nil {
				return err
			}
			confirm, err := ui.PromptPassword("Confirm master password: ")
			if err != nil {
				return err
			}

			if !crypto.SecureCompare(password, confirm) {
				fmt.Fprintf(ui.Out, "%sError: Passwords do not match%s\n", ui.DangerBright, ui.Reset)
				return nil
			}

			salt, err := crypto.GenerateSalt(16)
			if err != nil {
				fmt.Fprintf(ui.Out, "%sError: Failed to generate salt: %v%s\n", ui.DangerBright, err, ui.Reset)
				return nil
			}

			v := vault.NewVault(salt)
			err = vault.SaveVault(vaultPath, v, password)
			if err != nil {
				fmt.Fprintf(ui.Out, "%sError: Failed to save vault: %v%s\n", ui.DangerBright, err, ui.Reset)
				return nil
			}

			fmt.Fprintf(ui.Out, "%s✓ Vault created successfully at %s%s\n", ui.SuccessBright, vaultPath, ui.Reset)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing vault")
	return cmd
}

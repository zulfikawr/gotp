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
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath := config.GetVaultPath()

			if _, err := os.Stat(vaultPath); err == nil && !force {
				return fmt.Errorf("vault already exists at %s. Use --force to overwrite", vaultPath)
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
				return fmt.Errorf("passwords do not match")
			}

			salt, err := crypto.GenerateSalt(16)
			if err != nil {
				return err
			}

			v := vault.NewVault(salt)
			err = vault.SaveVault(vaultPath, v, password)
			if err != nil {
				return err
			}

			fmt.Fprintf(ui.Out, "%s✓ Vault created successfully at %s%s\n", ui.SuccessBright, vaultPath, ui.Reset)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing vault")
	return cmd
}

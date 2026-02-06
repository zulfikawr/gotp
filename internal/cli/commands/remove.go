package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zulfikawr/gotp/internal/cli/ui"
	"github.com/zulfikawr/gotp/internal/config"
	"github.com/zulfikawr/gotp/internal/vault"
)

func NewRemoveCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove an account from the vault",
		Long:  `Permanently remove a TOTP account from your secure vault. Requires a confirmation unless the --force flag is used.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			vaultPath := config.GetVaultPath()

			v, key, err := vault.LoadVaultInteractive(vaultPath, ui.PromptPassword)
			if err != nil {
				return err
			}

			index := -1
			for i := range v.Accounts {
				if strings.EqualFold(v.Accounts[i].Name, name) {
					index = i
					break
				}
			}

			if index == -1 {
				return fmt.Errorf("account %q not found", name)
			}

			if !force {
				confirm := ui.PromptConfirm(fmt.Sprintf("Are you sure you want to remove %q?", name), false)
				if !confirm {
					fmt.Fprintln(ui.Out, "Operation cancelled.")
					return nil
				}
			}

			v.Accounts = append(v.Accounts[:index], v.Accounts[index+1:]...)

			if err := vault.CreateBackup(vaultPath, 3); err != nil {
				fmt.Fprintf(ui.Out, "%sWarning: failed to create backup: %v%s\n", ui.WarningBright, err, ui.Reset)
			}

			if err := vault.SaveVaultWithKey(vaultPath, v, key); err != nil {
				return err
			}

			fmt.Fprintf(ui.Out, "%s✓ Removed account: %s%s\n", ui.SuccessBright, name, ui.Reset)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation")
	return cmd
}

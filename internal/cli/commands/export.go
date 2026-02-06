package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zulfikawr/gotp/internal/cli/ui"
	"github.com/zulfikawr/gotp/internal/config"
	"github.com/zulfikawr/gotp/internal/crypto"
	"github.com/zulfikawr/gotp/internal/vault"
)

func NewExportCmd() *cobra.Command {
	var format, outputPath string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export accounts for backup",
		Long:  `Export your stored accounts for backup or migration. Supports JSON, otpauth:// URIs, and password-protected encrypted formats.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath := config.GetVaultPath()

			v, _, err := vault.LoadVaultInteractive(vaultPath, ui.PromptPassword)
			if err != nil {
				return err
			}

			var output []byte
			switch format {
			case "json":
				if !ui.PromptConfirm("⚠ WARNING: This will export secrets in PLAIN TEXT. Continue?", false) {
					fmt.Fprintln(ui.Out, "Export cancelled.")
					return nil
				}
				output, err = json.MarshalIndent(v.Accounts, "", "  ")
				if err != nil {
					return err
				}
			case "uri":
				if !ui.PromptConfirm("⚠ WARNING: This will export secrets in PLAIN TEXT. Continue?", false) {
					fmt.Fprintln(ui.Out, "Export cancelled.")
					return nil
				}
				var uris string
				for _, acc := range v.Accounts {
					uris += acc.ToURI() + "\n"
				}
				output = []byte(uris)
			case "encrypted":
				exportPass, err := ui.PromptPassword("Enter password for encrypted export: ")
				if err != nil {
					return err
				}
				confirmPass, _ := ui.PromptPassword("Confirm export password: ")
				if !crypto.SecureCompare(exportPass, confirmPass) {
					return fmt.Errorf("passwords do not match")
				}

				salt, _ := crypto.GenerateSalt(16)
				exportVault := vault.NewVault(salt)
				exportVault.Accounts = v.Accounts

				ciphertext, err := exportVault.Marshal(exportPass)
				if err != nil {
					return fmt.Errorf("encryption failed: %v", err)
				}

				metadata := vault.VaultMetadata{
					Salt:       exportVault.Salt,
					KDFParams:  exportVault.KDFParams,
					Ciphertext: ciphertext,
				}
				output, _ = json.Marshal(metadata)

			default:
				return fmt.Errorf("unsupported format: %s. Use 'json', 'uri', or 'encrypted'", format)
			}

			if outputPath != "" {
				err = os.WriteFile(outputPath, output, 0600)
				if err != nil {
					return err
				}
				fmt.Fprintf(ui.Out, "%s✓ Exported %d accounts to %s%s\n", ui.SuccessBright, len(v.Accounts), outputPath, ui.Reset)
			} else {
				fmt.Fprintln(ui.Out, string(output))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "json", "Export format (json, uri, encrypted)")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path (default: stdout)")

	return cmd
}

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

			var output []byte
			switch format {
			case "json":
				if !ui.PromptConfirm("⚠ WARNING: This will export secrets in PLAIN TEXT. Continue?", false) {
					fmt.Fprintln(ui.Out, "Export cancelled.")
					return nil
				}
				output, err = json.MarshalIndent(v.Accounts, "", "  ")
				if err != nil {
					fmt.Fprintf(ui.Out, "%sError: Failed to marshal JSON: %v%s\n", ui.DangerBright, err, ui.Reset)
					return nil
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
					fmt.Fprintf(ui.Out, "%sError: Passwords do not match%s\n", ui.DangerBright, ui.Reset)
					return nil
				}

				salt, _ := crypto.GenerateSalt(16)
				exportVault := vault.NewVault(salt)
				exportVault.Accounts = v.Accounts

				ciphertext, err := exportVault.Marshal(exportPass)
				if err != nil {
					fmt.Fprintf(ui.Out, "%sError: Encryption failed: %v%s\n", ui.DangerBright, err, ui.Reset)
					return nil
				}

				metadata := vault.VaultMetadata{
					Salt:       exportVault.Salt,
					KDFParams:  exportVault.KDFParams,
					Ciphertext: ciphertext,
				}
				output, _ = json.Marshal(metadata)

			default:
				fmt.Fprintf(ui.Out, "%sError: Unsupported format: %s%s\n", ui.DangerBright, format, ui.Reset)
				fmt.Fprintf(ui.Out, "%sTip: Use 'json', 'uri', or 'encrypted'.%s\n", ui.TextMuted, ui.Reset)
				return nil
			}

			if outputPath != "" {
				err = os.WriteFile(outputPath, output, 0600)
				if err != nil {
					fmt.Fprintf(ui.Out, "%sError: Failed to write file: %v%s\n", ui.DangerBright, err, ui.Reset)
					return nil
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

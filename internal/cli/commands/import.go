package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/zulfikawr/gotp/internal/cli/ui"
	"github.com/zulfikawr/gotp/internal/config"
	"github.com/zulfikawr/gotp/internal/importers"
	"github.com/zulfikawr/gotp/internal/vault"
)

func NewImportCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "import <file>",
		Short: "Import accounts from a file",
		Long:  `Import TOTP accounts into your secure vault from a file. Supports Aegis, Authy, Google Authenticator, JSON, otpauth:// URIs, and password-protected encrypted exports.`,
		Args:  cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			vaultPath := config.GetVaultPath()

			// Check if vault exists first
			if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
				fmt.Fprintf(ui.Out, "%sError: Vault file not found at %s%s\n", ui.DangerBright, vaultPath, ui.Reset)
				fmt.Fprintf(ui.Out, "%sTip: Run '%s%sgotp %sinit%s' to create a new secure vault.%s\n", ui.TextMuted, ui.Reset, ui.SuccessBright, ui.WarningBright, ui.TextMuted, ui.Reset)
				return nil
			}

			v, key, err := vault.LoadVaultInteractive(vaultPath, ui.PromptPassword)
			if err != nil {
				fmt.Fprintf(ui.Out, "%sError: %v%s\n", ui.DangerBright, err, ui.Reset)
				return nil
			}

			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Fprintf(ui.Out, "%sError: Failed to read file: %v%s\n", ui.DangerBright, err, ui.Reset)
				return nil
			}

			var importedAccounts []vault.Account
			switch format {
			case "json":
				if err := json.Unmarshal(data, &importedAccounts); err != nil {
					fmt.Fprintf(ui.Out, "%sError: Failed to parse JSON: %v%s\n", ui.DangerBright, err, ui.Reset)
					return nil
				}
			case "uri":
				scanner := bufio.NewScanner(strings.NewReader(string(data)))
				for scanner.Scan() {
					line := strings.TrimSpace(scanner.Text())
					if line == "" {
						continue
					}
					// Check if it's a Google migration URI (otpauth-migration://)
					if strings.HasPrefix(line, "otpauth-migration://") {
						// Parse Google migration URI
						accounts, err := importers.ParseGoogleExport([]byte(line))
						if err != nil {
							fmt.Fprintf(ui.Out, "%sWarning: skipping invalid migration URI: %v%s\n", ui.WarningBright, err, ui.Reset)
							continue
						}
						importedAccounts = append(importedAccounts, accounts...)
						continue
					}
					// Regular otpauth:// URI
					acc, err := vault.FromURI(line)
					if err != nil {
						fmt.Fprintf(ui.Out, "%sWarning: skipping invalid URI: %v%s\n", ui.WarningBright, err, ui.Reset)
						continue
					}
					importedAccounts = append(importedAccounts, *acc)
				}
			case "encrypted":
				exportPass, err := ui.PromptPassword("Enter password for encrypted import: ")
				if err != nil {
					return err
				}

				var metadata vault.VaultMetadata
				if err := json.Unmarshal(data, &metadata); err != nil {
					fmt.Fprintf(ui.Out, "%sError: Failed to parse metadata: %v%s\n", ui.DangerBright, err, ui.Reset)
					return nil
				}

				impVault, err := vault.UnmarshalVault(metadata.Ciphertext, exportPass, metadata.Salt, metadata.KDFParams)
				if err != nil {
					fmt.Fprintf(ui.Out, "%sError: Import decryption failed: %v%s\n", ui.DangerBright, err, ui.Reset)
					return nil
				}
				importedAccounts = impVault.Accounts

			case "aegis":
				importedAccounts, err = importers.ImportData(data, importers.FormatAegis)
				if err != nil {
					fmt.Fprintf(ui.Out, "%sError: Aegis import failed: %v%s\n", ui.DangerBright, err, ui.Reset)
					return nil
				}

			case "authy":
				importedAccounts, err = importers.ImportData(data, importers.FormatAuthy)
				if err != nil {
					fmt.Fprintf(ui.Out, "%sError: Authy import failed: %v%s\n", ui.DangerBright, err, ui.Reset)
					return nil
				}

			case "google":
				importedAccounts, err = importers.ImportData(data, importers.FormatGoogle)
				if err != nil {
					fmt.Fprintf(ui.Out, "%sError: Google import failed: %v%s\n", ui.DangerBright, err, ui.Reset)
					return nil
				}

			case "auto":
				// Auto-detect format
				detectedFormat := importers.DetectFormat(data)
				fmt.Fprintf(ui.Out, "%sDetected format: %s%s\n", ui.InfoBright, detectedFormat, ui.Reset)
				importedAccounts, err = importers.ImportData(data, detectedFormat)
				if err != nil {
					fmt.Fprintf(ui.Out, "%sError: Import failed: %v%s\n", ui.DangerBright, err, ui.Reset)
					return nil
				}

			default:
				fmt.Fprintf(ui.Out, "%sError: Unsupported format: %s%s\n", ui.DangerBright, format, ui.Reset)
				return nil
			}

			count := 0
			skipped := 0
			for _, impAcc := range importedAccounts {
				isDuplicate := false
				for _, existing := range v.Accounts {
					if strings.EqualFold(existing.Name, impAcc.Name) &&
						strings.EqualFold(existing.Issuer, impAcc.Issuer) &&
						strings.EqualFold(existing.Username, impAcc.Username) {
						isDuplicate = true
						break
					}
				}

				if isDuplicate {
					skipped++
					continue
				}

				if impAcc.ID == "" {
					impAcc.ID = uuid.New().String()
				}
				v.Accounts = append(v.Accounts, impAcc)
				count++
			}

			if count > 0 {
				if err := vault.CreateBackup(vaultPath, 3); err != nil {
					fmt.Fprintf(ui.Out, "%sWarning: failed to create backup: %v%s\n", ui.WarningBright, err, ui.Reset)
				}

				if err := vault.SaveVaultWithKey(vaultPath, v, key); err != nil {
					fmt.Fprintf(ui.Out, "%sError: Failed to save vault: %v%s\n", ui.DangerBright, err, ui.Reset)
					return nil
				}
			}

			fmt.Fprintf(ui.Out, "%s✓ Imported %d accounts, skipped %d duplicates.%s\n", ui.SuccessBright, count, skipped, ui.Reset)
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "auto", "Import format (auto, json, uri, encrypted, aegis, authy, google)")

	return cmd
}

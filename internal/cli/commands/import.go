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
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			vaultPath := config.GetVaultPath()

			v, key, err := vault.LoadVaultInteractive(vaultPath, ui.PromptPassword)
			if err != nil {
				return err
			}

			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			var importedAccounts []vault.Account
			switch format {
			case "json":
				if err := json.Unmarshal(data, &importedAccounts); err != nil {
					return err
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
					return err
				}

				impVault, err := vault.UnmarshalVault(metadata.Ciphertext, exportPass, metadata.Salt, metadata.KDFParams)
				if err != nil {
					return err
				}
				importedAccounts = impVault.Accounts

			case "aegis":
				importedAccounts, err = importers.ImportData(data, importers.FormatAegis)
				if err != nil {
					return err
				}

			case "authy":
				importedAccounts, err = importers.ImportData(data, importers.FormatAuthy)
				if err != nil {
					return err
				}

			case "google":
				importedAccounts, err = importers.ImportData(data, importers.FormatGoogle)
				if err != nil {
					return err
				}

			case "auto":
				// Auto-detect format
				detectedFormat := importers.DetectFormat(data)
				fmt.Fprintf(ui.Out, "%sDetected format: %s%s\n", ui.InfoBright, detectedFormat, ui.Reset)
				importedAccounts, err = importers.ImportData(data, detectedFormat)
				if err != nil {
					return err
				}

			default:
				return fmt.Errorf("unsupported format: %s. Use 'json', 'uri', 'encrypted', 'aegis', 'authy', 'google', or 'auto'", format)
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
					return err
				}
			}

			fmt.Fprintf(ui.Out, "%s✓ Imported %d accounts, skipped %d duplicates.%s\n", ui.SuccessBright, count, skipped, ui.Reset)
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "auto", "Import format (auto, json, uri, encrypted, aegis, authy, google)")

	return cmd
}

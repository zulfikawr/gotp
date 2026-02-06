package commands

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zulfikawr/gotp/internal/cli/ui"
	"github.com/zulfikawr/gotp/internal/config"
	"github.com/zulfikawr/gotp/internal/totp"
	"github.com/zulfikawr/gotp/internal/vault"
	"github.com/zulfikawr/gotp/pkg/base32"
)

func NewListCmd() *cobra.Command {
	var filterTag string
	var sortBy string
	var withCodes bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all stored accounts",
		Long:  `Display all TOTP accounts stored in your secure vault. Supports filtering by tags and various sorting options.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath := config.GetVaultPath()
			isJSON, _ := cmd.Flags().GetBool("json")

			v, _, err := vault.LoadVaultInteractive(vaultPath, ui.PromptPassword)
			if err != nil {
				return err
			}

			accounts := v.Accounts

			if filterTag != "" {
				var filtered []vault.Account
				for _, acc := range accounts {
					for _, tag := range acc.Tags {
						if strings.EqualFold(tag, filterTag) {
							filtered = append(filtered, acc)
							break
						}
					}
				}
				accounts = filtered
			}

			sort.Slice(accounts, func(i, j int) bool {
				switch strings.ToLower(sortBy) {
				case "issuer":
					return accounts[i].Issuer < accounts[j].Issuer
				case "username":
					return accounts[i].Username < accounts[j].Username
				default:
					return accounts[i].Name < accounts[j].Name
				}
			})

			if isJSON {
				data, _ := json.Marshal(accounts)
				fmt.Fprintln(ui.Out, string(data))
				return nil
			}

			if len(accounts) == 0 {
				fmt.Fprintln(ui.Out, "No accounts found.")
				return nil
			}

			headers := []string{"NAME", "ISSUER", "USERNAME"}
			if withCodes {
				headers = append(headers, "CODE")
			}
			headers = append(headers, "TAGS")

			rows := [][]string{}
			now := time.Now()

			for _, acc := range accounts {
				row := []string{acc.Name, acc.Issuer, acc.Username}
				if withCodes {
					code := "ERROR"
					secretBytes, err := base32.Decode(string(acc.Secret))
					if err == nil {
						code, _ = totp.GenerateTOTP(totp.TOTPParams{
							Secret:    secretBytes,
							Timestamp: now,
							Period:    acc.Period,
							Digits:    acc.Digits,
							Algorithm: acc.Algorithm,
						})
					}
					row = append(row, code)
				}
				row = append(row, strings.Join(acc.Tags, ", "))
				rows = append(rows, row)
			}

			ui.PrintTable(headers, rows)
			fmt.Fprintf(ui.Out, "\nTotal: %d accounts\n", len(accounts))
			return nil
		},
	}

	cmd.Flags().StringVarP(&filterTag, "filter", "f", "", "Filter by tag")
	cmd.Flags().StringVar(&sortBy, "sort", "name", "Sort by (name, issuer, username)")
	cmd.Flags().BoolVar(&withCodes, "with-codes", false, "Show current TOTP codes")

	return cmd
}

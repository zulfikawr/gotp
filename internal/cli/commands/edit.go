package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zulfikawr/gotp/internal/cli/ui"
	"github.com/zulfikawr/gotp/internal/config"
	"github.com/zulfikawr/gotp/internal/vault"
	"github.com/zulfikawr/gotp/pkg/base32"
)

func NewEditCmd() *cobra.Command {
	var newName, username, issuer, secret string
	var tags, addTags, removeTags []string

	cmd := &cobra.Command{
		Use:   "edit <name>",
		Short: "Edit an existing account",
		Long:  `Modify the details of an existing TOTP account. You can use flags for specific changes or enter interactive mode for a guided experience.`,
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

			acc := &v.Accounts[index]

			flagsProvided := cmd.Flags().Changed("name") ||
				cmd.Flags().Changed("username") ||
				cmd.Flags().Changed("issuer") ||
				cmd.Flags().Changed("secret") ||
				cmd.Flags().Changed("tags") ||
				cmd.Flags().Changed("add-tag") ||
				cmd.Flags().Changed("remove-tag")

			if !flagsProvided {
				fmt.Fprintln(ui.Out, ui.Dimmed("No edit flags provided. Entering interactive mode..."))
				for {
					fmt.Fprintf(ui.Out, "\nEditing account: %s%s%s\n", ui.Bold, acc.Name, ui.Reset)
					fmt.Fprintln(ui.Out, "1. Name")
					fmt.Fprintln(ui.Out, "2. Issuer")
					fmt.Fprintln(ui.Out, "3. Username")
					fmt.Fprintln(ui.Out, "4. Secret")
					fmt.Fprintln(ui.Out, "5. Tags")
					fmt.Fprintln(ui.Out, "0. Save and Exit")

					choice := ui.PromptString("Select field to edit (0-5)", "0")

					switch choice {
					case "1":
						acc.Name = ui.PromptRequired("New Name")
					case "2":
						acc.Issuer = ui.PromptString("New Issuer", acc.Issuer)
					case "3":
						acc.Username = ui.PromptString("New Username", acc.Username)
					case "4":
						newSecret := ui.PromptValidate("New Secret (Base32)", func(s string) error {
							s = strings.ToUpper(strings.ReplaceAll(s, " ", ""))
							_, err := base32.Decode(s)
							return err
						})
						if ui.PromptConfirm("Are you sure you want to change the secret?", false) {
							acc.Secret = vault.Secret(strings.ToUpper(strings.ReplaceAll(newSecret, " ", "")))
						}
					case "5":
						currentTags := strings.Join(acc.Tags, ", ")
						newTagsStr := ui.PromptString("New Tags (comma-separated)", currentTags)
						acc.Tags = strings.Split(newTagsStr, ",")
						for i := range acc.Tags {
							acc.Tags[i] = strings.TrimSpace(acc.Tags[i])
						}
					case "0":
						goto save
					default:
						fmt.Fprintln(ui.Out, "Invalid choice.")
					}
				}
			} else {
				if newName != "" {
					acc.Name = newName
				}
				if username != "" {
					acc.Username = username
				}
				if issuer != "" {
					acc.Issuer = issuer
				}
				if secret != "" {
					if ui.PromptConfirm("Are you sure you want to change the secret?", false) {
						acc.Secret = vault.Secret(strings.ToUpper(strings.ReplaceAll(secret, " ", "")))
					}
				}
				if len(tags) > 0 {
					acc.Tags = tags
				}
				for _, t := range addTags {
					found := false
					for _, existing := range acc.Tags {
						if existing == t {
							found = true
							break
						}
					}
					if !found {
						acc.Tags = append(acc.Tags, t)
					}
				}
				for _, t := range removeTags {
					for i, existing := range acc.Tags {
						if existing == t {
							acc.Tags = append(acc.Tags[:i], acc.Tags[i+1:]...)
							break
						}
					}
				}
			}

		save:
			if err := vault.SaveVaultWithKey(vaultPath, v, key); err != nil {
				return err
			}

			fmt.Fprintf(ui.Out, "%s✓ Updated account: %s%s\n", ui.SuccessBright, acc.Name, ui.Reset)
			return nil
		},
	}

	cmd.Flags().StringVar(&newName, "name", "", "New account name")
	cmd.Flags().StringVar(&username, "username", "", "New username")
	cmd.Flags().StringVar(&issuer, "issuer", "", "New issuer")
	cmd.Flags().StringVar(&secret, "secret", "", "New secret")
	cmd.Flags().StringSliceVar(&tags, "tags", []string{}, "Replace tags")
	cmd.Flags().StringSliceVar(&addTags, "add-tag", []string{}, "Add tags")
	cmd.Flags().StringSliceVar(&removeTags, "remove-tag", []string{}, "Remove tags")

	return cmd
}

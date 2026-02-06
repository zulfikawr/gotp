package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zulfikawr/gotp/internal/cli/commands"
	"github.com/zulfikawr/gotp/internal/cli/ui"
	"github.com/zulfikawr/gotp/internal/config"
	"golang.org/x/term"
)

var (
	vaultPath  string
	configPath string
	jsonOutput bool
	noColor    bool
)

// wrapText wraps a single-line string into multiple lines based on terminal width.
func wrapText(text string, width int, indent string) string {
	if width < 10 {
		width = 80
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	currentLine := indent + words[0]
	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) > width {
			result.WriteString(currentLine + "\n")
			currentLine = indent + word
		} else {
			currentLine += " " + word
		}
	}
	result.WriteString(currentLine)
	return result.String()
}

// NewRootCmd creates the root command for the gotp CLI.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gotp",
		Short: "gotp - A terminal-based TOTP authenticator",
		Long:  `gotp is a secure, cross-platform, terminal-based TOTP authenticator that allows you to manage your two-factor authentication codes.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if vaultPath != "" {
				config.SetVaultPathOverride(vaultPath)
			}
			if noColor {
				ui.SetColor(false)
			}
		},
	}

	// Custom Help Logic
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		out := cmd.OutOrStdout()
		width, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			width = 80
		}

		// 1. App Name and Wrapped Description
		fmt.Fprintf(out, "%s%s%s\n", ui.PrimaryBright+ui.Bold, cmd.Name(), ui.Reset)
		fmt.Fprintf(out, "%s%s%s\n\n", ui.TextMuted, wrapText(cmd.Long, width, "  "), ui.Reset)

		// 2. Usage
		fmt.Fprintf(out, "%sUsage:%s\n", ui.PrimaryBright+ui.Bold, ui.Reset)

		commandPath := cmd.CommandPath()
		parts := strings.Split(commandPath, " ")
		rootName := parts[0]
		subPath := ""
		if len(parts) > 1 {
			subPath = " " + ui.WarningBright + strings.Join(parts[1:], " ") + ui.Reset
		}

		usageStr := "  " + ui.SuccessBright + rootName + ui.Reset + subPath

		if cmd.HasAvailableSubCommands() {
			usageStr += " " + ui.WarningBright + "[command]" + ui.Reset
		}
		if cmd.HasAvailableLocalFlags() {
			usageStr += " " + ui.InfoBright + "[flags]" + ui.Reset
		}
		if cmd.HasAvailableInheritedFlags() {
			usageStr += " " + ui.InfoBright + "[global flags]" + ui.Reset
		}
		fmt.Fprintln(out, usageStr)
		fmt.Fprintln(out)
		// 3. Commands
		if cmd.HasAvailableSubCommands() {
			fmt.Fprintf(out, "%sAvailable Commands:%s\n", ui.PrimaryBright+ui.Bold, ui.Reset)
			for _, sub := range cmd.Commands() {
				if sub.IsAvailableCommand() || sub.Name() == "help" {
					fmt.Fprintf(out, "  %s%-12s%s %s%s%s\n", ui.WarningBright, sub.Name(), ui.Reset, ui.TextMuted, sub.Short, ui.Reset)
				}
			}
			if cmd.LocalFlags().HasFlags() || cmd.InheritedFlags().HasFlags() {
				fmt.Fprintln(out)
			}
		}

		// 4. Flags
		printFlags := func(title string, fs *pflag.FlagSet, hasNext bool) {
			if fs.HasFlags() {
				fmt.Fprintf(out, "%s%s%s\n", ui.PrimaryBright+ui.Bold, title, ui.Reset)
				fs.VisitAll(func(f *pflag.Flag) {
					if f.Hidden {
						return
					}
					shorthand := ""
					if f.Shorthand != "" {
						shorthand = "-" + f.Shorthand + ", "
					}
					fmt.Fprintf(out, "  %s%-4s--%-12s%s %s%s%s\n", ui.InfoBright, shorthand, f.Name, ui.Reset, ui.TextMuted, f.Usage, ui.Reset)
				})
				if hasNext {
					fmt.Fprintln(out)
				}
			}
		}

		printFlags("Flags:", cmd.LocalFlags(), cmd.InheritedFlags().HasFlags())
		printFlags("Global Flags:", cmd.InheritedFlags(), cmd.Parent() == nil)

		// 5. Footer (Only for root command)
		if cmd.HasAvailableSubCommands() && cmd.Parent() == nil {
			fmt.Fprintf(out, "\n%sUse \"%s%s%s %s[command]%s %s--help%s%s\" for more information about a command.%s\n",
				ui.TextMuted,
				ui.SuccessBright, cmd.Root().Name(), ui.Reset,
				ui.WarningBright, ui.Reset,
				ui.InfoBright, ui.Reset,
				ui.TextMuted,
				ui.Reset)
		}
	})

	// Persistent Flags (Global)
	rootCmd.PersistentFlags().StringVarP(&vaultPath, "vault", "v", "", "Path to vault file")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "C", "", "Path to config file")
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")

	// Register subcommands
	rootCmd.AddCommand(commands.NewInitCmd())
	rootCmd.AddCommand(commands.NewAddCmd())
	rootCmd.AddCommand(commands.NewListCmd())
	rootCmd.AddCommand(commands.NewGetCmd())
	rootCmd.AddCommand(commands.NewRemoveCmd())
	rootCmd.AddCommand(commands.NewEditCmd())
	rootCmd.AddCommand(commands.NewExportCmd())
	rootCmd.AddCommand(commands.NewImportCmd())
	rootCmd.AddCommand(commands.NewPasswdCmd())
	rootCmd.AddCommand(commands.NewQrCmd())
	rootCmd.AddCommand(commands.NewCompletionCmd())

	return rootCmd
}

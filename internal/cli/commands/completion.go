package commands

import (
	"os"

	"github.com/spf13/cobra"
)

func NewCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

  $ source <(gotp completion bash)

  # To load completions for each session, add to your .bashrc:
  # gotp completion bash > /etc/bash_completion.d/gotp

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, add to your .zshrc:
  $ gotp completion zsh > "${fpath[1]}/_gotp"

Fish:

  $ gotp completion fish | source

  # To load completions for each session, add to your fish config:
  $ gotp completion fish > ~/.config/fish/completions/gotp.fish

PowerShell:

  PS> gotp completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> gotp completion powershell > gotp.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			switch args[0] {
			case "bash":
				err = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				err = cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				err = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				err = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			if err != nil {
				os.Exit(1)
			}
		},
	}
}

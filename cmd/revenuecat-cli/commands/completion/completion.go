package completion

import (
	"os"

	"github.com/spf13/cobra"
)

// CompletionCmd generates shell completion scripts
var CompletionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for rc.

To load completions:

Bash:
  $ source <(rc completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ rc completion bash > /etc/bash_completion.d/rc
  # macOS:
  $ rc completion bash > $(brew --prefix)/etc/bash_completion.d/rc

Zsh:
  $ source <(rc completion zsh)
  # To load completions for each session, execute once:
  $ rc completion zsh > "${fpath[1]}/_rc"

Fish:
  $ rc completion fish | source
  # To load completions for each session, execute once:
  $ rc completion fish > ~/.config/fish/completions/rc.fish

PowerShell:
  PS> rc completion powershell | Out-String | Invoke-Expression
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

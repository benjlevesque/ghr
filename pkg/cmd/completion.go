package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func makeCompletionCmd() *cobra.Command {
	var completionCmd = &cobra.Command{
		Use: "completion [bash|zsh|fish|powershell]",
		Long: `To load completions:

Bash:

$ source <(ghr completion bash)

# To load completions for each session, execute once:
Linux:
  $ ghr completion bash > /etc/bash_completion.d/ghr
MacOS:
  $ ghr completion bash > /usr/local/etc/bash_completion.d/ghr

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ ghr completion zsh > "${fpath[1]}/_ghr"

# You will need to start a new shell for this setup to take effect.

Fish:

$ ghr completion fish | source

# To load completions for each session, execute once:
$ ghr completion fish > ~/.config/fish/completions/ghr.fish
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		},
	}
	return completionCmd
}

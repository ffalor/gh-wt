package cmd

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/ffalor/gh-wt/internal/completion"
	"github.com/spf13/cobra"
)

// NewCompletionCommand creates the completion command with install/uninstall subcommands.
func NewCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [shell]",
		Short: "Generate shell completion scripts for gh wt commands",
		Long: heredoc.Doc(`
			Generate shell completion scripts to enable tab completion for gh wt commands.

			Tab completion provides:
			- Command name completion (add, remove, run, action)
			- Subcommand completion (install, uninstall)
			- Flag completion

			Supported shells: bash, zsh, fish, powershell
		`),
		Example: heredoc.Doc(`
			# Generate completion script for bash
			gh wt completion bash

			# Generate completion script for zsh
			gh wt completion zsh

			# Generate completion script for fish
			gh wt completion fish

			# Generate completion script for PowerShell
			gh wt completion powershell

			# Install completions automatically (detects your shell)
			gh wt completion install

			# Uninstall completions
			gh wt completion uninstall
		`),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := args[0]

			switch shell {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletion(os.Stdout)
			default:
				return cmd.Help()
			}
		},
	}

	// Add install subcommand
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install shell completion for the detected shell",
		Long: heredoc.Doc(`
			Automatically install shell completion for your current shell.

			This command detects your shell (bash, zsh, fish, or PowerShell) and installs
			the completion script to the appropriate location. After installation, restart
			your shell or source your shell configuration file.

			Supported shells:
			  - Bash: Installs to ~/.bash_completion.d/ or /etc/bash_completion.d/
			  - Zsh: Installs to ~/.zsh/completions/
			  - Fish: Installs to ~/.config/fish/completions/
			  - PowerShell: Provides instructions to add to profile
		`),
		Example: heredoc.Doc(`
			gh wt completion install
			gh wt completion install --verbose
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return completion.InstallShellCompletion(Log, cmd.Root())
		},
	}

	// Add uninstall subcommand
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall shell completion for the detected shell",
		Long: heredoc.Doc(`
			Automatically uninstall shell completion for your current shell.

			This command detects your shell and removes the completion script from the
			appropriate location. After uninstallation, restart your shell or source
			your shell configuration file.
		`),
		Example: heredoc.Doc(`
			gh wt completion uninstall
			gh wt completion uninstall --verbose
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return completion.UninstallShellCompletion(Log)
		},
	}

	cmd.AddCommand(installCmd)
	cmd.AddCommand(uninstallCmd)

	return cmd
}

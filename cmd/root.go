package cmd

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/ffalor/gh-wt/internal/config"
	"github.com/ffalor/gh-wt/internal/logger"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	forceFlag bool
	verbose   bool
	noColor   bool
	cliArgs   string
)

// Version is the current version of the CLI.
var Version = "dev"

// Commit is the git commit hash.
var Commit = ""

// Date is the build date.
var Date = ""

// BuiltBy is the builder.
var BuiltBy = ""

func buildVersion(version, commit, date, builtBy string) string {
	result := version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	result = fmt.Sprintf("%s\ngoos: %s\ngoarch: %s", result, runtime.GOOS, runtime.GOARCH)
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
		result = fmt.Sprintf("%s\nmodule version: %s, checksum: %s", result, info.Main.Version, info.Main.Sum)
	}

	return result
}

// Log is the package-level logger instance.
var Log = logger.NewLogger(false, true)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use: "gh wt",
	Annotations: map[string]string{
		cobra.CommandDisplayNameAnnotation: "gh wt",
	},
	Short: "Create and manage git worktrees",
	Long: `gh wt is a GitHub CLI extension that helps you create git worktrees. A GitHub pull request or issue url can also be used.

Examples:
  # Create worktree from PR URL
  gh wt add https://github.com/owner/repo/pull/123 -action claude -- "/review"

  # Create worktree from Issue URL
  gh wt add https://github.com/owner/repo/issues/456 -action claude -- "implement issue #456"

  # Create a worktree
  gh wt add my-feature-branch

  # Remove a worktree
  gh wt rm pr_123`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.Load()
		if err != nil {
			return err
		}
		Log = logger.NewLogger(verbose, !noColor)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Find and store arguments after --
	dashDashIndex := -1
	for i, arg := range os.Args {
		if arg == "--" {
			dashDashIndex = i
			break
		}
	}

	if dashDashIndex != -1 {
		cliArgs = strings.Join(os.Args[dashDashIndex+1:], " ")
		os.Args = os.Args[:dashDashIndex]
	}

	err := rootCmd.Execute()
	if err != nil {
		if Log != nil {
			Log.Errorf("Error: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&forceFlag, "force", "f", false, "force operation without prompts")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color output")

	// Version flag
	rootCmd.Version = buildVersion(Version, Commit, Date, BuiltBy)
	rootCmd.SetVersionTemplate(`gh-wt {{printf "version %s\n" .Version}}`)
}

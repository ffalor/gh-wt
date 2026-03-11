package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/ffalor/gh-wt/internal/config"
	"github.com/ffalor/gh-wt/internal/git"
	"github.com/ffalor/gh-wt/internal/logger"
	"github.com/spf13/cobra"
)

// listCmd represents the list command.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List managed worktrees",
	Long: heredoc.Doc(`
		List all worktrees managed by gh-wt (those under the configured worktree directory).
		Displays the worktree name and associated branch.
	`),
	Example: heredoc.Doc(`
		# List all worktrees
		gh wt list

		# Using the alias
		gh wt ls
	`),
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	GroupID: "worktrees",
	RunE:    runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	worktrees, err := git.GetWorktreeInfo()
	if err != nil {
		return fmt.Errorf("failed to list worktrees: %w", err)
	}

	filtered := filterWorktreesByBase(worktrees, cfg.WorktreeBase)

	if len(filtered) == 0 {
		Log.Warnf("No worktrees found under %s\n", cfg.WorktreeBase)
		return nil
	}

	// Build entries and compute max name width
	type entry struct{ name, branch string }
	entries := make([]entry, 0, len(filtered))
	maxWidth := len("NAME")
	for _, wt := range filtered {
		name := getWorktreeDisplayName(wt.Path)
		branch := wt.Branch
		if branch == "" {
			branch = "(detached)"
		}
		if len(name) > maxWidth {
			maxWidth = len(name)
		}
		entries = append(entries, entry{name, branch})
	}

	// Header
	Log.Outf(logger.Default, "%-*s%s\n", maxWidth+4, "NAME", "BRANCH")

	// Rows
	for _, e := range entries {
		Log.Outf(logger.Green, "%-*s", maxWidth+4, e.name)
		Log.Outf(logger.Default, "%s\n", e.branch)
	}

	return nil
}

// filterWorktreesByBase returns only worktrees under the configured base directory.
func filterWorktreesByBase(worktrees []git.WorktreeInfo, base string) []git.WorktreeInfo {
	var filtered []git.WorktreeInfo
	prefix := base + string(os.PathSeparator)

	for _, wt := range worktrees {
		if strings.HasPrefix(wt.Path, prefix) {
			filtered = append(filtered, wt)
		}
	}

	return filtered
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/ffalor/gh-wt/internal/config"
	"github.com/ffalor/gh-wt/internal/git"
	"github.com/ffalor/gh-wt/internal/logger"
	"github.com/spf13/cobra"
)

var allFlag bool

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

		# List worktrees across all repos
		gh wt list --all

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
	listCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "list worktrees for all repos")
}

func runList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if allFlag {
		return runListAll(cfg)
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

func runListAll(cfg config.Config) error {
	worktrees, err := git.ListAllWorktrees(cfg.WorktreeBase)
	if err != nil {
		return fmt.Errorf("failed to list all worktrees: %w", err)
	}

	if len(worktrees) == 0 {
		Log.Warnf("No worktrees found under %s\n", cfg.WorktreeBase)
		return nil
	}

	groups := groupWorktreesByRepo(worktrees, cfg.WorktreeBase)

	// Compute a global max name width across all groups for consistent column alignment
	maxWidth := len("NAME")
	for _, wt := range worktrees {
		name := filepath.Base(wt.Path)
		if len(name) > maxWidth {
			maxWidth = len(name)
		}
	}

	for i, group := range groups {
		if i > 0 {
			Log.Plainf("\n")
		}

		// Repo name header
		Log.Outf(logger.Default, "%s\n", group.repo)

		// Indented header
		Log.Outf(logger.Default, "  %-*s%s\n", maxWidth+4, "NAME", "BRANCH")

		// Indented rows
		for _, wt := range group.worktrees {
			name := filepath.Base(wt.Path)
			branch := wt.Branch
			if branch == "" {
				branch = "(detached)"
			}
			Log.Plainf("  ")
			Log.Outf(logger.Green, "%-*s", maxWidth+4, name)
			Log.Outf(logger.Default, "%s\n", branch)
		}
	}

	return nil
}

type repoGroup struct {
	repo      string
	worktrees []git.WorktreeInfo
}

// groupWorktreesByRepo groups worktrees by their parent repo directory,
// preserving discovery order.
func groupWorktreesByRepo(worktrees []git.WorktreeInfo, baseDir string) []repoGroup {
	orderMap := make(map[string]int)
	var groups []repoGroup

	for _, wt := range worktrees {
		// Extract the repo name: path relative to baseDir, then take the first component
		rel, err := filepath.Rel(baseDir, wt.Path)
		if err != nil {
			continue
		}
		parts := strings.SplitN(filepath.ToSlash(rel), "/", 2)
		repo := parts[0]

		idx, exists := orderMap[repo]
		if !exists {
			idx = len(groups)
			orderMap[repo] = idx
			groups = append(groups, repoGroup{repo: repo})
		}
		groups[idx].worktrees = append(groups[idx].worktrees, wt)
	}

	return groups
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

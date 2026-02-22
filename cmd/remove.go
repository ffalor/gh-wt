package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/ffalor/gh-wt/internal/git"
	"github.com/ffalor/gh-wt/internal/logger"
	"github.com/ffalor/gh-wt/internal/worktree"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command.
var rmCmd = &cobra.Command{
	Use:   "rm [worktree-name]",
	Short: "Remove a worktree and its associated branch",
	Long: heredoc.Doc(`
		Remove a worktree and its associated branch. Will prompt if there are
		uncommitted changes (unless --force is used).
	`),
	Example: heredoc.Doc(`
		# Remove a worktree by name
		gh wt rm pr_123

		# Remove a worktree with force
		gh wt rm issue_456 --force
	`),
	Aliases: []string{"remove"},
	Args:    cobra.ExactArgs(1),
	RunE:    runRm,
	GroupID: "worktrees",
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

func runRm(cmd *cobra.Command, args []string) error {
	worktreeName := args[0]

	// Require being in a git repository (consistent with create command)
	if !git.IsGitRepository(".") {
		return fmt.Errorf("not in a git repository")
	}

	// Find the worktree by name using the shared helper
	matches, err := worktree.FindByName(worktreeName)
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		Log.Warnf("Worktree '%s' not found in this repository.\n", worktreeName)
		return nil
	}

	// If multiple matches, prompt user to select one
	var targetWorktree git.WorktreeInfo
	if len(matches) == 1 {
		targetWorktree = matches[0]
	} else {
		options := make([]string, len(matches))
		for i, wt := range matches {
			options[i] = wt.Path
		}
		p := prompter.New(os.Stdin, os.Stdout, os.Stderr)
		idx, err := p.Select("Multiple worktrees match '"+worktreeName+"'. Select one:", "", options)
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}
		targetWorktree = matches[idx]
	}

	// Handle uncommitted changes prompt.
	force := forceFlag
	if !force && git.HasUncommittedChanges(targetWorktree.Path) {
		p := prompter.New(os.Stdin, os.Stdout, os.Stderr)
		confirm, err := p.Confirm("Worktree has uncommitted changes. Remove anyway?", false)
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}
		if !confirm {
			Log.Warnf("Cancelled - no changes made\n")
			return nil
		}
		force = true // User confirmed.
	}

	// Extract the worktree name from the path for display
	worktreeDisplayName := getWorktreeDisplayName(targetWorktree.Path)
	worktreePathDisplay := getTildePath(targetWorktree.Path)

	// Print the header line
	Log.Infof("Removing worktree %s...\n", worktreeDisplayName)

	// 1. Remove the worktree directory and git metadata.
	if err := worktree.Remove(targetWorktree.Path, force); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	Log.Outf(logger.Default, "Worktree: %s\n", worktreePathDisplay)

	if targetWorktree.Branch != "" {
		Log.Outf(logger.Default, "Branch: %s\n", targetWorktree.Branch)
	} else {
		Log.Outf(logger.Default, "Branch: <none>\n")
	}

	// 2. Delete the associated branch if we found one.
	if targetWorktree.Branch != "" {
		if err := git.BranchDelete(targetWorktree.Branch, true); err != nil {
			// This is not a fatal error, as the primary goal (removing the worktree) succeeded.
			// The branch might be the main branch or have other worktrees, so git will prevent its deletion.
			Log.Warnf("Failed to delete branch '%s': %v. You may need to remove it manually.\n", targetWorktree.Branch, err)
		}
	}

	// Print the details and success message

	Log.Outf(logger.Green, "✓ Worktree removed successfully!\n")

	return nil
}

// getWorktreeDisplayName extracts a short name from the worktree path for display.
func getWorktreeDisplayName(path string) string {
	// Get the last two components of the path (repo/worktree-name)
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2] + "/" + parts[len(parts)-1]
	}
	return filepath.Base(path)
}

// getTildePath replaces the home directory with ~ for display.
func getTildePath(path string) string {
	home, err := os.UserHomeDir()
	if err == nil && strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}

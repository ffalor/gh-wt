package cmd

import (
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/ffalor/gh-wt/internal/git"
	"github.com/ffalor/gh-wt/internal/logger"
	"github.com/ffalor/gh-wt/internal/worktree"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command.
var rmCmd = &cobra.Command{
	Use:     "rm [worktree-name]",
	Short:   "Remove a worktree and its associated branch",
	Long:    `Remove a worktree and its associated branch. Will prompt if there are uncommitted changes (unless --force is used).`,
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

	// 1. Remove the worktree directory and git metadata.
	Log.Infof("Removing worktree '%s'...\n", targetWorktree.Path)
	if err := worktree.Remove(targetWorktree.Path, force); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}
	Log.Outf(logger.Green, "Successfully removed worktree directory.\n")

	// 2. Delete the associated branch if we found one.
	if targetWorktree.Branch != "" {
		Log.Infof("Deleting branch '%s'...\n", targetWorktree.Branch)
		if err := git.BranchDelete(targetWorktree.Branch, true); err != nil {
			// This is not a fatal error, as the primary goal (removing the worktree) succeeded.
			// The branch might be the main branch or have other worktrees, so git will prevent its deletion.
			return fmt.Errorf("worktree removed, but failed to delete branch '%s': %w. You may need to remove it manually", targetWorktree.Branch, err)
		}
		Log.Outf(logger.Green, "Successfully deleted branch '%s'.\n", targetWorktree.Branch)
	}

	Log.Outf(logger.Green, "\nWorktree '%s' and branch '%s' removed successfully.\n", targetWorktree.Path, targetWorktree.Branch)
	return nil
}

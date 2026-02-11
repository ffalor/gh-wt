package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/ffalor/gh-worktree/internal/config"
	"github.com/ffalor/gh-worktree/internal/git"
	"github.com/ffalor/gh-worktree/internal/worktree"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove [worktree-name]",
	Short: "Remove a worktree",
	Long:  `Remove a worktree and its associated branch. Will prompt if there are uncommitted changes (unless --force is used).`,
	Args:  cobra.ExactArgs(1),
	RunE:  runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	worktreeName := args[0]
	baseDir := config.GetWorktreeBase()
	
	// Find the worktree across all repos
	repoName, worktreePath, err := findWorktree(baseDir, worktreeName)
	if err != nil {
		return err
	}
	
	repoPath := filepath.Join(baseDir, repoName, ".base")
	
	// Get branch name
	branch, err := git.GetCurrentBranch(worktreePath)
	if err != nil {
		return fmt.Errorf("failed to get branch name: %w", err)
	}
	
	// Check for uncommitted changes if not forced
	if !forceFlag && git.HasUncommittedChanges(worktreePath) {
		var confirm bool
		err := huh.NewConfirm().
			Title(fmt.Sprintf("Worktree '%s' has uncommitted changes. Remove anyway?", worktreeName)).
			Affirmative("Yes, remove").
			Negative("No, keep it").
			Value(&confirm).
			Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}
		if !confirm {
			fmt.Println("Operation cancelled")
			return nil
		}
		// Force remove since user confirmed
		forceFlag = true
	}
	
	// Remove worktree
	if err := worktree.Remove(repoPath, worktreePath, branch, forceFlag); err != nil {
		return err
	}
	
	fmt.Printf("Removed worktree: %s\n", worktreeName)
	return nil
}

func findWorktree(baseDir, name string) (repoName, worktreePath string, err error) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return "", "", fmt.Errorf("worktree base directory not found: %w", err)
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		possiblePath := filepath.Join(baseDir, entry.Name(), name)
		if _, err := os.Stat(possiblePath); err == nil {
			return entry.Name(), possiblePath, nil
		}
	}
	
	return "", "", fmt.Errorf("worktree '%s' not found", name)
}

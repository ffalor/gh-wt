package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/ffalor/gh-wt/internal/action"
	"github.com/ffalor/gh-wt/internal/execext"
	"github.com/ffalor/gh-wt/internal/git"
	"github.com/ffalor/gh-wt/internal/logger"
	"github.com/ffalor/gh-wt/internal/worktree"
	"github.com/spf13/cobra"
)

// runCmd represents the run command.
var runCmd = &cobra.Command{
	Use:   "run <worktree> [action] [-- command]",
	Short: "Run an action or command in an existing worktree",
	Long: `Run an action or command in an existing worktree.

Use this command to:
- Run configured actions on worktrees that were created without an action
- Run commands directly in a worktree

Examples:
  # Run named action on worktree
  gh wt run pr_123 claude -- fix issue #456

  # Run command directly in worktree
  gh wt run pr_123 -- ls

  # Show help
  gh wt run pr_123`,
	Args:    cobra.RangeArgs(1, 2),
	RunE:    runRun,
	GroupID: "worktrees",
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// runRun is the main function for the run command.
func runRun(cmd *cobra.Command, args []string) error {
	worktreeName := args[0]
	var actionName string

	// Determine if we have an action name or just CLI args
	if len(args) > 1 {
		actionName = args[1]
	}

	// Find the worktree path
	wt, err := findWorktree(worktreeName)
	if err != nil {
		return err
	}

	// Check if worktree exists
	if !worktree.Exists(wt.Path) {
		return fmt.Errorf("worktree '%s' does not exist at %s", worktreeName, wt.Path)
	}

	info := &worktree.WorktreeInfo{
		WorktreeName: worktreeName,
		BranchName:   wt.Branch,
	}

	// Get repo name for worktree info - try GitHub API first, fallback to cwd
	repo, err := repository.Current()
	if err != nil {
		// Fallback to using current working directory name
		repoName, err := git.GetRepoName()
		if err != nil {
			return err
		}
		info.Repo = repoName
	} else {
		info.Repo = repo.Name
		info.Owner = repo.Owner
	}

	if actionName != "" {
		// Run the action
		Log.Outf(logger.Magenta, "Running action '%s' in %s...\n", actionName, wt.Path)

		if err := action.Execute(context.Background(), &action.ExecuteOptions{
			ActionName:   actionName,
			WorktreePath: wt.Path,
			Info:         info,
			CLIArgs:      cliArgs,
			Logger:       Log,
			Stdin:        os.Stdin,
			Stdout:       os.Stdout,
			Stderr:       os.Stderr,
			Env:          os.Environ(),
		}); err != nil {
			return fmt.Errorf("action '%s' failed: %w", actionName, err)
		}

		Log.Outf(logger.Green, "Action completed successfully.\n")
	} else if cliArgs != "" {
		// Run CLI args directly in the worktree
		Log.Outf(logger.Magenta, "Running in worktree: %s\n", cliArgs)

		if err := execext.RunCommand(context.Background(), &execext.RunCommandOptions{
			Command: cliArgs,
			Dir:     wt.Path,
			Env:     os.Environ(),
			Stdin:   os.Stdin,
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
		}); err != nil {
			return fmt.Errorf("command '%s' failed: %w", cliArgs, err)
		}
	} else {
		// No action or command provided, show help
		return cmd.Help()
	}

	return nil
}

// findWorktree finds the worktree based on the worktree name.
// It prompts if multiple matches.
func findWorktree(worktreeName string) (git.WorktreeInfo, error) {
	var info git.WorktreeInfo
	matches, err := worktree.FindByName(worktreeName)
	if err != nil {
		return info, err
	}

	if len(matches) == 0 {
		return info, fmt.Errorf("worktree '%s' not found", worktreeName)
	}

	// If multiple matches, prompt user to select one
	if len(matches) > 1 {
		options := make([]string, len(matches))
		for i, wt := range matches {
			options[i] = wt.Path
		}
		p := prompter.New(os.Stdin, os.Stdout, os.Stderr)
		idx, err := p.Select("Multiple worktrees match '"+worktreeName+"'. Select one:", "", options)
		if err != nil {
			return info, fmt.Errorf("prompt failed: %w", err)
		}
		return matches[idx], nil
	}

	// Return the single match
	return matches[0], nil
}

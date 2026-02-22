package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MakeNowJust/heredoc"
	gh "github.com/cli/go-gh/v2"
	"github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/ffalor/gh-wt/internal/action"
	"github.com/ffalor/gh-wt/internal/config"
	"github.com/ffalor/gh-wt/internal/execext"
	"github.com/ffalor/gh-wt/internal/git"
	"github.com/ffalor/gh-wt/internal/logger"
	"github.com/ffalor/gh-wt/internal/worktree"
	"github.com/spf13/cobra"
)

// addCmd represents the add command.
var addCmd = &cobra.Command{
	Use:   "add [url|name]",
	Short: "Add a new worktree",
	Long: heredoc.Doc(`
		Add a new git worktree from either:
		  - A GitHub pull request URL or number
		  - A GitHub issue URL or number
		  - A name to use for the new worktree and branch
	`),
	Example: heredoc.Doc(`
		# Create worktree from PR URL
		gh wt add https://github.com/owner/repo/pull/123

		# Create worktree from Issue URL
		gh wt add https://github.com/owner/repo/issues/456

		# Create a worktree from a local branch
		gh wt add my-feature-branch

		# Create worktree with custom name
		gh wt add https://github.com/owner/repo/pull/123 --name my-custom-name
	`),
	Aliases: []string{"create"},
	Args:    cobra.RangeArgs(0, 1),
	RunE:    runAdd,
	GroupID: "worktrees",
}

func init() {
	addCmd.Flags().StringVar(&prFlag, "pr", "", "PR number, PR URL, or git remote URL with PR ref")
	addCmd.Flags().StringVar(&issueFlag, "issue", "", "issue number, issue URL, or git remote URL with issue ref")
	addCmd.Flags().StringVarP(&nameFlag, "name", "n", "", "name to use for the worktree (overrides default for PR/Issue)")
	addCmd.Flags().StringVarP(&actionFlag, "action", "a", "", "action to run after worktree creation")
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	// Determine the type of input
	if prFlag != "" {
		return createFromPR(prFlag)
	}
	if issueFlag != "" {
		return createFromIssue(issueFlag)
	}
	if len(args) == 0 {
		return cmd.Help()
	}

	// This is the main entry point for creating a worktree
	arg := args[0]
	worktreeType, err := DetermineWorktreeType(arg)
	if err != nil {
		return err
	}

	switch worktreeType {
	case worktree.PR:
		return createFromPR(arg)
	case worktree.Issue:
		return createFromIssue(arg)
	default:
		return createFromLocal(arg)
	}
}

// createFromPR handles creation from a PR URL or number.
func createFromPR(value string) error {
	Log.Infof("Fetching Pull Request info...\n")
	args := []string{"pr", "view", value, "--json", "number,title,headRefName,url"}
	stdout, stderr, err := gh.Exec(args...)
	if err != nil {
		return fmt.Errorf("failed to fetch PR info: %w\n%s", err, stderr.String())
	}

	var prInfo struct {
		Number      int    `json:"number"`
		Title       string `json:"title"`
		HeadRefName string `json:"headRefName"`
		URL         string `json:"url"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &prInfo); err != nil {
		return fmt.Errorf("failed to parse PR info: %w", err)
	}

	repo, err := repository.Current()
	if err != nil {
		return err
	}

	worktreeName := fmt.Sprintf("pr_%d", prInfo.Number)
	if nameFlag != "" {
		worktreeName = nameFlag
	}

	info := &worktree.WorktreeInfo{
		Type:         worktree.PR,
		Owner:        repo.Owner,
		Repo:         repo.Name,
		Number:       prInfo.Number,
		BranchName:   prInfo.HeadRefName,
		WorktreeName: worktreeName,
	}

	Log.Outf(logger.Green, "Creating worktree for PR #%d: %s\n", info.Number, prInfo.Title)

	// Fetch the PR ref
	prRef := fmt.Sprintf("refs/pull/%d/head", info.Number)
	Log.Infof("Fetching PR #%d...\n", info.Number)
	if err := git.Fetch(prRef); err != nil {
		return fmt.Errorf("failed to fetch PR: %w", err)
	}

	return createWorktree(info, "FETCH_HEAD")
}

// createFromIssue handles creation from an Issue URL or number.
func createFromIssue(value string) error {
	Log.Infof("Fetching Issue info...\n")
	args := []string{"issue", "view", value, "--json", "number,title,url"}
	stdout, stderr, err := gh.Exec(args...)
	if err != nil {
		return fmt.Errorf("failed to fetch Issue info: %w\n%s", err, stderr.String())
	}

	var issueInfo struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		URL    string `json:"url"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &issueInfo); err != nil {
		return fmt.Errorf("failed to parse issue info: %w", err)
	}

	repo, err := repository.Current()
	if err != nil {
		return err
	}

	branchName := fmt.Sprintf("issue_%d", issueInfo.Number)
	worktreeName := branchName
	if nameFlag != "" {
		worktreeName = nameFlag
	}

	info := &worktree.WorktreeInfo{
		Type:         worktree.Issue,
		Owner:        repo.Owner,
		Repo:         repo.Name,
		Number:       issueInfo.Number,
		BranchName:   branchName,
		WorktreeName: worktreeName,
	}

	Log.Outf(logger.Green, "Creating worktree for Issue #%d: %s\n", info.Number, issueInfo.Title)
	return createWorktree(info, "HEAD") // Issues start from HEAD
}

// createFromLocal handles creation from a local branch name.
func createFromLocal(name string) error {
	if !git.IsGitRepository(".") {
		return fmt.Errorf("not in a git repository")
	}

	// Get repo name using the shared helper
	repoName, err := git.GetRepoName()
	if err != nil {
		return err
	}

	// Sanitize the name for the branch
	sanitizedBranchName := SanitizeBranchName(name)

	worktreeName := name
	if nameFlag != "" {
		worktreeName = nameFlag
		sanitizedBranchName = SanitizeBranchName(nameFlag)
	}

	info := &worktree.WorktreeInfo{
		Type:         worktree.Local,
		Repo:         repoName,
		BranchName:   sanitizedBranchName,
		WorktreeName: worktreeName,
	}

	return createWorktree(info, "HEAD")
}

func createWorktree(info *worktree.WorktreeInfo, startPoint string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	baseDir := cfg.WorktreeBase
	worktreePath := filepath.Join(baseDir, info.Repo, info.WorktreeName)
	absPath, _ := filepath.Abs(worktreePath)

	branchExists := git.BranchExists(info.BranchName)
	worktreeDirExists := worktree.Exists(worktreePath)
	worktreeGitRegistered := git.WorktreeIsRegistered(worktreePath)

	hasConflict := worktreeDirExists || worktreeGitRegistered || branchExists

	if hasConflict {
		if !forceFlag {
			message := buildConflictMessage(info, absPath, worktreePath, worktreeDirExists, worktreeGitRegistered, branchExists)
			p := prompter.New(os.Stdin, os.Stdout, os.Stderr)
			overwrite, err := p.Confirm(message, false)
			if err != nil {
				return fmt.Errorf("failed to read confirmation: %w", err)
			}
			if !overwrite {
				Log.Warnf("Cancelled - no changes made\n")
				return nil
			}
		}

		if err := performCleanup(worktreePath, worktreeDirExists, worktreeGitRegistered, branchExists, info.BranchName); err != nil {
			return err
		}
	}

	err = worktree.Create(worktreePath, info.BranchName, startPoint)
	if err != nil {
		if worktree.Exists(worktreePath) {
			os.RemoveAll(worktreePath)
		}
		return err
	}

	printSuccess(absPath)

	return executePostCreation(actionFlag, cliArgs, absPath, info)
}

func buildConflictMessage(info *worktree.WorktreeInfo, absPath, worktreePath string, worktreeDirExists, worktreeGitRegistered, branchExists bool) string {
	var message strings.Builder

	fmt.Fprintf(&message, "Target: create worktree for '%s'\n\nThis will:\n", info.BranchName)

	currentBranch := ""
	if worktreeGitRegistered {
		currentBranch, _ = git.GetWorktreeBranch(worktreePath)
	}

	if worktreeDirExists && worktreeGitRegistered {
		if currentBranch != "" {
			fmt.Fprintf(&message, "- Remove worktree at %s (currently on branch '%s')\n", absPath, currentBranch)
		} else {
			fmt.Fprintf(&message, "- Remove worktree at %s\n", absPath)
		}
	} else if worktreeGitRegistered {
		fmt.Fprintf(&message, "- Remove stale worktree record at %s\n", absPath)
	} else if worktreeDirExists {
		fmt.Fprintf(&message, "- Remove directory at %s\n", absPath)
	}

	if branchExists {
		fmt.Fprintf(&message, "- Delete existing branch '%s'\n", info.BranchName)
	}

	fmt.Fprintf(&message, "- Create worktree and branch for '%s'\n", info.BranchName)

	if worktreeDirExists && git.IsGitRepository(worktreePath) {
		if git.HasUncommittedChanges(worktreePath) {
			message.WriteString(fmt.Sprintf("\n⚠️  WARNING: Worktree at %s has uncommitted changes that will be PERMANENTLY DELETED. Consider committing or stashing changes first.\n", absPath))
		}
	}

	if branchExists {
		worktrees, err := git.GetWorktreeInfo()
		if err == nil {
			for _, wt := range worktrees {
				if wt.Branch == info.BranchName {
					if git.HasUncommittedChanges(wt.Path) {
						message.WriteString(fmt.Sprintf("\n⚠️ WARNING: Branch '%s' has uncommitted changes that will be PERMANENTLY DELETED. Consider committing or stashing changes first.\n", info.BranchName))
					}
					break
				}
			}
		}
	}

	message.WriteString("\nOverwrite?")
	return message.String()
}

func performCleanup(worktreePath string, worktreeDirExists, worktreeGitRegistered, branchExists bool, branchName string) error {
	if worktreeDirExists && worktreeGitRegistered {
		if err := git.WorktreeRemove(worktreePath, true); err != nil {
			return fmt.Errorf("failed to remove worktree: %w", err)
		}
	} else if worktreeDirExists {
		if err := os.RemoveAll(worktreePath); err != nil {
			return fmt.Errorf("failed to remove directory: %w", err)
		}
	} else if worktreeGitRegistered {
		if err := git.WorktreePrune(); err != nil {
			return fmt.Errorf("failed to prune worktree: %w", err)
		}
	}

	if branchExists {
		Log.Infof("Deleting existing branch '%s'...\n", branchName)
		if err := git.BranchDelete(branchName, true); err != nil {
			return fmt.Errorf("failed to delete branch: %w", err)
		}
	}

	return nil
}

func executePostCreation(actionFlag, cliArgs, absPath string, info *worktree.WorktreeInfo) error {
	if actionFlag != "" {
		if err := action.Execute(context.Background(), &action.ExecuteOptions{
			ActionName:   actionFlag,
			WorktreePath: absPath,
			Info:         info,
			CLIArgs:      cliArgs,
			Logger:       Log,
			Stdin:        os.Stdin,
			Stdout:       os.Stdout,
			Stderr:       os.Stderr,
			Env:          os.Environ(),
		}); err != nil {
			Log.Warnf("\n⚠️  Action '%s' failed: %v\n", actionFlag, err)
		}
	} else if cliArgs != "" {
		Log.Outf(logger.Magenta, "\nRunning in worktree: %s\n", cliArgs)

		if err := execext.RunCommand(context.Background(), &execext.RunCommandOptions{
			Command: cliArgs,
			Dir:     absPath,
			Env:     os.Environ(),
			Stdin:   os.Stdin,
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
		}); err != nil {
			Log.Warnf("\n⚠️  Command '%s' failed: %v\n", cliArgs, err)
		}
	}

	return nil
}

// printSuccess prints the final success message.
func printSuccess(path string) {
	Log.Outf(logger.Green, "\nWorktree created successfully!\n")
	Log.Outf(logger.Default, "Location: %s\n", path)
	Log.Outf(logger.Default, "\nTo switch to the worktree:\n")
	Log.Outf(logger.Cyan, "  cd %s\n", path)
}

// SanitizeBranchName is moved from types.go.
func SanitizeBranchName(name string) string {
	invalidChars := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return invalidChars.ReplaceAllString(name, "_")
}

// DetermineWorktreeType determines the type of worktree based on the input
// Returns the worktree type and an error message if invalid.
func DetermineWorktreeType(input string) (worktree.WorktreeType, error) {
	u, err := url.Parse(input)
	if err != nil {
		return worktree.Local, nil
	}

	if u.Scheme == "" {
		return worktree.Local, nil
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return worktree.Local, nil
	}

	prPattern := regexp.MustCompile(`^/[^/]+/[^/]+/pull/\d+(?:/.*)?$`)
	if prPattern.MatchString(u.Path) {
		return worktree.PR, nil
	}

	issuePattern := regexp.MustCompile(`^/[^/]+/[^/]+/issues/\d+(?:/.*)?$`)
	if issuePattern.MatchString(u.Path) {
		return worktree.Issue, nil
	}

	return worktree.Local, nil
}

var (
	prFlag     string
	issueFlag  string
	nameFlag   string
	actionFlag string
)

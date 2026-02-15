## gh wt

Create and manage git worktrees

### Synopsis

gh wt is a GitHub CLI extension that helps you create git worktrees. A GitHub pull request or issue url can also be used.

Examples:
  # Create worktree from PR URL
  gh wt add https://github.com/owner/repo/pull/123 -a claude -- "/review"

  # Create worktree from Issue URL
  gh wt add https://github.com/owner/repo/issues/456 -a claude -- "implement issue #456"

  # Create a worktree
  gh wt add my-feature-branch

  # Remove a worktree
  gh wt rm pr_123

### Options

```
  -f, --force      force operation without prompts
  -h, --help       help for gh wt
      --no-color   disable color output
      --verbose    verbose output
```

### SEE ALSO

* [gh wt action](gh_wt_action.md)	 - Manage and list actions
* [gh wt add](gh_wt_add.md)	 - Add a new worktree
* [gh wt completion](gh_wt_completion.md)	 - Generate shell completion scripts for gh wt commands
* [gh wt rm](gh_wt_rm.md)	 - Remove a worktree and its associated branch
* [gh wt run](gh_wt_run.md)	 - Run an action or command in an existing worktree


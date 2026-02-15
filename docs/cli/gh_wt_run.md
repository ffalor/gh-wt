## gh wt run

Run an action or command in an existing worktree

### Synopsis

Run an action or command in an existing worktree.

Use this command to:
- Run configured actions on worktrees that were created without an action
- Run commands directly in a worktree

Examples:
  # Run named action on worktree
  gh wt run pr_123 claude -- fix issue #456

  # Run command directly in worktree
  gh wt run pr_123 -- ls

  # Show help
  gh wt run pr_123

```
gh wt run <worktree> [action] [-- command] [flags]
```

### Options

```
  -h, --help   help for run
```

### Options inherited from parent commands

```
  -f, --force      force operation without prompts
      --no-color   disable color output
      --verbose    verbose output
```

### SEE ALSO

* [gh wt](gh_wt.md)	 - Create and manage git worktrees


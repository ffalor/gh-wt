## gh wt add

Add a new worktree

### Synopsis

Add a new git worktree from either:
  - A GitHub pull request URL or number
  - A GitHub issue URL or number
  - A name to use for the new worktree and branch


```
gh wt add [url|name] [flags]
```

### Options

```
  -a, --action string   action to run after worktree creation
  -h, --help            help for add
      --issue string    issue number, issue URL, or git remote URL with issue ref
      --pr string       PR number, PR URL, or git remote URL with PR ref
```

### Options inherited from parent commands

```
  -f, --force      force operation without prompts
      --no-color   disable color output
      --verbose    verbose output
```

### SEE ALSO

* [gh wt](gh_wt.md)	 - Create and manage git worktrees


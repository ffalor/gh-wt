## gh wt completion

Generate shell completion scripts for gh wt commands

### Synopsis

Generate shell completion scripts to enable tab completion for gh wt commands.

Tab completion provides:
- Command name completion (add, remove, run, action)
- Subcommand completion (install, uninstall)
- Flag completion

Supported shells: bash, zsh, fish, powershell

Examples:
  # Generate completion script for bash
  gh wt completion bash

  # Generate completion script for zsh
  gh wt completion zsh

  # Generate completion script for fish
  gh wt completion fish

  # Generate completion script for PowerShell
  gh wt completion powershell

  # Install completions automatically (detects your shell)
  gh wt completion install

  # Uninstall completions
  gh wt completion uninstall

```
gh wt completion [shell] [flags]
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
  -f, --force      force operation without prompts
      --no-color   disable color output
      --verbose    verbose output
```

### SEE ALSO

* [gh wt](gh_wt.md)	 - Create and manage git worktrees
* [gh wt completion install](gh_wt_completion_install.md)	 - Install shell completion for the detected shell
* [gh wt completion uninstall](gh_wt_completion_uninstall.md)	 - Uninstall shell completion for the detected shell


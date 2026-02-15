## gh wt completion install

Install shell completion for the detected shell

### Synopsis

Automatically install shell completion for your current shell.

This command detects your shell (bash, zsh, fish, or PowerShell) and installs
the completion script to the appropriate location. After installation, restart
your shell or source your shell configuration file.

Supported shells:
  - Bash: Installs to ~/.bash_completion.d/ or /etc/bash_completion.d/
  - Zsh: Installs to ~/.zsh/completions/
  - Fish: Installs to ~/.config/fish/completions/
  - PowerShell: Provides instructions to add to profile

Examples:
  gh wt completion install
	gh wt completion install --verbose

```
gh wt completion install [flags]
```

### Options

```
  -h, --help   help for install
```

### Options inherited from parent commands

```
  -f, --force      force operation without prompts
      --no-color   disable color output
      --verbose    verbose output
```

### SEE ALSO

* [gh wt completion](gh_wt_completion.md)	 - Generate shell completion scripts for gh wt commands


## gh wt completion uninstall

Uninstall shell completion for the detected shell

### Synopsis

Automatically uninstall shell completion for your current shell.

This command detects your shell and removes the completion script from the
appropriate location. After uninstallation, restart your shell or source
your shell configuration file.

Examples:
  gh wt completion uninstall
  gh wt completion uninstall --verbose

```
gh wt completion uninstall [flags]
```

### Options

```
  -h, --help   help for uninstall
```

### Options inherited from parent commands

```
  -f, --force      force operation without prompts
      --no-color   disable color output
      --verbose    verbose output
```

### SEE ALSO

* [gh wt completion](gh_wt_completion.md)	 - Generate shell completion scripts for gh wt commands


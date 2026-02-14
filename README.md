# gh-worktree

`gh-worktree` is a GitHub CLI extension for creating and cleaning up Git worktrees from pull requests, issues, or your current HEAD. It enhances your development workflow by enabling customizable post-create actions to automatically set up your environment, like launching tmux, interacting with AI tools, or running project bootstrap commands.

## Requirements

- `git`
- GitHub CLI (`gh`) authenticated for the target repo

## Install

Install from GitHub:

```bash
gh extension install ffalor/gh-worktree
```

## Quick Start

Current `gh worktree --help` output:

```text
gh-worktree is a GitHub CLI extension that helps you create git worktrees
from GitHub pull requests, issues, or local branch names.

Examples:
  # Create worktree from PR URL
  gh worktree https://github.com/owner/repo/pull/123 -action claude -- "/review"

  # Create worktree from Issue URL
  gh worktree https://github.com/owner/repo/issues/456 -action claude -- "implement issue #456"

  # Create a worktree
  gh worktree my-feature-branch

  # Remove a worktree
  gh worktree remove pr_123

Usage:
  gh-worktree [url|name] [flags]
  gh-worktree [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create a new worktree from a GitHub URL or branch name
  help        Help about any command
  remove      Remove a worktree and its associated branch

Flags:
  -f, --force      force operation without prompts
  -h, --help       help for gh-worktree
      --no-color   disable color output
  -v, --verbose    verbose output

Use "gh-worktree [command] --help" for more information about a command.
```

## Configuration

Config file path:
- `~/.config/gh-worktree/config.yaml`

Environment variables:
- Prefix: `GH_WT_`
- Example: `GH_WT_WORKTREE_DIR=~/github/worktree`

Minimal config:

```yaml
worktree_dir: "~/github/worktree"
```

### Actions

Actions are named command lists you can run with `--action <name>` after a worktree is created.

Example `config.yaml`:

```yaml
worktree_dir: "~/github/worktree"

actions:
  - name: tmux
    cmds:
      - tmux new-session -d -s "{{.BranchName}}" -c "{{.WorktreePath}}"
      - tmux send-keys -t "{{.BranchName}}":0 "git status" C-m
      - tmux split-window -h -t "{{.BranchName}}":0 -c "{{.WorktreePath}}"
      - tmux send-keys -t "{{.BranchName}}":0.1 "nvim ." C-m
      - tmux attach -t "{{.BranchName}}"

  - name: claude
    dir: "{{.WorktreePath}}"
    cmds:
      - claude -p "{{.CLI_ARGS}}"

  - name: dev_bootstrap
    cmds:
      - |
        if [ -f package.json ]; then
          if [ -f pnpm-lock.yaml ] && command -v pnpm >/dev/null 2>&1; then
            pnpm install
          elif [ -f yarn.lock ] && command -v yarn >/dev/null 2>&1; then
            yarn install
          elif [ -f package-lock.json ] && command -v npm >/dev/null 2>&1; then
            npm ci
          elif command -v npm >/dev/null 2>&1; then
            npm install
          fi
        fi
```

Run an action:

```bash
gh worktree create 123 --action tmux
```

Pass extra args to actions after `--`:

```bash
gh worktree create 123 --action claude -- "fix issue #456"
```

## Action Template Variables

Available in action `cmds` and optional `dir`:

- `{{.WorktreePath}}`
- `{{.BranchName}}`
- `{{.Action}}`
- `{{.CLI_ARGS}}`
- `{{.OS}}`
- `{{.ARCH}}`
- `{{.ROOT_DIR}}`
- `{{.Type}}`
- `{{.Owner}}`
- `{{.Repo}}`
- `{{.Number}}`
- `{{.WorktreeName}}`

## Behavior Notes

- On create conflicts (existing worktree/branch/path), the CLI prompts before destructive cleanup.
- `--force` skips these prompts.
- If an action fails, worktree creation still succeeds and the action failure is shown as a warning.

## Development

Build/install from source:

```bash
git clone https://github.com/ffalor/gh-worktree.git
cd gh-worktree
task install
```

Build:

```bash
task build
```

Run local binary:

```bash
task dev -- create 123
```

Test:

```bash
go test -v ./...
```

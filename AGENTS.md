# AGENTS.md - Agentic Coding Guidelines

`gh wt` is a GitHub CLI extension for managing git worktrees from GitHub PRs, Issues, and local branches. Built with Go 1.25.6 using cobra, viper, and go-gh.

## Build, Lint, and Test Commands

### Building the Project

```bash
task build
```

### Running the Extension

```bash
task install   # Install locally as gh extension
task run [args]    # Run the installed extension
task dev [args]    # Build and run for development
```

### Testing

```bash
go test -v -run TestFunctionName ./path/to/package  # Run a single test
go test -v ./...           # Run all tests with verbose output
go test -cover ./...       # Run tests with coverage
```

### Linting and Formatting

```bash
go fmt ./...          # Format code (required before committing)
go vet ./...          # Run go vet
go vet -shadow ./...  # Check for static analysis issues
```

### Development Tasks

```bash
task clean    # Clean built binary
task remove   # Remove installed extension
```

## Code Style Guidelines

### Imports

Group imports in order with blank lines between groups: standard library, external packages, internal packages.

```go
import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/ffalor/gh-wt/internal/config"
    "github.com/ffalor/gh-wt/internal/logger"
)
```

### Formatting

- Use `go fmt` for code formatting (enforced)
- Maximum line length: ~100 characters (soft limit)
- Use tabs for indentation

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Packages | lowercase | `git`, `config` |
| Functions | PascalCase | `CreateWorktree` |
| Variables | camelCase | `worktreePath` |
| Constants | PascalCase | `DefaultWorktreeBase` |
| Interfaces | PascalCase + -er | `Reader` |
| Errors | PascalCase + Err | `ErrNotFound` |

### Error Handling

- Use `fmt.Errorf("context: %w", err)` for wrapped errors
- Define sentinel errors: `var ErrCancelled = errors.New("cancelled")`
- Check with `errors.Is(err, ErrNotFound)`
- Avoid bare `panic()` except for unrecoverable conditions

### Command Structure (Cobra)

- Commands go in `cmd/` package
- Use `RunE` for commands that can fail
- Group flags in `init()` functions

```go
var createCmd = &cobra.Command{
    Use:   "create [url|name]",
    Short: "Create a new worktree",
    RunE:  runCreate,
}

func init() {
    rootCmd.AddCommand(createCmd)
    createCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "force operation")
}
```

### Configuration

- Use Viper for configuration (`internal/config`)
- Config file: `~/.config/gh-wt/config.yaml`
- Support config file, environment variables (prefix: `GH_WT_`), and flags
- Provide sensible defaults
- Use `config.Get()` to retrieve typed configuration

### GitHub CLI Integration

- Use `github.com/cli/go-gh/v2` for API calls
- Parse JSON with `json.Unmarshal`
- Use `gh.Exec(args...)` for running gh CLI commands

### Logging and Output

- Use the `logger` package (`internal/logger`) for user-facing output
- Use `Log.Outf()` with colors: `logger.Default`, `logger.Green`, `logger.Red`, `logger.Yellow`, `logger.Cyan`, `logger.Blue`, `logger.Magenta`
- Use `Log.Errorf()` for errors, `Log.Warnf()` for warnings, `Log.Infof()` for info

### Shell Execution

- Use `mvdan.cc/sh/v3` for shell command execution in actions
- Access via `internal/execext` package

## Commit Messages

Use Conventional Commits format:

```
<type>(<scope>): <description>

[optional body]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Example: `feat(worktree): add support for GitHub issues`

## Project Structure

```
.
├── main.go              # Entry point
├── cmd/                 # CLI commands (cobra)
│   ├── root.go          # Root command, flags: --force, --verbose, --no-color
│   ├── create.go        # Create worktree from PR, Issue, or local branch
│   └── remove.go        # Remove worktree and associated branch
├── internal/
│   ├── action/          # Post-creation action execution with templating
│   ├── config/         # Viper configuration management
│   ├── execext/        # Shell command execution (mvdan/sh)
│   ├── git/            # Git operations (branch, worktree)
│   ├── logger/         # Colored logging output
│   └── worktree/       # Worktree creation/removal logic
├── .agents/skills/     # Agent skills
├── Taskfile.yml        # Development tasks
├── go.mod              # Go module definition
└── config.example.yaml # Example configuration
```

## Commands

### Add (`gh wt add`)

- Create worktree from PR URL/number: `gh wt https://github.com/owner/repo/pull/123`
- Create worktree from Issue URL/number: `gh wt https://github.com/owner/repo/issues/456`
- Create local worktree: `gh wt my-feature-branch`
- Flags: `--pr`, `--issue`, `-a, --action`, `--use-existing`

### Remove (`gh wt rm`)

- Remove worktree by name: `gh wt rm <worktree-name>`
- Flags: `--force` to skip confirmation

### Root

- Can also invoke add directly: `gh wt <url|name>`

## Configuration

Example `~/.config/gh-wt/config.yaml`:

```yaml
worktree_dir: "~/github/worktree"

actions:
  - name: tmux
    cmds:
      - tmux new-session -d -s {{.BranchName}}
      - tmux send-keys -t {{.BranchName}} "cd {{.WorktreePath}} && vim ." C-m
```

Action template variables:

- `{{.WorktreePath}}` - Path to worktree
- `{{.WorktreeName}}` - Name of the worktree directory
- `{{.BranchName}}` - Branch name
- `{{.Action}}` - Action name
- `{{.CLI_ARGS}}` - CLI arguments after --
- `{{.OS}}` - Operating system
- `{{.ARCH}}` - Architecture
- `{{.ROOT_DIR}}` - Git root directory

## Additional Notes

- This is a GitHub CLI extension - install with `gh extension install .`
- Requires `gh` CLI to be installed
- Worktrees stored in `~/github/worktree` by default (configurable)
- Uses bare repositories for remote PR/Issue worktrees
- Supports post-creation actions with templating

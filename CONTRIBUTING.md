# Contributing to gh-wt

Thank you for your interest in contributing to gh-wt!

## Requirements

- [Go](https://go.dev/) 1.25.6+
- [GitHub CLI](https://cli.github.com/) (`gh`)
- [Task](https://taskfile.dev/) (optional, for development tasks)

## Development Setup

```bash
git clone https://github.com/ffalor/gh-worktree.git
cd gh-worktree
task install
```

## Making Changes

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. Run formatting and tests:
   ```bash
   go fmt ./...
   go vet ./...
   go test -v ./...
   ```
5. Commit using [Conventional Commits](https://www.conventionalcommits.org/)
6. Push and open a Pull Request

## Project Structure

- `cmd/` - CLI commands (Cobra)
- `internal/` - Core packages
  - `action/` - Post-creation action execution
  - `config/` - Configuration management
  - `execext/` - Shell command execution
  - `git/` - Git operations
  - `logger/` - Logging output
  - `worktree/` - Worktree management

## Code Style

- Follow the conventions in `AGENTS.md`
- Use `go fmt` for formatting
- Group imports: standard library, external packages, internal packages
- Wrap errors with context: `fmt.Errorf("context: %w", err)`

## Testing

Run all tests:
```bash
go test -v ./...
```

Run a specific test:
```bash
go test -v -run TestFunctionName ./path/to/package
```

## Commands

Development tasks are defined in `Taskfile.yml`:

```bash
task build    # Build the binary
task install # Install as gh extension
task dev     # Build and run for development
task clean   # Clean built binary
```

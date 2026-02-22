package action

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/ffalor/gh-wt/internal/config"
	"github.com/ffalor/gh-wt/internal/execext"
	"github.com/ffalor/gh-wt/internal/git"
	"github.com/ffalor/gh-wt/internal/logger"
	"github.com/ffalor/gh-wt/internal/worktree"
)

var (
	// ErrNilOptions is returned when Execute receives nil options.
	ErrNilOptions = errors.New("action: nil options given")
	// ErrNilLogger is returned when ExecuteOptions.Logger is nil.
	ErrNilLogger = errors.New("action: nil logger given")
)

// ExecuteOptions contains dependencies and context for running an action.
type ExecuteOptions struct {
	ActionName   string
	WorktreePath string
	Info         *worktree.WorktreeInfo
	CLIArgs      string
	Logger       *logger.Logger
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
	Env          []string
}

// Execute runs the specified action after templating its commands.
func Execute(ctx context.Context, opts *ExecuteOptions) error {
	if opts == nil {
		return ErrNilOptions
	}
	if opts.Logger == nil {
		return ErrNilLogger
	}
	if strings.TrimSpace(opts.ActionName) == "" {
		return fmt.Errorf("action: action name is required")
	}
	if strings.TrimSpace(opts.WorktreePath) == "" {
		return fmt.Errorf("action: worktree path is required")
	}
	if opts.Info == nil {
		return fmt.Errorf("action: worktree info is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	stdin := opts.Stdin
	if stdin == nil {
		stdin = os.Stdin
	}

	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	env := opts.Env
	if len(env) == 0 {
		env = os.Environ()
	}

	cfg, err := config.Get()
	if err != nil {
		return err
	}

	var action *config.Action
	for i := range cfg.Actions {
		if cfg.Actions[i].Name == opts.ActionName {
			action = &cfg.Actions[i]
			break
		}
	}

	if action == nil {
		return fmt.Errorf("action '%s' not found in config", opts.ActionName)
	}

	// Get git root directory
	rootDir, err := git.GetGitRoot()
	if err != nil {
		return fmt.Errorf("failed to get git root directory: %w", err)
	}

	// Prepare data for template
	data := struct {
		WorktreePath string
		WorktreeName string
		Action       string
		CLI_ARGS     string
		OS           string
		ARCH         string
		ROOT_DIR     string
		*worktree.WorktreeInfo
	}{
		WorktreePath: opts.WorktreePath,
		WorktreeName: filepath.Base(opts.WorktreePath),
		Action:       opts.ActionName,
		CLI_ARGS:     opts.CLIArgs,
		OS:           runtime.GOOS,
		ARCH:         runtime.GOARCH,
		ROOT_DIR:     rootDir,
		WorktreeInfo: opts.Info,
	}

	runDir := opts.WorktreePath

	if action.Dir != "" {
		tmpl, err := template.New("dir").Parse(action.Dir)
		if err != nil {
			return fmt.Errorf("failed to parse action directory template: %w", err)
		}
		var renderedDir bytes.Buffer
		if err := tmpl.Execute(&renderedDir, data); err != nil {
			return fmt.Errorf("failed to render action directory template: %w", err)
		}
		runDir = renderedDir.String()
	}

	opts.Logger.Outf(logger.Magenta, "\nRunning action '%s' in %s...\n", opts.ActionName, runDir)

	for _, cmdStr := range action.Cmds {
		tmpl, err := template.New("cmd").Parse(cmdStr)
		if err != nil {
			return fmt.Errorf("failed to parse command template: %w", err)
		}

		var renderedCmd bytes.Buffer
		if err := tmpl.Execute(&renderedCmd, data); err != nil {
			return fmt.Errorf("failed to render command template: %w", err)
		}

		finalCmd := renderedCmd.String()
		opts.Logger.Outf(logger.Magenta, "[%s]: %s\n", opts.ActionName, finalCmd)

		if err := execext.RunCommand(ctx, &execext.RunCommandOptions{
			Command: finalCmd,
			Dir:     runDir,
			Env:     env,
			Stdin:   stdin,
			Stdout:  stdout,
			Stderr:  stderr,
		}); err != nil {
			return fmt.Errorf("command '%s' failed: %w", finalCmd, err)
		}
	}

	opts.Logger.Outf(logger.Green, "Action finished successfully.\n")
	return nil
}

package action

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"

	"github.com/ffalor/gh-worktree/internal/config"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// WorktreeType represents the type of worktree
type WorktreeType string

const (
	Issue WorktreeType = "issue"
	PR    WorktreeType = "pr"
	Local WorktreeType = "local"
)

// WorktreeInfo is a copy of the struct from cmd/create.go,
// moved here for shared access.
type WorktreeInfo struct {
	Type         WorktreeType // Changed from string to WorktreeType
	Owner        string
	Repo         string
	Number       int
	BranchName   string
	WorktreeName string
}

// Execute runs the specified action after templating its commands.
func Execute(actionName, worktreePath string, info *WorktreeInfo, cliArgs string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	var action *config.Action
	for i := range cfg.Actions {
		if cfg.Actions[i].Name == actionName {
			action = &cfg.Actions[i]
			break
		}
	}

	if action == nil {
		return fmt.Errorf("action '%s' not found in config", actionName)
	}

	// Get git root directory
	gitRoot, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return fmt.Errorf("failed to get git root directory: %w", err)
	}
	rootDir := strings.TrimSpace(string(gitRoot))

	// Prepare data for template
	data := struct {
		WorktreePath string
		Action       string
		CLI_ARGS     string
		OS           string
		ARCH         string
		ROOT_DIR     string
		*WorktreeInfo
	}{
		WorktreePath: worktreePath,
		Action:       actionName,
		CLI_ARGS:     cliArgs,
		OS:           runtime.GOOS,
		ARCH:         runtime.GOARCH,
		ROOT_DIR:     rootDir,
		WorktreeInfo: info,
	}

	runDir := worktreePath // Default to worktreePath

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

	fmt.Printf("\nRunning action '%s' in %s...\n", actionName, runDir)
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
		fmt.Printf("$ %s\n", finalCmd)

		parser := syntax.NewParser()
		prog, err := parser.Parse(strings.NewReader(finalCmd), "")
		if err != nil {
			return fmt.Errorf("failed to parse shell command: %w", err)
		}

		runner, err := interp.New(
			interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
			interp.Dir(runDir),
			interp.Env(expand.ListEnviron(os.Environ()...)),
		)
		if err != nil {
			return fmt.Errorf("failed to create shell interpreter: %w", err)
		}

		if err := runner.Run(context.Background(), prog); err != nil {
			return fmt.Errorf("command '%s' failed: %w", finalCmd, err)
		}
	}

	fmt.Println("Action finished successfully.")
	return nil
}

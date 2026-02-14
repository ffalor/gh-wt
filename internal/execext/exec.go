package execext

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// ErrNilOptions is returned when nil options are provided.
var ErrNilOptions = errors.New("execext: nil options given")

// RunCommandOptions configures shell command execution.
type RunCommandOptions struct {
	Command   string
	Dir       string
	Env       []string
	PosixOpts []string
	BashOpts  []string
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
}

// RunCommand runs a shell command with mvdan/sh.
func RunCommand(ctx context.Context, opts *RunCommandOptions) error {
	if opts == nil {
		return ErrNilOptions
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

	environ := opts.Env
	if len(environ) == 0 {
		environ = os.Environ()
	}

	posixOpts := append([]string{}, opts.PosixOpts...)
	posixOpts = append(posixOpts, "e")

	var params []string
	for _, opt := range posixOpts {
		if len(opt) == 1 {
			params = append(params, fmt.Sprintf("-%s", opt))
			continue
		}
		params = append(params, "-o", opt)
	}

	runner, err := interp.New(
		interp.Params(params...),
		interp.Env(expand.ListEnviron(environ...)),
		interp.StdIO(stdin, stdout, stderr),
		interp.Dir(opts.Dir),
	)
	if err != nil {
		return err
	}

	parser := syntax.NewParser()

	if len(opts.BashOpts) > 0 {
		shoptCmd := fmt.Sprintf("shopt -s %s", strings.Join(opts.BashOpts, " "))
		prog, err := parser.Parse(strings.NewReader(shoptCmd), "")
		if err != nil {
			return err
		}
		if err := runner.Run(ctx, prog); err != nil {
			return err
		}
	}

	prog, err := parser.Parse(strings.NewReader(opts.Command), "")
	if err != nil {
		return err
	}
	return runner.Run(ctx, prog)
}

package logger

import (
	"fmt"
	"io"
	"os"
)

type (
	Color     func() PrintFunc
	PrintFunc func(io.Writer, string, ...any)
)

func None() PrintFunc {
	return func(w io.Writer, s string, args ...any) {
		fmt.Fprintf(w, s, args...)
	}
}

func Default() PrintFunc {
	return colorPrint("0")
}

func Blue() PrintFunc {
	return colorPrint("34")
}

func Green() PrintFunc {
	return colorPrint("32")
}

func Cyan() PrintFunc {
	return colorPrint("36")
}

func Yellow() PrintFunc {
	return colorPrint("33")
}

func Magenta() PrintFunc {
	return colorPrint("35")
}

func Red() PrintFunc {
	return colorPrint("31")
}

func colorPrint(code string) PrintFunc {
	return func(w io.Writer, s string, args ...any) {
		fmt.Fprintf(w, "\x1b[%sm", code)
		fmt.Fprintf(w, s, args...)
		fmt.Fprint(w, "\x1b[0m")
	}
}

// Logger is a wrapper that prints stuff to STDOUT or STDERR,
// with optional color and verbosity.
type Logger struct {
	Stdout  io.Writer
	Stderr  io.Writer
	Verbose bool
	Color   bool
}

// NewLogger creates a new Logger instance.
func NewLogger(verbose, useColor bool) *Logger {
	return &Logger{
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Verbose: verbose,
		Color:   useColor,
	}
}

// Outf prints stuff to STDOUT.
func (l *Logger) Outf(c Color, s string, args ...any) {
	l.FOutf(l.Stdout, c, s, args...)
}

// FOutf prints stuff to the given writer.
func (l *Logger) FOutf(w io.Writer, c Color, s string, args ...any) {
	if len(args) == 0 {
		s, args = "%s", []any{s}
	}
	if !l.Color {
		c = None
	}
	print := c()
	print(w, s, args...)
}

// VerboseOutf prints stuff to STDOUT if verbose mode is enabled.
func (l *Logger) VerboseOutf(c Color, s string, args ...any) {
	if l.Verbose {
		l.Outf(c, s, args...)
	}
}

// Errf prints stuff to STDERR.
func (l *Logger) Errf(c Color, s string, args ...any) {
	if len(args) == 0 {
		s, args = "%s", []any{s}
	}
	if !l.Color {
		c = None
	}
	print := c()
	print(l.Stderr, s, args...)
}

// VerboseErrf prints stuff to STDERR if verbose mode is enabled.
func (l *Logger) VerboseErrf(c Color, s string, args ...any) {
	if l.Verbose {
		l.Errf(c, s, args...)
	}
}

// Warnf prints a warning message to STDERR.
func (l *Logger) Warnf(s string, args ...any) {
	l.Errf(Yellow, s, args...)
}

// Errorf prints an error message to STDERR.
func (l *Logger) Errorf(s string, args ...any) {
	l.Errf(Red, s, args...)
}

// Infof prints an informational message to STDOUT.
func (l *Logger) Infof(s string, args ...any) {
	l.Outf(Cyan, s, args...)
}

// Plainf prints a plain message to STDOUT.
func (l *Logger) Plainf(s string, args ...any) {
	if len(args) == 0 {
		s, args = "%s", []any{s}
	}
	fmt.Fprintf(l.Stdout, s, args...)
}

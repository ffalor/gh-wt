package completion

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ffalor/gh-wt/internal/logger"
	"github.com/spf13/cobra"
)

// ShellType represents the detected shell type
type ShellType string

const (
	ShellBash       ShellType = "bash"
	ShellZsh        ShellType = "zsh"
	ShellFish       ShellType = "fish"
	ShellPowerShell ShellType = "powershell"
	ShellUnknown    ShellType = "unknown"
)

// DetectShell detects the current shell from environment variables
func DetectShell(Log *logger.Logger) ShellType {
	Log.Plainf("Detecting current shell\n")

	// Check shell-specific version variables first (most reliable)
	if os.Getenv("ZSH_VERSION") != "" {
		Log.Plainf("Detected zsh from ZSH_VERSION\n")
		return ShellZsh
	}
	if os.Getenv("BASH_VERSION") != "" {
		Log.Plainf("Detected bash from BASH_VERSION\n")
		return ShellBash
	}
	if os.Getenv("FISH_VERSION") != "" {
		Log.Plainf("Detected fish from FISH_VERSION\n")
		return ShellFish
	}

	// Fall back to $SHELL environment variable
	shell := os.Getenv("SHELL")
	if shell == "" {
		Log.Plainf("SHELL environment variable not set, checking platform\n")
		// On Windows, check for PowerShell
		if runtime.GOOS == "windows" {
			Log.Plainf("Detected Windows, assuming PowerShell\n")
			return ShellPowerShell
		}
		Log.Plainf("Could not detect shell\n")
		return ShellUnknown
	}

	Log.Plainf("SHELL environment variable: %s\n", shell)

	// Extract shell name from path
	shellName := filepath.Base(shell)
	Log.Plainf("Shell base name: %s\n", shellName)

	switch {
	case strings.Contains(shellName, "bash"):
		Log.Plainf("Detected bash from SHELL\n")
		return ShellBash
	case strings.Contains(shellName, "zsh"):
		Log.Plainf("Detected zsh from SHELL\n")
		return ShellZsh
	case strings.Contains(shellName, "fish"):
		Log.Plainf("Detected fish from SHELL\n")
		return ShellFish
	case strings.Contains(shellName, "pwsh") || strings.Contains(shellName, "powershell"):
		Log.Plainf("Detected PowerShell from SHELL\n")
		return ShellPowerShell
	default:
		Log.Plainf("Unknown shell: %s\n", shellName)
		return ShellUnknown
	}
}

// InstallShellCompletion installs shell completion for the detected shell
func InstallShellCompletion(Log *logger.Logger, cmd *cobra.Command) error {
	Log.VerboseOutf(logger.Default, "Starting shell completion installation\n")

	shellType := DetectShell(Log)
	Log.VerboseOutf(logger.Default, "Detected shell type: %s\n", shellType)

	if shellType == ShellUnknown {
		return fmt.Errorf("could not detect shell type. Please install completions manually using: gh wt completion <shell>")
	}

	Log.Outf(logger.Default, "Detected shell: %s\n", shellType)

	switch shellType {
	case ShellBash:
		return installBashCompletion(Log, cmd)
	case ShellZsh:
		return installZshCompletion(Log, cmd)
	case ShellFish:
		return installFishCompletion(Log, cmd)
	case ShellPowerShell:
		return installPowerShellCompletion(Log, cmd)
	default:
		return fmt.Errorf("shell completion not supported for: %s", shellType)
	}
}

// installBashCompletion installs bash completion
func installBashCompletion(Log *logger.Logger, cmd *cobra.Command) error {
	Log.VerboseOutf(logger.Default, "[verbose] Installing bash completion\n")

	// Generate completion script using Cobra
	var buf strings.Builder
	if err := cmd.GenBashCompletion(&buf); err != nil {
		return fmt.Errorf("failed to generate bash completion: %w", err)
	}

	completionScript := buf.String()

	// Determine installation path
	var completionPath string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Try to determine the best location for bash completions
	if runtime.GOOS == "darwin" {
		// macOS with Homebrew
		brewPrefix := os.Getenv("HOMEBREW_PREFIX")
		if brewPrefix == "" {
			// Try common locations
			for _, prefix := range []string{"/opt/homebrew", "/usr/local"} {
				if _, err := os.Stat(filepath.Join(prefix, "etc", "bash_completion.d")); err == nil {
					brewPrefix = prefix
					break
				}
			}
		}
		if brewPrefix != "" {
			completionPath = filepath.Join(brewPrefix, "etc", "bash_completion.d", "gh-wt")
		} else {
			completionPath = filepath.Join(homeDir, ".bash_completion.d", "gh-wt")
		}
	} else {
		// Linux
		if _, err := os.Stat("/etc/bash_completion.d"); err == nil {
			completionPath = "/etc/bash_completion.d/gh-wt"
		} else {
			completionPath = filepath.Join(homeDir, ".bash_completion.d", "gh-wt")
		}
	}

	// Create directory if needed (for user-level installations)
	completionDir := filepath.Dir(completionPath)
	if strings.HasPrefix(completionDir, homeDir) {
		// Use restrictive permissions (0750) following principle of least privilege
		if err := os.MkdirAll(completionDir, 0750); err != nil {
			return fmt.Errorf("failed to create completion directory: %w", err)
		}
	}

	// Try to write completion file
	// Use restrictive permissions (0600) following principle of least privilege
	err = os.WriteFile(completionPath, []byte(completionScript), 0600)
	if err != nil && strings.HasPrefix(completionPath, "/etc") {
		// If system-wide installation fails, fall back to user directory
		Log.VerboseOutf(logger.Default, "Failed to install system-wide, falling back to user directory: %v\n", err)
		completionPath = filepath.Join(homeDir, ".bash_completion.d", "gh-wt")
		// Use restrictive permissions (0750) following principle of least privilege
		if err := os.MkdirAll(filepath.Dir(completionPath), 0750); err != nil {
			return fmt.Errorf("failed to create user completion directory: %w", err)
		}
		// Use restrictive permissions (0600) following principle of least privilege
		if err := os.WriteFile(completionPath, []byte(completionScript), 0600); err != nil {
			return fmt.Errorf("failed to write completion file: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to write completion file: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Installed bash completion to: %s\n", completionPath)

	// Check if .bashrc sources completions
	bashrcPath := filepath.Join(homeDir, ".bashrc")
	if strings.HasPrefix(completionPath, homeDir) {
		// For user-level installations, check if .bashrc sources the completion directory
		// Clean and validate the path to prevent path traversal
		cleanBashrcPath := filepath.Clean(bashrcPath)
		if !filepath.IsAbs(cleanBashrcPath) {
			Log.VerboseOutf(logger.Default, "Invalid bashrc path (not absolute): %s\n", bashrcPath)
			return fmt.Errorf("invalid bashrc path: %s", bashrcPath)
		}
		bashrcContent, err := os.ReadFile(cleanBashrcPath)
		needsSourceLine := true
		if err == nil {
			if strings.Contains(string(bashrcContent), ".bash_completion.d") ||
				strings.Contains(string(bashrcContent), completionPath) {
				needsSourceLine = false
			}
		}

		if needsSourceLine {
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "To enable completions, add the following to your ~/.bashrc:")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintf(os.Stderr, "  for f in ~/.bash_completion.d/*; do [ -f \"$f\" ] && source \"$f\"; done\n")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "Then restart your shell or run: source ~/.bashrc")
		} else {
			fmt.Fprintln(os.Stderr, "Please restart your shell for completions to take effect")
		}
	} else {
		fmt.Fprintln(os.Stderr, "Please restart your shell for completions to take effect")
	}

	return nil
}

// installZshCompletion installs zsh completion
func installZshCompletion(Log *logger.Logger, cmd *cobra.Command) error {
	Log.VerboseOutf(logger.Default, "Installing zsh completion\n")

	// Generate completion script using Cobra
	var buf strings.Builder
	if err := cmd.GenZshCompletion(&buf); err != nil {
		return fmt.Errorf("failed to generate zsh completion: %w", err)
	}

	completionScript := buf.String()

	// Determine installation path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check for fpath directories
	// Try user's local completion directory first
	userCompletionDir := filepath.Join(homeDir, ".zsh", "completions")
	// Use restrictive permissions (0750) following principle of least privilege
	if err := os.MkdirAll(userCompletionDir, 0750); err != nil {
		return fmt.Errorf("failed to create completion directory: %w", err)
	}
	completionPath := filepath.Join(userCompletionDir, "_gh-wt")

	// Write completion file
	// Use restrictive permissions (0600) following principle of least privilege
	if err := os.WriteFile(completionPath, []byte(completionScript), 0600); err != nil {
		return fmt.Errorf("failed to write completion file: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Installed zsh completion to: %s\n", completionPath)

	// Check if .zshrc configures fpath
	zshrcPath := filepath.Join(homeDir, ".zshrc")
	// Clean and validate the path to prevent path traversal
	cleanZshrcPath := filepath.Clean(zshrcPath)
	if !filepath.IsAbs(cleanZshrcPath) {
		Log.VerboseOutf(logger.Default, "Invalid zshrc path (not absolute): %s\n", zshrcPath)
		return fmt.Errorf("invalid zshrc path: %s", zshrcPath)
	}
	zshrcContent, err := os.ReadFile(cleanZshrcPath)
	needsFpath := true
	if err == nil {
		if strings.Contains(string(zshrcContent), userCompletionDir) {
			needsFpath = false
		}
	}

	if needsFpath {
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "To enable completions, add the following to your ~/.zshrc:")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintf(os.Stderr, "  fpath=(~/.zsh/completions $fpath)\n")
		fmt.Fprintf(os.Stderr, "  autoload -Uz compinit && compinit\n")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Then restart your shell or run: source ~/.zshrc")
	} else {
		fmt.Fprintln(os.Stderr, "Please restart your shell for completions to take effect")
	}

	return nil
}

// installFishCompletion installs fish completion
func installFishCompletion(Log *logger.Logger, cmd *cobra.Command) error {
	Log.VerboseOutf(logger.Default, "Installing fish completion\n")

	// Generate completion script using Cobra
	var buf strings.Builder
	if err := cmd.GenFishCompletion(&buf, true); err != nil {
		return fmt.Errorf("failed to generate fish completion: %w", err)
	}

	completionScript := buf.String()

	// Determine installation path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Fish completion directory
	completionDir := filepath.Join(homeDir, ".config", "fish", "completions")
	// Use restrictive permissions (0750) following principle of least privilege
	if err := os.MkdirAll(completionDir, 0750); err != nil {
		return fmt.Errorf("failed to create completion directory: %w", err)
	}

	completionPath := filepath.Join(completionDir, "gh-wt.fish")

	// Write completion file
	// Use restrictive permissions (0600) following principle of least privilege
	if err := os.WriteFile(completionPath, []byte(completionScript), 0600); err != nil {
		return fmt.Errorf("failed to write completion file: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Installed fish completion to: %s\n", completionPath)
	fmt.Fprintln(os.Stderr, "Fish will automatically load completions on next shell start")

	return nil
}

// installPowerShellCompletion installs PowerShell completion
func installPowerShellCompletion(Log *logger.Logger, cmd *cobra.Command) error {
	Log.VerboseOutf(logger.Default, "Installing PowerShell completion\n")

	// Determine PowerShell profile path
	var profileCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		profileCmd = exec.Command("powershell", "-NoProfile", "-Command", "echo $PROFILE")
	} else {
		profileCmd = exec.Command("pwsh", "-NoProfile", "-Command", "echo $PROFILE")
	}

	var profileBuf strings.Builder
	profileCmd.Stdout = &profileBuf
	if err := profileCmd.Run(); err != nil {
		return fmt.Errorf("failed to get PowerShell profile path: %w", err)
	}

	profilePath := strings.TrimSpace(profileBuf.String())

	fmt.Fprintf(os.Stderr, "PowerShell profile path: %s\n", profilePath)
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "To enable completions, add the following to your PowerShell profile:")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  gh wt completion powershell | Out-String | Invoke-Expression")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Or run the following command to append it automatically:")
	fmt.Fprintln(os.Stderr, "")
	if runtime.GOOS == "windows" {
		fmt.Fprintln(os.Stderr, "  gh wt completion powershell >> $PROFILE")
	} else {
		fmt.Fprintln(os.Stderr, "echo 'gh wt completion powershell | Out-String | Invoke-Expression' >> $PROFILE")
	}
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Then restart your shell or run: . $PROFILE")

	return nil
}

// UninstallShellCompletion uninstalls shell completion for the detected shell
func UninstallShellCompletion(Log *logger.Logger) error {
	Log.VerboseOutf(logger.Default, "Starting shell completion uninstallation\n")

	shellType := DetectShell(Log)
	Log.VerboseOutf(logger.Default, "Detected shell type: %s\n", shellType)

	if shellType == ShellUnknown {
		return fmt.Errorf("could not detect shell type. Please uninstall completions manually")
	}

	Log.Outf(logger.Default, "Detected shell: %s\n", shellType)

	switch shellType {
	case ShellBash:
		return uninstallBashCompletion(Log)
	case ShellZsh:
		return uninstallZshCompletion(Log)
	case ShellFish:
		return uninstallFishCompletion(Log)
	case ShellPowerShell:
		return uninstallPowerShellCompletion(Log)
	default:
		return fmt.Errorf("shell completion not supported for: %s", shellType)
	}
}

// uninstallBashCompletion uninstalls bash completion
func uninstallBashCompletion(Log *logger.Logger) error {
	Log.VerboseOutf(logger.Default, "Uninstalling bash completion\n")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check all possible locations where completion might be installed
	var possiblePaths []string

	// User-level installations
	possiblePaths = append(possiblePaths, filepath.Join(homeDir, ".bash_completion.d", "gh-wt"))

	// macOS with Homebrew
	if runtime.GOOS == "darwin" {
		brewPrefix := os.Getenv("HOMEBREW_PREFIX")
		if brewPrefix == "" {
			for _, prefix := range []string{"/opt/homebrew", "/usr/local"} {
				if _, err := os.Stat(filepath.Join(prefix, "etc", "bash_completion.d")); err == nil {
					possiblePaths = append(possiblePaths, filepath.Join(prefix, "etc", "bash_completion.d", "gh-wt"))
				}
			}
		} else {
			possiblePaths = append(possiblePaths, filepath.Join(brewPrefix, "etc", "bash_completion.d", "gh-wt"))
		}
	}

	// System-wide installations (Linux)
	if runtime.GOOS == "linux" {
		possiblePaths = append(possiblePaths, "/etc/bash_completion.d/gh-wt")
	}

	removed := false
	var lastErr error

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			Log.VerboseOutf(logger.Default, "Found completion file at: %s\n", path)
			if err := os.Remove(path); err != nil {
				Log.VerboseOutf(logger.Default, "Failed to remove %s: %v\n", path, err)
				lastErr = err
				continue
			}
			fmt.Fprintf(os.Stderr, "Removed bash completion from: %s\n", path)
			removed = true
		}
	}

	if !removed {
		return fmt.Errorf("no bash completion file found to remove")
	}

	if lastErr != nil {
		fmt.Fprintln(os.Stderr, "Some completion files could not be removed (may require elevated permissions)")
	}

	fmt.Fprintln(os.Stderr, "Please restart your shell for changes to take effect")

	return nil
}

// uninstallZshCompletion uninstalls zsh completion
func uninstallZshCompletion(Log *logger.Logger) error {
	Log.VerboseOutf(logger.Default, "Uninstalling zsh completion\n")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check possible locations
	completionPath := filepath.Join(homeDir, ".zsh", "completions", "_gh-wt")

	if _, err := os.Stat(completionPath); err != nil {
		return fmt.Errorf("no zsh completion file found at: %s", completionPath)
	}

	Log.VerboseOutf(logger.Default, "Found completion file at: %s\n", completionPath)

	if err := os.Remove(completionPath); err != nil {
		return fmt.Errorf("failed to remove completion file: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Removed zsh completion from: %s\n", completionPath)
	fmt.Fprintln(os.Stderr, "Please restart your shell for changes to take effect")

	return nil
}

// uninstallFishCompletion uninstalls fish completion
func uninstallFishCompletion(Log *logger.Logger) error {
	Log.VerboseOutf(logger.Default, "Uninstalling fish completion\n")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	completionPath := filepath.Join(homeDir, ".config", "fish", "completions", "gh-wt.fish")

	if _, err := os.Stat(completionPath); err != nil {
		return fmt.Errorf("no fish completion file found at: %s", completionPath)
	}

	Log.VerboseOutf(logger.Default, "Found completion file at: %s\n", completionPath)

	if err := os.Remove(completionPath); err != nil {
		return fmt.Errorf("failed to remove completion file: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Removed fish completion from: %s\n", completionPath)
	fmt.Fprintln(os.Stderr, "Fish will automatically detect the removal on next shell start")

	return nil
}

// uninstallPowerShellCompletion uninstalls PowerShell completion
func uninstallPowerShellCompletion(Log *logger.Logger) error {
	Log.VerboseOutf(logger.Default, "Uninstalling PowerShell completion\n")

	// Determine PowerShell profile path
	var profileCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		profileCmd = exec.Command("powershell", "-NoProfile", "-Command", "echo $PROFILE")
	} else {
		profileCmd = exec.Command("pwsh", "-NoProfile", "-Command", "echo $PROFILE")
	}

	var profileBuf strings.Builder
	profileCmd.Stdout = &profileBuf
	if err := profileCmd.Run(); err != nil {
		return fmt.Errorf("failed to get PowerShell profile path: %w", err)
	}

	profilePath := strings.TrimSpace(profileBuf.String())

	fmt.Fprintf(os.Stderr, "PowerShell profile path: %s\n", profilePath)
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "To uninstall completions, remove the following line from your PowerShell profile:")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  gh wt completion powershell | Out-String | Invoke-Expression")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Then restart your shell or run: . $PROFILE")

	return nil
}

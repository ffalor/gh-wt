//go:build integration

package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCommandIntegration(t *testing.T) {
	setup := setupIntegrationTest(t)

	t.Run("list with no managed worktrees shows warning", func(t *testing.T) {
		cmd := exec.Command(setup.binaryPath, "list")
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		// Should succeed (exit 0) even with no worktrees
		require.NoError(t, err, "list should succeed: %s", outputStr)

		// Should show warning about no worktrees
		assert.Contains(t, outputStr, "No worktrees found")
	})

	t.Run("ls alias works", func(t *testing.T) {
		cmd := exec.Command(setup.binaryPath, "ls")
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		// Should succeed (exit 0)
		require.NoError(t, err, "ls alias should succeed: %s", outputStr)

		// Should show same warning
		assert.Contains(t, outputStr, "No worktrees found")
	})

	t.Run("list rejects arguments", func(t *testing.T) {
		cmd := exec.Command(setup.binaryPath, "list", "extra-arg")
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		// Should fail
		require.Error(t, err, "list with args should fail")

		// Should show error about unknown command or too many args
		assert.Contains(t, outputStr, "unknown command")
	})

	t.Run("list help shows usage", func(t *testing.T) {
		cmd := exec.Command(setup.binaryPath, "list", "--help")
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		require.NoError(t, err, "list --help should succeed: %s", outputStr)

		assert.Contains(t, outputStr, "List all worktrees managed by gh-wt")
		assert.Contains(t, outputStr, "gh wt list")
		assert.Contains(t, outputStr, "gh wt ls")
	})
}

func TestListWithWorktreesIntegration(t *testing.T) {
	setup := setupIntegrationTest(t)

	// Create a temporary directory structure for worktrees
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo")
	worktreeBase := filepath.Join(tmpDir, "worktrees")

	// Initialize a git repo
	require.NoError(t, os.MkdirAll(repoDir, 0o755))

	gitInit := exec.Command("git", "init")
	gitInit.Dir = repoDir
	require.NoError(t, gitInit.Run(), "git init failed")

	gitConfig1 := exec.Command("git", "config", "user.email", "test@example.com")
	gitConfig1.Dir = repoDir
	require.NoError(t, gitConfig1.Run())

	gitConfig2 := exec.Command("git", "config", "user.name", "Test User")
	gitConfig2.Dir = repoDir
	require.NoError(t, gitConfig2.Run())

	// Create initial commit
	testFile := filepath.Join(repoDir, "test.txt")
	require.NoError(t, os.WriteFile(testFile, []byte("test"), 0o644))

	gitAdd := exec.Command("git", "add", "test.txt")
	gitAdd.Dir = repoDir
	require.NoError(t, gitAdd.Run())

	gitCommit := exec.Command("git", "commit", "-m", "initial commit")
	gitCommit.Dir = repoDir
	require.NoError(t, gitCommit.Run())

	// Create worktrees under the worktree base
	wt1Path := filepath.Join(worktreeBase, "repo", "feature-1")
	gitWt1 := exec.Command("git", "worktree", "add", "-b", "feature-1", wt1Path)
	gitWt1.Dir = repoDir
	require.NoError(t, gitWt1.Run(), "failed to create worktree 1")

	wt2Path := filepath.Join(worktreeBase, "repo", "feature-2")
	gitWt2 := exec.Command("git", "worktree", "add", "-b", "feature-2", wt2Path)
	gitWt2.Dir = repoDir
	require.NoError(t, gitWt2.Run(), "failed to create worktree 2")

	t.Run("list shows managed worktrees", func(t *testing.T) {
		cmd := exec.Command(setup.binaryPath, "list")
		cmd.Dir = repoDir
		cmd.Env = append(os.Environ(), "GH_WT_WORKTREE_DIR="+worktreeBase)
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		require.NoError(t, err, "list should succeed: %s", outputStr)

		// Should show both worktrees
		assert.Contains(t, outputStr, "feature-1")
		assert.Contains(t, outputStr, "feature-2")

		// Should NOT contain "No worktrees found"
		assert.NotContains(t, outputStr, "No worktrees found")
	})

	t.Run("list with no-color flag", func(t *testing.T) {
		cmd := exec.Command(setup.binaryPath, "list", "--no-color")
		cmd.Dir = repoDir
		cmd.Env = append(os.Environ(), "GH_WT_WORKTREE_DIR="+worktreeBase)
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		require.NoError(t, err, "list --no-color should succeed: %s", outputStr)

		// Should not contain ANSI escape codes
		assert.NotContains(t, outputStr, "\x1b[")

		// Should still show worktree info
		assert.Contains(t, outputStr, "feature-1")
		assert.Contains(t, outputStr, "feature-2")
	})

	// Cleanup worktrees
	gitWtRm1 := exec.Command("git", "worktree", "remove", wt1Path, "--force")
	gitWtRm1.Dir = repoDir
	_ = gitWtRm1.Run()

	gitWtRm2 := exec.Command("git", "worktree", "remove", wt2Path, "--force")
	gitWtRm2.Dir = repoDir
	_ = gitWtRm2.Run()
}

func TestListAllIntegration(t *testing.T) {
	setup := setupIntegrationTest(t)

	tmpDir := t.TempDir()
	worktreeBase := filepath.Join(tmpDir, "worktrees")

	// Create two separate git repos and worktrees under the same base
	for _, repoName := range []string{"repo-a", "repo-b"} {
		repoDir := filepath.Join(tmpDir, repoName)
		require.NoError(t, os.MkdirAll(repoDir, 0o755))

		gitInit := exec.Command("git", "init")
		gitInit.Dir = repoDir
		require.NoError(t, gitInit.Run(), "git init failed for %s", repoName)

		gitCfg1 := exec.Command("git", "config", "user.email", "test@example.com")
		gitCfg1.Dir = repoDir
		require.NoError(t, gitCfg1.Run())

		gitCfg2 := exec.Command("git", "config", "user.name", "Test User")
		gitCfg2.Dir = repoDir
		require.NoError(t, gitCfg2.Run())

		testFile := filepath.Join(repoDir, "test.txt")
		require.NoError(t, os.WriteFile(testFile, []byte("test"), 0o644))

		gitAdd := exec.Command("git", "add", "test.txt")
		gitAdd.Dir = repoDir
		require.NoError(t, gitAdd.Run())

		gitCommit := exec.Command("git", "commit", "-m", "initial commit")
		gitCommit.Dir = repoDir
		require.NoError(t, gitCommit.Run())

		// Create a worktree for each repo
		wtPath := filepath.Join(worktreeBase, repoName, "feature-1")
		gitWt := exec.Command("git", "worktree", "add", "-b", "feature-1", wtPath)
		gitWt.Dir = repoDir
		require.NoError(t, gitWt.Run(), "failed to create worktree for %s", repoName)
	}

	t.Run("list --all shows worktrees grouped by repo", func(t *testing.T) {
		cmd := exec.Command(setup.binaryPath, "list", "--all", "--no-color")
		cmd.Env = append(os.Environ(), "GH_WT_WORKTREE_DIR="+worktreeBase)
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		require.NoError(t, err, "list --all should succeed: %s", outputStr)

		// Should contain both repo names as group headers
		assert.Contains(t, outputStr, "repo-a")
		assert.Contains(t, outputStr, "repo-b")

		// Should contain worktree names
		assert.Contains(t, outputStr, "feature-1")

		// Should contain indented headers
		assert.Contains(t, outputStr, "  NAME")
		assert.Contains(t, outputStr, "BRANCH")

		// Verify repo-a appears before repo-b (filesystem order)
		idxA := strings.Index(outputStr, "repo-a")
		idxB := strings.Index(outputStr, "repo-b")
		assert.Greater(t, idxB, idxA, "repo-a should appear before repo-b")
	})

	t.Run("list --all with -a shorthand", func(t *testing.T) {
		cmd := exec.Command(setup.binaryPath, "list", "-a", "--no-color")
		cmd.Env = append(os.Environ(), "GH_WT_WORKTREE_DIR="+worktreeBase)
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		require.NoError(t, err, "list -a should succeed: %s", outputStr)
		assert.Contains(t, outputStr, "repo-a")
		assert.Contains(t, outputStr, "repo-b")
	})

	t.Run("list --all with no worktrees shows warning", func(t *testing.T) {
		emptyBase := filepath.Join(tmpDir, "empty-base")
		require.NoError(t, os.MkdirAll(emptyBase, 0o755))

		cmd := exec.Command(setup.binaryPath, "list", "--all")
		cmd.Env = append(os.Environ(), "GH_WT_WORKTREE_DIR="+emptyBase)
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		require.NoError(t, err, "list --all with empty base should succeed: %s", outputStr)
		assert.Contains(t, outputStr, "No worktrees found")
	})
}

func TestListVisibilityIntegration(t *testing.T) {
	setup := setupIntegrationTest(t)

	cmd := exec.Command(setup.binaryPath, "--help")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Help command should succeed: %s", string(output))

	outputStr := string(output)

	// List command should appear in help under Worktrees group
	assert.Contains(t, outputStr, "list")
	assert.Contains(t, outputStr, "List managed worktrees")

	// Verify it's in the Worktrees section
	worktreesIdx := strings.Index(outputStr, "Worktrees")
	listIdx := strings.Index(outputStr, "list")
	assert.Greater(t, listIdx, worktreesIdx, "list should appear after Worktrees header")
}

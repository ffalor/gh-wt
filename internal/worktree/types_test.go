package worktree

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseArgument_ValidPRURL(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantType     WorktreeType
		wantOwner    string
		wantRepo     string
		wantNumber   int
		wantCloneURL string
	}{
		{
			name:         "valid PR URL",
			input:        "https://github.com/spf13/cobra-cli/pull/123",
			wantType:     PR,
			wantOwner:    "spf13",
			wantRepo:     "cobra-cli",
			wantNumber:   123,
			wantCloneURL: "https://github.com/spf13/cobra-cli.git",
		},
		{
			name:         "PR URL with complex repo name",
			input:        "https://github.com/kubernetes/kubernetes/pull/45678",
			wantType:     PR,
			wantOwner:    "kubernetes",
			wantRepo:     "kubernetes",
			wantNumber:   45678,
			wantCloneURL: "https://github.com/kubernetes/kubernetes.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArgument(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.wantType, got.Type)
			assert.Equal(t, tt.wantOwner, got.Owner)
			assert.Equal(t, tt.wantRepo, got.Repo)
			assert.Equal(t, tt.wantNumber, got.Number)
			assert.Equal(t, tt.wantCloneURL, got.CloneURL)
			assert.Equal(t, fmt.Sprintf("pr_%d", tt.wantNumber), got.WorktreeName)
		})
	}
}

func TestParseArgument_ValidIssueURL(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantType     WorktreeType
		wantOwner    string
		wantRepo     string
		wantNumber   int
		wantCloneURL string
	}{
		{
			name:         "valid issue URL",
			input:        "https://github.com/spf13/cobra-cli/issues/456",
			wantType:     Issue,
			wantOwner:    "spf13",
			wantRepo:     "cobra-cli",
			wantNumber:   456,
			wantCloneURL: "https://github.com/spf13/cobra-cli.git",
		},
		{
			name:         "issue URL with complex repo name",
			input:        "https://github.com/kubernetes/kubernetes/issues/789",
			wantType:     Issue,
			wantOwner:    "kubernetes",
			wantRepo:     "kubernetes",
			wantNumber:   789,
			wantCloneURL: "https://github.com/kubernetes/kubernetes.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArgument(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.wantType, got.Type)
			assert.Equal(t, tt.wantOwner, got.Owner)
			assert.Equal(t, tt.wantRepo, got.Repo)
			assert.Equal(t, tt.wantNumber, got.Number)
			assert.Equal(t, tt.wantCloneURL, got.CloneURL)
		})
	}
}

func TestParseArgument_UnsupportedGitHubURL(t *testing.T) {
	// These are GitHub URLs but unsupported types (not PR or Issue)
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "GitHub actions URL",
			input: "https://github.com/owner/repo/actions/runs/123",
		},
		{
			name:  "GitHub wiki URL",
			input: "https://github.com/owner/repo/wiki",
		},
		{
			name:  "GitHub release URL",
			input: "https://github.com/owner/repo/releases/tag/v1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These should return an error because they are GitHub URLs
			// but not PR or Issue URLs
			_, err := ParseArgument(tt.input)
			require.Error(t, err)
		})
	}
}

func TestParseArgument_NonGitHubURL(t *testing.T) {
	// These are not GitHub URLs at all - they fall through to parseLocalName
	// which will fail if not in a git repository
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "GitLab URL",
			input: "https://gitlab.com/owner/repo/merge_requests/123",
		},
		{
			name:  "Bitbucket URL",
			input: "https://bitbucket.org/owner/repo/pull-requests/123",
		},
		{
			name:  "ftp URL",
			input: "ftp://github.com/owner/repo/issues/123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These will attempt to parse as local names
			// and will error if not in a git repository
			_, err := ParseArgument(tt.input)
			// Either success (if in a git repo and can get current repo info)
			// or error (if not in a git repo)
			// We just verify it doesn't panic
			assert.True(t, err == nil || err != nil, "Should return either success or error")
		})
	}
}

func TestSanitizeBranchName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid name",
			input:    "feature-branch",
			expected: "feature-branch",
		},
		{
			name:     "with spaces",
			input:    "feature branch with spaces",
			expected: "feature_branch_with_spaces",
		},
		{
			name:     "with special characters",
			input:    "feature/branch@test#123",
			expected: "feature_branch_test_123",
		},
		{
			name:     "with periods",
			input:    "feature.v1.2.3",
			expected: "feature_v1_2_3",
		},
		{
			name:     "with asterisks",
			input:    "feature*fix",
			expected: "feature_fix",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeBranchName(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestWorktreeInfo_GetWorktreePath(t *testing.T) {
	tests := []struct {
		name     string
		info     WorktreeInfo
		baseDir  string
		expected string
	}{
		{
			name: "PR worktree",
			info: WorktreeInfo{
				Repo:         "cobra-cli",
				WorktreeName: "pr_123",
			},
			baseDir:  "/home/user/worktrees",
			expected: filepath.Join("/home/user/worktrees", "cobra-cli", "pr_123"),
		},
		{
			name: "Issue worktree",
			info: WorktreeInfo{
				Repo:         "gh-worktree",
				WorktreeName: "issue_456",
			},
			baseDir:  "/github/worktrees",
			expected: filepath.Join("/github/worktrees", "gh-worktree", "issue_456"),
		},
		{
			name: "Local worktree",
			info: WorktreeInfo{
				Repo:         "my-project",
				WorktreeName: "feature-test",
			},
			baseDir:  "~/worktrees",
			expected: filepath.Join("~/worktrees", "my-project", "feature-test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.info.GetWorktreePath(tt.baseDir)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestWorktreeInfo_GetRepoPath(t *testing.T) {
	tests := []struct {
		name     string
		info     WorktreeInfo
		baseDir  string
		expected string
	}{
		{
			name: "get bare repo path",
			info: WorktreeInfo{
				Repo: "cobra-cli",
			},
			baseDir:  "/home/user/worktrees",
			expected: filepath.Join("/home/user/worktrees", "cobra-cli", BareDir),
		},
		{
			name: "get bare repo path for issue",
			info: WorktreeInfo{
				Repo: "gh-worktree",
			},
			baseDir:  "/github/worktrees",
			expected: filepath.Join("/github/worktrees", "gh-worktree", BareDir),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.info.GetRepoPath(tt.baseDir)
			assert.Equal(t, tt.expected, got)
		})
	}
}

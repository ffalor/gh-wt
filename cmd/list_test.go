package cmd

import (
	"testing"

	"github.com/ffalor/gh-wt/internal/git"
)

func TestListCmd_Structure(t *testing.T) {
	if listCmd == nil {
		t.Fatal("listCmd is nil")
	}

	if listCmd.Use != "list" {
		t.Errorf("expected Use to be 'list', got %q", listCmd.Use)
	}

	if len(listCmd.Aliases) != 1 || listCmd.Aliases[0] != "ls" {
		t.Errorf("expected Aliases to be ['ls'], got %v", listCmd.Aliases)
	}

	if listCmd.GroupID != "worktrees" {
		t.Errorf("expected GroupID to be 'worktrees', got %q", listCmd.GroupID)
	}

	if listCmd.Args == nil {
		t.Error("expected Args to be defined")
	}

	if listCmd.RunE == nil {
		t.Error("expected RunE to be defined")
	}
}

func TestListCmd_AllFlag(t *testing.T) {
	flag := listCmd.Flags().Lookup("all")
	if flag == nil {
		t.Fatal("expected --all flag to be defined")
	}

	if flag.Shorthand != "a" {
		t.Errorf("expected shorthand 'a', got %q", flag.Shorthand)
	}

	if flag.DefValue != "false" {
		t.Errorf("expected default value 'false', got %q", flag.DefValue)
	}
}

func TestFilterWorktreesByBase(t *testing.T) {
	tests := []struct {
		name        string
		worktrees   []git.WorktreeInfo
		base        string
		expectedLen int
	}{
		{
			name: "filters by prefix",
			worktrees: []git.WorktreeInfo{
				{Path: "/home/user/worktrees/repo/feature-1", Branch: "feature-1"},
				{Path: "/home/user/worktrees/repo/feature-2", Branch: "feature-2"},
				{Path: "/home/user/other/repo", Branch: "main"},
			},
			base:        "/home/user/worktrees",
			expectedLen: 2,
		},
		{
			name: "no matches",
			worktrees: []git.WorktreeInfo{
				{Path: "/home/user/other/repo", Branch: "main"},
			},
			base:        "/home/user/worktrees",
			expectedLen: 0,
		},
		{
			name: "all match",
			worktrees: []git.WorktreeInfo{
				{Path: "/home/user/worktrees/repo/feature-1", Branch: "feature-1"},
				{Path: "/home/user/worktrees/repo/feature-2", Branch: "feature-2"},
			},
			base:        "/home/user/worktrees",
			expectedLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterWorktreesByBase(tt.worktrees, tt.base)
			if len(result) != tt.expectedLen {
				t.Errorf("expected %d worktrees, got %d", tt.expectedLen, len(result))
			}
		})
	}
}

func TestGetWorktreeDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "standard path",
			path:     "/home/user/worktrees/repo/feature-branch",
			expected: "repo/feature-branch",
		},
		{
			name:     "short path",
			path:     "/feature",
			expected: "/feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getWorktreeDisplayName(tt.path)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGroupWorktreesByRepo(t *testing.T) {
	tests := []struct {
		name      string
		worktrees []git.WorktreeInfo
		baseDir   string
		expected  []repoGroup
	}{
		{
			name: "groups by repo",
			worktrees: []git.WorktreeInfo{
				{Path: "/base/repo-a/feature-1", Branch: "feature-1"},
				{Path: "/base/repo-a/bugfix-2", Branch: "bugfix-2"},
				{Path: "/base/repo-b/testing123", Branch: "testing123"},
			},
			baseDir: "/base",
			expected: []repoGroup{
				{
					repo: "repo-a",
					worktrees: []git.WorktreeInfo{
						{Path: "/base/repo-a/feature-1", Branch: "feature-1"},
						{Path: "/base/repo-a/bugfix-2", Branch: "bugfix-2"},
					},
				},
				{
					repo: "repo-b",
					worktrees: []git.WorktreeInfo{
						{Path: "/base/repo-b/testing123", Branch: "testing123"},
					},
				},
			},
		},
		{
			name:      "empty input",
			worktrees: nil,
			baseDir:   "/base",
			expected:  nil,
		},
		{
			name: "single repo",
			worktrees: []git.WorktreeInfo{
				{Path: "/base/my-repo/wt-1", Branch: "wt-1"},
				{Path: "/base/my-repo/wt-2", Branch: "wt-2"},
			},
			baseDir: "/base",
			expected: []repoGroup{
				{
					repo: "my-repo",
					worktrees: []git.WorktreeInfo{
						{Path: "/base/my-repo/wt-1", Branch: "wt-1"},
						{Path: "/base/my-repo/wt-2", Branch: "wt-2"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := groupWorktreesByRepo(tt.worktrees, tt.baseDir)

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d groups, got %d", len(tt.expected), len(result))
			}

			for i, group := range result {
				if group.repo != tt.expected[i].repo {
					t.Errorf("group %d: expected repo %q, got %q", i, tt.expected[i].repo, group.repo)
				}
				if len(group.worktrees) != len(tt.expected[i].worktrees) {
					t.Errorf("group %d: expected %d worktrees, got %d",
						i, len(tt.expected[i].worktrees), len(group.worktrees))
					continue
				}
				for j, wt := range group.worktrees {
					if wt.Path != tt.expected[i].worktrees[j].Path {
						t.Errorf("group %d, wt %d: expected path %q, got %q",
							i, j, tt.expected[i].worktrees[j].Path, wt.Path)
					}
					if wt.Branch != tt.expected[i].worktrees[j].Branch {
						t.Errorf("group %d, wt %d: expected branch %q, got %q",
							i, j, tt.expected[i].worktrees[j].Branch, wt.Branch)
					}
				}
			}
		})
	}
}

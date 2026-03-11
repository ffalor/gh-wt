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
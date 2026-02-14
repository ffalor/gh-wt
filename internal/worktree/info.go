package worktree

// WorktreeType represents the type of worktree being created.
type WorktreeType string

const (
	Issue WorktreeType = "issue"
	PR    WorktreeType = "pr"
	Local WorktreeType = "local"
)

// WorktreeInfo contains metadata used during creation and action templating.
type WorktreeInfo struct {
	Type         WorktreeType
	Owner        string
	Repo         string
	Number       int
	BranchName   string
	WorktreeName string
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ffalor/gh-worktree/internal/config"
	"github.com/ffalor/gh-worktree/internal/worktree"
	"github.com/spf13/cobra"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type worktreeItem struct {
	worktree.WorktreeListItem
}

func (i worktreeItem) FilterValue() string { return i.Name }

type model struct {
	list     list.Model
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			m.quitting = true
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			selected := m.list.SelectedItem().(worktreeItem)
			fmt.Printf("\ncd %s\n", selected.Path)
			return m, tea.Quit
		}
		if msg.String() == "d" {
			selected := m.list.SelectedItem().(worktreeItem)
			fmt.Printf("\nDelete worktree: %s\n", selected.Name)
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	return docStyle.Render(m.list.View())
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [repo]",
	Short: "List all worktrees",
	Long:  `List all worktrees for a repository with an interactive interface.`,
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	baseDir := config.GetWorktreeBase()

	// Determine which repo to list
	var repoPath string
	if len(args) > 0 {
		repoPath = filepath.Join(baseDir, args[0], worktree.BareDir)
	} else {
		// Try to find repos
		entries, err := os.ReadDir(baseDir)
		if err != nil {
			return fmt.Errorf("no worktrees found in %s", baseDir)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				possiblePath := filepath.Join(baseDir, entry.Name(), worktree.BareDir)
				if _, err := os.Stat(possiblePath); err == nil {
					repoPath = possiblePath
					break
				}
			}
		}
	}

	if repoPath == "" {
		return fmt.Errorf("no repositories found in %s", baseDir)
	}

	items, err := worktree.List(repoPath)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		fmt.Println("No worktrees found")
		return nil
	}

	// Convert to list items
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = worktreeItem{WorktreeListItem: item}
	}

	// Check if we're in a TTY
	if !isTTY() {
		// Simple text output
		for _, item := range items {
			status := "clean"
			if item.HasChanges {
				status = "modified"
			}
			fmt.Printf("%-20s %-30s %-10s %s\n",
				item.Name,
				item.Branch,
				status,
				time.Unix(item.LastModTime, 0).Format("2006-01-02 15:04"))
		}
		return nil
	}

	// Interactive bubble tea UI
	delegate := list.NewDefaultDelegate()

	// Customize styles based on worktree type and status
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 0, 0, 2)

	delegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AFAFAF"))

	l := list.New(listItems, delegate, 0, 0)
	l.Title = "Worktrees"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#25A065")).
		Padding(0, 1)

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "cd to worktree")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		}
	}

	m := model{list: l}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

func isTTY() bool {
	fi, _ := os.Stdout.Stat()
	return fi.Mode()&os.ModeCharDevice != 0
}

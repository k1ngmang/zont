package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type listItem string

func (i listItem) FilterValue() string { return string(i) }
func (i listItem) Title() string       { return string(i) }
func (i listItem) Description() string { return "base model" }

func createList(items []list.Item, title string) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), 40, 20)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)
	return l
}

func renderAssetsTab(m model) string {
	listView := m.modelList.View()

	previewStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(25).
		Height(20)

	previewContent := previewStyle.Render(m.previewContent)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		listView,
		"  ",
		previewContent,
	)
}

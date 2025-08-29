package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderTabs() string {
	var tabs []string
	for i, tab := range m.tabs {
		style := lipgloss.NewStyle().Padding(0, 2)
		if i == m.activeTab {
			style = style.
				Foreground(lipgloss.Color("205")).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("205"))
		} else {
			style = style.Foreground(lipgloss.Color("240"))
		}
		tabs = append(tabs, style.Render(tab))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func renderTabsContent(m model) string {
	var content string

	switch m.activeTab {
	case 0:
		content = renderWelcomeTab()
	case 1:
		content = renderAssetsTab(m)
	case 2:
		content = renderConfigTab(m)
	}

	status := fmt.Sprintf("Selected file: %s", m.selectedFile)
	if m.selectedFile == "" {
		status = "File not selected"
	}

	statusBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1).
		Render(status)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.renderTabs(),
		content,
		statusBar,
	)
}

func renderWelcomeTab() string {
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render("Welcome to Zont!\n\n" +
			"In the 'Assets' tab you can find pre-made models or load your own.\n\n" +
			"Control:\n" +
			"←/→ - switch tabs\n" +
			"↑/↓ - list navigation\n" +
			"Enter - model selection/loading\n" +
			"q - exit")
}

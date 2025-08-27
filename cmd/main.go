package main

import (
	"fmt"
	"log"
	"os"
	"zontengine/internal/matrix"
	"zontengine/internal/render"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func RenderFromFile(file string) {
	matrix := matrix.NewMatrix(20, 20)
	renderer := render.NewRender(matrix)

	verts, err := render.LoadOBJ(file)
	if err != nil {
		log.Fatalf("Error %s: %v", file, err)
	}

	renderer.Render(verts)
}

type model struct {
	activeTab        int
	tabs             []string
	modelList        list.Model
	selectedFile     string
	previewContent   string
	lastSelectedItem string
	renderingMode    bool
}

type listItem string

func (i listItem) FilterValue() string { return string(i) }
func (i listItem) Title() string       { return string(i) }
func (i listItem) Description() string { return "base model" }

func initialModel() model {
	tabs := []string{"Zont", "Assets", "Render"}

	items := []list.Item{
		listItem("cube.obj"),
		listItem("monkey.obj"),
		listItem("Render Selected Model"),
	}

	l := list.New(items, list.NewDefaultDelegate(), 40, 20)
	l.Title = "Select model"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	return model{
		activeTab:      1,
		tabs:           tabs,
		modelList:      l,
		selectedFile:   "",
		previewContent: "",
		renderingMode:  false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	prevSelected := m.lastSelectedItem

	if m.renderingMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "esc", " ":
				m.renderingMode = false
				return m, tea.ClearScreen
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "left", "h":
			if m.activeTab > 0 {
				m.activeTab--
			}

		case "right", "l":
			if m.activeTab < len(m.tabs)-1 {
				m.activeTab++
			}

		case "enter":
			if m.activeTab == 1 {
				if selected, ok := m.modelList.SelectedItem().(listItem); ok {
					selectedStr := string(selected)

					if selectedStr == "Render Selected Model" {
						if m.selectedFile != "" {
							m.renderingMode = true
							return m, tea.Sequence(
								tea.ExitAltScreen,
								tea.ClearScreen,
								func() tea.Msg {
									RenderFromFile("models/" + m.selectedFile)
									return nil
								},
								tea.EnterAltScreen,
							)
						}
					} else {
						m.selectedFile = selectedStr
						fmt.Printf("File selected: %s\n", m.selectedFile)
					}
				}
			}
		}
	}

	if m.activeTab == 1 && !m.renderingMode {
		m.modelList, cmd = m.modelList.Update(msg)

		if selected, ok := m.modelList.SelectedItem().(listItem); ok {
			currentSelected := string(selected)
			if currentSelected != prevSelected && currentSelected != "Render Selected Model" {
				m.lastSelectedItem = currentSelected
				m.previewContent = m.renderModelFace(currentSelected)
			}
		}
	}

	return m, cmd
}

func (m model) renderModelFace(filename string) string {
	if filename == "" {
		return "Select a model to preview"
	}

	file := "models/" + filename
	matrix := matrix.NewMatrix(20, 20)
	renderer := render.NewRender(matrix)

	verts, err := render.LoadOBJ(file)
	if err != nil {
		return fmt.Sprintf("Loading error: %v", err)
	}

	str := renderer.RenderFrontFace(verts)
	return str
}

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

func (m model) View() string {
	if m.renderingMode {
		return ""
	}

	var content string

	switch m.activeTab {
	case 0:
		content = lipgloss.NewStyle().
			Padding(1, 2).
			Render("Welcome to Zont!\n\n" +
				"In the 'Assets' tab you can find pre-made models.\n\n" +
				"Control:\n" +
				"←/→ - switch tabs\n" +
				"↑/↓ - list navigation\n" +
				"Enter - model selection/rendering\n" +
				"q - exit")

	case 1:
		listView := m.modelList.View()

		previewStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1).
			Width(25).
			Height(20)

		previewContent := previewStyle.Render(m.previewContent)

		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			listView,
			"  ",
			previewContent,
		)

	case 2:
		content = lipgloss.NewStyle().
			Padding(1, 2).
			Render("Settings\n\n" +
				"TODO :).")
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

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

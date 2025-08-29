package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"zontengine/internal/config"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type textInput struct {
	textinput.Model
	label string
}

func createTextInput(placeholder, label string) textInput {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 3
	ti.Width = 10
	return textInput{Model: ti, label: label}
}

func handleRenderingMode(m model, msg tea.Msg) (model, tea.Cmd) {
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

func handleCustomFileLoading(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.customFilePath != "" {
				return processCustomFile(m)
			}
		case "esc":
			m.loadingCustom = false
			return m, nil
		default:
			if len(msg.String()) == 1 {
				m.customFilePath += msg.String()
			}
		}
	}
	return m, nil
}

func processCustomFile(m model) (model, tea.Cmd) {
	if _, err := os.Stat(m.customFilePath); os.IsNotExist(err) {
		m.configError = fmt.Sprintf("File not found: %s", m.customFilePath)
		m.loadingCustom = false
		return m, nil
	}

	if !strings.HasSuffix(m.customFilePath, ".obj") {
		m.configError = "Only .obj files are supported"
		m.loadingCustom = false
		return m, nil
	}

	filename := filepath.Base(m.customFilePath)
	m.selectedFile = filename

	modelsPath := "models/" + filename
	if _, err := os.Stat(modelsPath); os.IsNotExist(err) {
		data, err := os.ReadFile(m.customFilePath)
		if err != nil {
			m.configError = fmt.Sprintf("Error reading file: %v", err)
			m.loadingCustom = false
			return m, nil
		}

		err = os.WriteFile(modelsPath, data, 0644)
		if err != nil {
			m.configError = fmt.Sprintf("Error copying file: %v", err)
			m.loadingCustom = false
			return m, nil
		}
	}

	if width, height, err := m.getConfigValues(); err == nil {
		config.Save(width, height, filename)
	}

	m.loadingCustom = false
	m.configError = ""
	return m, nil
}

func handleConfigTab(m model, msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.widthInput.Focused() {
				m.widthInput.Blur()
				m.heightInput.Focus()
			} else {
				m.heightInput.Blur()
				m.widthInput.Focus()
			}
		case "enter":
			width, height, err := m.getConfigValues()
			if err != nil {
				m.configError = err.Error()
				m.configSaved = false
				return m, nil
			}

			err = config.Save(width, height, m.selectedFile)
			if err != nil {
				m.configError = fmt.Sprintf("Save error: %v", err)
				m.configSaved = false
			} else {
				m.configError = ""
				m.configSaved = true
			}
		}
	}

	if m.widthInput.Focused() {
		m.widthInput.Model, cmd = m.widthInput.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.heightInput.Focused() {
		m.heightInput.Model, cmd = m.heightInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func handleNavigation(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left", "h":
			if m.activeTab > 0 {
				m.activeTab--
				m.widthInput.Blur()
				m.heightInput.Blur()
				m.configSaved = false
			}
		case "right", "l":
			if m.activeTab < len(m.tabs)-1 {
				m.activeTab++
				m.widthInput.Blur()
				m.heightInput.Blur()
				m.configSaved = false
			}
		case "enter":
			if m.activeTab == 1 {
				if selected, ok := m.modelList.SelectedItem().(listItem); ok {
					selectedStr := string(selected)
					if selectedStr == "Load .obj file" {
						m.loadingCustom = true
						m.customFilePath = ""
						return m, nil
					} else {
						m.selectedFile = selectedStr
						if width, height, err := m.getConfigValues(); err == nil {
							config.Save(width, height, m.selectedFile)
						}
					}
				}
			}
		}
	}
	return m, nil
}

func renderCustomFileLoading(m model) string {
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(fmt.Sprintf(
			"Enter path to .obj file:\n\n%s\n\nPress Enter to confirm, Esc to cancel\n\n%s",
			m.customFilePath,
			m.configError,
		))
}

func renderConfigTab(m model) string {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Width(10)

	inputStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(12)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("82")).
		Bold(true)

	modelInfo := "No model selected"
	if m.selectedFile != "" {
		modelInfo = fmt.Sprintf("Selected model: %s", m.selectedFile)
	}

	configContent := lipgloss.JoinVertical(
		lipgloss.Left,
		"Render Configuration:",
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			labelStyle.Render("Width:"),
			" ",
			inputStyle.Render(m.widthInput.View()),
		),
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			labelStyle.Render("Height:"),
			" ",
			inputStyle.Render(m.heightInput.View()),
		),
		"",
		modelInfo,
		"",
		"Controls:",
		"  Tab - switch between fields",
		"  Enter - save configuration",
		"",
	)

	if m.configError != "" {
		configContent = lipgloss.JoinVertical(
			lipgloss.Left,
			configContent,
			errorStyle.Render("Error: "+m.configError),
		)
	} else if m.configSaved {
		configContent = lipgloss.JoinVertical(
			lipgloss.Left,
			configContent,
			successStyle.Render("✓ Configuration saved successfully!"),
		)
	}

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(configContent)
}

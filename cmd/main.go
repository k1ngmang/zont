package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"zontengine/internal/matrix"
	"zontengine/internal/render"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	ModelFile string `json:"model_file"`
}

const configFileName = "render_config.json"

func RenderFromFile(file string, width, height int) {
	matrix := matrix.NewMatrix(width, height)
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
	widthInput       textinput.Model
	heightInput      textinput.Model
	configError      string
	configSaved      bool
	loadingCustom    bool
	customFilePath   string
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
		listItem("Load .obj file"), // Заменено на загрузку пользовательского файла
	}

	l := list.New(items, list.NewDefaultDelegate(), 40, 20)
	l.Title = "Select model"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	config := loadConfig()

	widthInput := textinput.New()
	widthInput.Placeholder = "20"
	widthInput.CharLimit = 3
	widthInput.Width = 10
	widthInput.SetValue(fmt.Sprintf("%d", config.Width))
	widthInput.Focus()

	heightInput := textinput.New()
	heightInput.Placeholder = "20"
	heightInput.CharLimit = 3
	heightInput.Width = 10
	heightInput.SetValue(fmt.Sprintf("%d", config.Height))

	selectedFile := config.ModelFile

	return model{
		activeTab:      1,
		tabs:           tabs,
		modelList:      l,
		selectedFile:   selectedFile,
		previewContent: "",
		renderingMode:  false,
		widthInput:     widthInput,
		heightInput:    heightInput,
		configError:    "",
		configSaved:    false,
		loadingCustom:  false,
		customFilePath: "",
	}
}

func loadConfig() Config {
	config := Config{Width: 20, Height: 20, ModelFile: ""}

	data, err := os.ReadFile(configFileName)
	if err != nil {
		return config
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return Config{Width: 20, Height: 20, ModelFile: ""}
	}

	return config
}

func saveConfig(width, height int, modelFile string) error {
	config := Config{
		Width:     width,
		Height:    height,
		ModelFile: modelFile,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFileName, data, 0644)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

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

	if m.loadingCustom {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				if m.customFilePath != "" {
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
						saveConfig(width, height, filename)
					}

					m.loadingCustom = false
					m.configError = ""
					return m, nil
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

	if m.activeTab == 2 {
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

				err = saveConfig(width, height, m.selectedFile)
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
			m.widthInput, cmd = m.widthInput.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.heightInput.Focused() {
			m.heightInput, cmd = m.heightInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

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
						fmt.Printf("File selected: %s\n", m.selectedFile)
						if width, height, err := m.getConfigValues(); err == nil {
							saveConfig(width, height, m.selectedFile)
						}
					}
				}
			}
		}
	}

	if m.activeTab == 1 && !m.renderingMode && !m.loadingCustom {
		m.modelList, cmd = m.modelList.Update(msg)
		cmds = append(cmds, cmd)

		if selected, ok := m.modelList.SelectedItem().(listItem); ok {
			currentSelected := string(selected)
			if currentSelected != prevSelected && currentSelected != "Load .obj file" {
				m.lastSelectedItem = currentSelected
				m.previewContent = m.renderModelFace(currentSelected)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) getConfigValues() (int, int, error) {
	widthStr := m.widthInput.Value()
	if widthStr == "" {
		widthStr = m.widthInput.Placeholder
	}
	heightStr := m.heightInput.Value()
	if heightStr == "" {
		heightStr = m.heightInput.Placeholder
	}

	width, err := strconv.Atoi(widthStr)
	if err != nil || width <= 0 {
		return 0, 0, fmt.Errorf("invalid width: must be positive integer")
	}

	height, err := strconv.Atoi(heightStr)
	if err != nil || height <= 0 {
		return 0, 0, fmt.Errorf("invalid height: must be positive integer")
	}

	return width, height, nil
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

	if m.loadingCustom {
		return lipgloss.NewStyle().
			Padding(1, 2).
			Render(fmt.Sprintf(
				"Enter path to .obj file:\n\n%s\n\nPress Enter to confirm, Esc to cancel\n\n%s",
				m.customFilePath,
				m.configError,
			))
	}

	var content string

	switch m.activeTab {
	case 0:
		content = lipgloss.NewStyle().
			Padding(1, 2).
			Render("Welcome to Zont!\n\n" +
				"In the 'Assets' tab you can find pre-made models or load your own.\n\n" +
				"Control:\n" +
				"←/→ - switch tabs\n" +
				"↑/↓ - list navigation\n" +
				"Enter - model selection/loading\n" +
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

		content = lipgloss.NewStyle().
			Padding(1, 2).
			Render(configContent)
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
	if len(os.Args) > 1 && os.Args[1] == "render" {
		config := loadConfig()

		if config.ModelFile == "" {
			fmt.Println("Error: No model selected in configuration")
			fmt.Println("Please run the program without 'render' flag first to select a model")
			os.Exit(1)
		}

		modelFile := "models/" + config.ModelFile
		if _, err := os.Stat(modelFile); os.IsNotExist(err) {
			fmt.Printf("Model file not found: %s\n", modelFile)
			fmt.Printf("Please check if '%s' exists in models/ directory\n", config.ModelFile)
			os.Exit(1)
		}

		fmt.Printf("Rendering %s with config: width=%d, height=%d\n",
			config.ModelFile, config.Width, config.Height)

		RenderFromFile(modelFile, config.Width, config.Height)
		return
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

package tui

import (
	"fmt"
	"os"
	"strconv"

	"zontengine/internal/config"
	"zontengine/internal/matrix"
	"zontengine/internal/render"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	activeTab        int
	tabs             []string
	modelList        list.Model
	selectedFile     string
	previewContent   string
	lastSelectedItem string
	renderingMode    bool
	widthInput       textInput
	heightInput      textInput
	configError      string
	configSaved      bool
	loadingCustom    bool
	customFilePath   string
}

func initialModel() model {
	tabs := []string{"Zont", "Assets", "Render"}

	items := []list.Item{
		listItem("cube.obj"),
		listItem("monkey.obj"),
		listItem("Load .obj file"),
	}

	l := createList(items, "Select model")

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Warning: %v\n", err)
		cfg = config.Config{Width: 20, Height: 20, ModelFile: ""}
	}

	widthInput := createTextInput("20", "Width")
	widthInput.SetValue(fmt.Sprintf("%d", cfg.Width))
	widthInput.Focus()

	heightInput := createTextInput("20", "Height")
	heightInput.SetValue(fmt.Sprintf("%d", cfg.Height))

	return model{
		activeTab:      1,
		tabs:           tabs,
		modelList:      l,
		selectedFile:   cfg.ModelFile,
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

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	prevSelected := m.lastSelectedItem

	if m.renderingMode {
		return handleRenderingMode(m, msg)
	}

	if m.loadingCustom {
		return handleCustomFileLoading(m, msg)
	}

	if m.activeTab == 2 {
		return handleConfigTab(m, msg)
	}

	m, cmd = handleNavigation(m, msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
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

func (m model) View() string {
	if m.renderingMode {
		return ""
	}

	if m.loadingCustom {
		return renderCustomFileLoading(m)
	}

	return renderTabsContent(m)
}

func Run() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

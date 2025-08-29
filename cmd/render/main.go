package main

import (
	"fmt"
	"log"
	"os"

	"zontengine/internal/config"
	"zontengine/internal/matrix"
	"zontengine/internal/render"
	"zontengine/internal/tui"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "render" {
		err := renderFromConfig()
		if err != nil {
			log.Fatal("Render error: ", err)
		}
	} else {
		log.Fatal("Incorrect arguments. Usage: program render")
	}
	tui.Run()
}

func renderFromConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if cfg.ModelFile == "" {
		return fmt.Errorf("no model selected in configuration - run TUI interface first")
	}

	modelFile := "models/" + cfg.ModelFile
	if _, err := os.Stat(modelFile); os.IsNotExist(err) {
		return fmt.Errorf("model file not found: %s (check models/ directory)", modelFile)
	}

	matrix := matrix.NewMatrix(cfg.Width, cfg.Height)
	renderer := render.NewRender(matrix)

	verts, err := render.LoadOBJ(modelFile)
	if err != nil {
		return fmt.Errorf("loading OBJ %s: %w", modelFile, err)
	}

	renderer.Render(verts)
	return nil
}

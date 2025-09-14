package config

/**
* Handles configuration management for the 3D rendering application.
* Provides functionality to load and save rendering settings to
* a JSON configuration file.
 */

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = "render_config.json"

type Config struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	ModelFile string `json:"model_file"`
}

func Load() (Config, error) {
	config := Config{Width: 20, Height: 20, ModelFile: ""}

	data, err := os.ReadFile(configFileName)
	if err != nil {
		return config, nil
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{Width: 20, Height: 20, ModelFile: ""},
			fmt.Errorf("error reading config: %w", err)
	}

	return config, nil
}

func Save(width, height int, modelFile string) error {
	config := Config{
		Width:     width,
		Height:    height,
		ModelFile: modelFile,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	return os.WriteFile(configFileName, data, 0644)
}

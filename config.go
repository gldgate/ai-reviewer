package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ModelMapping       map[string]ModelConfig `yaml:"model_mapping"`
	GlobalInstructions string                 `yaml:"global_instructions"`
}

type ModelConfig struct {
	Provider              string  `yaml:"provider"`
	Model                 string  `yaml:"model"`
	MaxTokens             int     `yaml:"max_tokens"`
	InputPricePerMillion  float64 `yaml:"input_price_per_million"`
	OutputPricePerMillion float64 `yaml:"output_price_per_million"`
}

func LoadConfig(searchPaths []string, repo string) (*Config, error) {
	finalConfig := &Config{
		ModelMapping: make(map[string]ModelConfig),
	}

	found := false
	for _, base := range searchPaths {
		configPath := filepath.Join(base, ".ai-review", repo, "config.yaml")
		data, err := os.ReadFile(configPath)
		if err != nil {
			// Fallback to global config if repo-specific doesn't exist?
			// User said "the tool should look in the local directory .ai-review/gonka-ai/gonka/"
			// Let's try to look in the old location as well for backward compatibility?
			// "The local .ai-review should directory should be specific to a specific owner/repo."
			// If I just change it to look in repo specific, it matches the requirement.
			configPath = filepath.Join(base, ".ai-review/config.yaml")
			data, err = os.ReadFile(configPath)
			if err != nil {
				continue
			}
		}

		fmt.Printf("    -> Loading config from: %s\n", configPath)
		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			fmt.Printf("Warning: error parsing config at %s: %v\n", configPath, err)
			continue
		}

		found = true
		// Merge model mappings
		for k, v := range cfg.ModelMapping {
			finalConfig.ModelMapping[k] = v
		}
		if cfg.GlobalInstructions != "" {
			finalConfig.GlobalInstructions = cfg.GlobalInstructions
		}
	}

	if !found {
		return nil, os.ErrNotExist
	}

	return finalConfig, nil
}

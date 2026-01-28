package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
)

type Persona struct {
	ID             string   `yaml:"id"`
	ModelCategory  string   `yaml:"model_category"`
	MaxTokens      int      `yaml:"max_tokens"`
	PathFilters    []string `yaml:"path_filters"`
	ExcludeFilters []string `yaml:"exclude_filters"`
	Role           string   `yaml:"role"`  // reviewer (default) | explainer
	Stage          string   `yaml:"stage"` // pre | post
	Instructions   string
}

func LoadPersonas(searchPaths []string, repo string) ([]Persona, error) {
	personaMap := make(map[string]Persona)
	foundAny := false

	for _, base := range searchPaths {
		// Try repo-specific personas
		personaDir := filepath.Join(base, ".ai-review", repo, "personas")
		files, _ := filepath.Glob(filepath.Join(personaDir, "*.md"))

		// Also try global personas
		globalPersonaDir := filepath.Join(base, ".ai-review/personas")
		globalFiles, _ := filepath.Glob(filepath.Join(globalPersonaDir, "*.md"))

		allFiles := append(files, globalFiles...)
		if len(allFiles) > 0 {
			foundAny = true
		}

		for _, file := range allFiles {
			f, err := os.Open(file)
			if err != nil {
				fmt.Printf("Warning: could not open persona file %s: %v\n", file, err)
				continue
			}

			var p Persona
			rest, err := frontmatter.Parse(f, &p)
			f.Close()
			if err != nil {
				fmt.Printf("Warning: error parsing frontmatter in %s: %v\n", file, err)
				continue
			}
			p.Instructions = string(rest)
			if p.ID == "" {
				// Fallback to filename without extension as ID
				p.ID = strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
			}
			if p.Role == "" {
				p.Role = "reviewer"
			}
			personaMap[p.ID] = p
		}
	}

	if !foundAny {
		return nil, fmt.Errorf("no personas found in any of the search paths")
	}

	var personas []Persona
	for _, p := range personaMap {
		personas = append(personas, p)
	}

	return personas, nil
}

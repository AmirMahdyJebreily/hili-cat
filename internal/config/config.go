package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Default configuration file location
const DefaultConfigPath = "/etc/highlight/config.json"

// Config represents the structure of the configuration file
type Config struct {
	Languages map[string]Language `json:"languages"`
}

// Language represents the syntax highlighting rules for a specific language
type Language struct {
	Extensions []string          `json:"extensions"`
	Rules      []HighlightRule   `json:"rules"`
	Styles     map[string]string `json:"styles"`
}

// HighlightRule defines a pattern to match and the style to apply
type HighlightRule struct {
	Name    string `json:"name"`
	Pattern string `json:"pattern"`
	Style   string `json:"style"`
}

// Load loads the syntax highlighting configuration from a file
func Load(configPath string) (Config, error) {
	var config Config

	file, err := os.Open(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, fmt.Errorf("failed to parse config file: %v", err)
	}

	return config, nil
}

// EnsureExists creates a default config file if it doesn't exist
func EnsureExists(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create directory structure if needed
		dir := filepath.Dir(configPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}

		// Create a basic default configuration
		defaultConfig := Config{
			Languages: map[string]Language{
				"go": {
					Extensions: []string{"go"},
					Rules: []HighlightRule{
						{Name: "keywords", Pattern: `\b(func|package|import|var|const|type|struct|interface|map|chan|go|defer|if|else|switch|case|for|range|return|break|continue)\b`, Style: "keyword"},
						{Name: "strings", Pattern: `"[^"]*"`, Style: "string"},
						{Name: "comments", Pattern: `//.*|/\*[\s\S]*?\*/`, Style: "comment"},
						{Name: "numbers", Pattern: `\b\d+\b`, Style: "number"},
					},
					Styles: map[string]string{
						"keyword": "cyan",
						"string":  "green",
						"comment": "yellow",
						"number":  "magenta",
					},
				},
				"json": {
					Extensions: []string{"json"},
					Rules: []HighlightRule{
						{Name: "keys", Pattern: `"[^"]*"\s*:`, Style: "key"},
						{Name: "strings", Pattern: `:\s*"[^"]*"`, Style: "string"},
						{Name: "numbers", Pattern: `:\s*\d+`, Style: "number"},
						{Name: "booleans", Pattern: `:\s*(true|false|null)`, Style: "boolean"},
					},
					Styles: map[string]string{
						"key":     "cyan",
						"string":  "green",
						"number":  "magenta",
						"boolean": "yellow",
					},
				},
			},
		}

		file, err := os.Create(configPath)
		if err != nil {
			return fmt.Errorf("failed to create config file: %v", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(defaultConfig); err != nil {
			return fmt.Errorf("failed to write default config: %v", err)
		}
	}

	return nil
}

// DetectLanguage tries to determine the language based on file extension
func DetectLanguage(cfg Config, fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" {
		return ""
	}

	// Remove the dot from extension if present
	if ext[0] == '.' {
		ext = ext[1:]
	}

	for lang, language := range cfg.Languages {
		for _, supportedExt := range language.Extensions {
			if supportedExt == ext {
				return lang
			}
		}
	}

	return ""
}

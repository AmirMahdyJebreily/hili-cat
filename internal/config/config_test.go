package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectLanguage(t *testing.T) {
	cfg := Config{
		Languages: map[string]Language{
			"go": {
				Extensions: []string{"go"},
			},
			"json": {
				Extensions: []string{"json"},
			},
			"multi": {
				Extensions: []string{"js", "ts", "jsx"},
			},
		},
	}

	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "go file",
			filename: "main.go",
			want:     "go",
		},
		{
			name:     "json file",
			filename: "config.json",
			want:     "json",
		},
		{
			name:     "typescript file",
			filename: "app.ts",
			want:     "multi",
		},
		{
			name:     "unknown extension",
			filename: "text.txt",
			want:     "",
		},
		{
			name:     "no extension",
			filename: "README",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectLanguage(cfg, tt.filename); got != tt.want {
				t.Errorf("DetectLanguage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnsureExists(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testConfigPath := filepath.Join(tmpDir, "test_config.json")

	// Test creating a new config file
	err = EnsureExists(testConfigPath)
	if err != nil {
		t.Errorf("EnsureExists() error = %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(testConfigPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created")
	}

	// Test with existing file (should not error)
	err = EnsureExists(testConfigPath)
	if err != nil {
		t.Errorf("EnsureExists() with existing file error = %v", err)
	}
}

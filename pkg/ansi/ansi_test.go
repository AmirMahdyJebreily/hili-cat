package ansi

import (
	"strings"
	"testing"
)

func TestGetStyleCode(t *testing.T) {
	tests := []struct {
		name     string
		style    string
		expected string
	}{
		{"Valid style", "red", Red},
		{"Valid bold style", "bold", Bold},
		{"Valid bright style", "brightgreen", BrightGreen},
		{"Invalid style", "nonexistent", Reset},
		{"Empty style", "", Reset},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStyleCode(tt.style)
			if result != tt.expected {
				t.Errorf("GetStyleCode(%q) = %q, want %q", tt.style, result, tt.expected)
			}
		})
	}
}

func TestColorize(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		style    string
		expected string
	}{
		{"Red text", "test", "red", Red + "test" + Reset},
		{"Bold text", "test", "bold", Bold + "test" + Reset},
		{"Invalid style", "test", "nonexistent", Reset + "test" + Reset},
		{"Empty text", "", "red", Red + "" + Reset},
		{"Empty style", "test", "", Reset + "test" + Reset},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Colorize(tt.text, tt.style)
			if result != tt.expected {
				t.Errorf("Colorize(%q, %q) = %q, want %q",
					tt.text, tt.style, result, tt.expected)
			}
		})
	}
}

func TestStyleMapCompleteness(t *testing.T) {
	// Check that all constants have corresponding entries in the StyleMap
	expectedStyles := []struct {
		name string
		code string
	}{
		{"reset", Reset},
		{"bold", Bold},
		{"italic", Italic},
		{"underline", Underline},
		{"black", Black},
		{"red", Red},
		{"green", Green},
		{"yellow", Yellow},
		{"blue", Blue},
		{"magenta", Magenta},
		{"cyan", Cyan},
		{"white", White},
		{"brightblack", BrightBlack},
		{"brightred", BrightRed},
		{"brightgreen", BrightGreen},
		{"brightyellow", BrightYellow},
		{"brightblue", BrightBlue},
		{"brightmagenta", BrightMagenta},
		{"brightcyan", BrightCyan},
		{"brightwhite", BrightWhite},
	}

	for _, style := range expectedStyles {
		code, ok := StyleMap[style.name]
		if !ok {
			t.Errorf("StyleMap is missing entry for %q", style.name)
			continue
		}
		if code != style.code {
			t.Errorf("StyleMap[%q] = %q, want %q", style.name, code, style.code)
		}
	}

	// Check for any extra entries in StyleMap
	for name := range StyleMap {
		found := false
		for _, style := range expectedStyles {
			if name == style.name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("StyleMap has unexpected entry %q", name)
		}
	}
}

func TestCodeValidity(t *testing.T) {
	// Ensure all ANSI codes follow the correct format
	for name, code := range StyleMap {
		if !strings.HasPrefix(code, "\033[") || !strings.HasSuffix(code, "m") {
			t.Errorf("ANSI code for %q doesn't follow expected format: %q", name, code)
		}
	}
}

package highlighter

import (
	"regexp"
	"strings"
	"testing"
)

func TestHighlightLine(t *testing.T) {
	// Create a simple test highlighter
	h := &Highlighter{
		styles: map[string]string{
			"keyword": "cyan",
			"string":  "green",
		},
		rules: []CompiledRule{
			{
				Name:    "keyword",
				Pattern: regexp.MustCompile(`\b(func|package|import)\b`),
				Style:   "keyword",
			},
			{
				Name:    "string",
				Pattern: regexp.MustCompile(`"[^"]*"`),
				Style:   "string",
			},
		},
	}

	// Test highlighting
	line := `func main() { fmt.Println("hello world") }`
	highlighted := h.highlightLine(line)
	
	// Basic verification that highlighting was applied
	if highlighted == line {
		t.Error("No highlighting was applied")
	}
	
	// Check that ANSI codes are present
	if !strings.Contains(highlighted, "\033[") {
		t.Error("No ANSI codes found in highlighted output")
	}
}

func TestProcessContent(t *testing.T) {
	// Create a simple test highlighter
	h := &Highlighter{
		lineEnding: "\n",
		options: Options{
			NumberLines: true,
		},
		styles: map[string]string{
			"keyword": "cyan",
		},
		rules: []CompiledRule{
			{
				Name:    "keyword",
				Pattern: regexp.MustCompile(`\b(func|package|import)\b`),
				Style:   "keyword",
			},
		},
	}
	
	input := []byte("package main\n\nfunc main() {}\n")
	output := h.ProcessContent(input)
	
	// Check that line numbers were added
	if !strings.Contains(output, "    1") {
		t.Error("Line numbers not added as expected")
	}
	
	// Check that the right number of lines are in the output
	lines := strings.Split(output, "\n")
	// We expect 3 lines of content plus one empty line after splitting
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines in output, got %d", len(lines))
	}
}

func TestAnsiStyle(t *testing.T) {
	h := &Highlighter{}
	
	testCases := []struct {
		name     string
		style    string
		expected string
	}{
		{"valid color", "red", Red},
		{"valid formatting", "bold", Bold},
		{"invalid style", "nonexistent", ""},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := h.ansiStyle(tc.style)
			if result != tc.expected {
				t.Errorf("ansiStyle(%s) = %s, want %s", tc.style, result, tc.expected)
			}
		})
	}
}

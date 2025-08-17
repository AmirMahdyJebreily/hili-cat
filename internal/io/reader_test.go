package io

import (
	"testing"
)

func TestDetectLineEnding(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "empty input",
			input:    []byte{},
			expected: "\n",
		},
		{
			name:     "LF ending",
			input:    []byte("line1\nline2\n"),
			expected: "\n",
		},
		{
			name:     "CRLF ending",
			input:    []byte("line1\r\nline2\r\n"),
			expected: "\r\n",
		},
		{
			name:     "mixed endings with CRLF",
			input:    []byte("line1\nline2\r\n"),
			expected: "\r\n", // Should detect CRLF even if mixed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectLineEnding(tt.input)
			if result != tt.expected {
				t.Errorf("DetectLineEnding() got %q, want %q", result, tt.expected)
			}
		})
	}
}

// Note: Testing OpenFile and ProcessFile requires more sophisticated setup
// with mock syscalls, which we'd implement in a production environment

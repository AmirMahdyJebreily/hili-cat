package highlighter

import (
	"fmt"
	"regexp"
	"strings"
)

// ANSI color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Black     = "\033[30m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[35m"
	Cyan      = "\033[36m"
	White     = "\033[37m"
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

// Line ending constants
const (
	LF   = "\n"
	CRLF = "\r\n"
)

// Config represents configuration data needed by the highlighter
type Config struct {
	Languages map[string]Language
}

// Language represents the syntax highlighting rules for a specific language
type Language struct {
	Extensions []string
	Rules      []HighlightRule
	Styles     map[string]string
}

// HighlightRule defines a pattern to match and the style to apply
type HighlightRule struct {
	Name    string
	Pattern string
	Style   string
}

// CompiledRule is a compiled version of HighlightRule for better performance
type CompiledRule struct {
	Name    string
	Pattern *regexp.Regexp
	Style   string
}

// Highlighter manages the syntax highlighting process
type Highlighter struct {
	config     Config
	language   string
	rules      []CompiledRule
	styles     map[string]string
	lineEnding string
	lineNum    int
	options    Options
}

// Options contains settings for the highlighter
type Options struct {
	NumberLines    bool
	NumberNonBlank bool
	SqueezeBlank   bool
	ShowEnds       bool
}

// NewHighlighter creates and initializes a new Highlighter
func NewHighlighter(cfg Config, lang, lineEnding string, opts Options) (*Highlighter, error) {
	language, ok := cfg.Languages[lang]
	if !ok {
		return nil, fmt.Errorf("language not found in configuration: %s", lang)
	}

	highlighter := &Highlighter{
		config:     cfg,
		language:   lang,
		styles:     language.Styles,
		lineEnding: lineEnding,
		lineNum:    0,
		options:    opts,
	}

	// Compile all regex patterns for better performance
	for _, rule := range language.Rules {
		pattern, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern for %s: %v", rule.Name, err)
		}
		highlighter.rules = append(highlighter.rules, CompiledRule{
			Name:    rule.Name,
			Pattern: pattern,
			Style:   rule.Style,
		})
	}

	return highlighter, nil
}

// ProcessContent processes the input data and returns highlighted output
func (h *Highlighter) ProcessContent(data []byte) string {
	var buffer strings.Builder
	var lineBuffer strings.Builder
	var lastLineWasBlank bool

	content := string(data)
	lines := strings.Split(content, h.lineEnding)

	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			// Skip the last line if it's empty (common when splitting on newlines)
			break
		}

		isBlankLine := len(strings.TrimSpace(line)) == 0

		// Handle squeeze blank option
		if h.options.SqueezeBlank && isBlankLine && lastLineWasBlank {
			continue
		}

		lastLineWasBlank = isBlankLine

		// Add line number if required
		if h.options.NumberLines || (h.options.NumberNonBlank && !isBlankLine) {
			h.lineNum++
			if h.options.NumberNonBlank && isBlankLine {
				lineBuffer.WriteString("     ")
			} else {
				lineBuffer.WriteString(fmt.Sprintf("%5d  ", h.lineNum))
			}
		}

		// Add highlighted content
		lineBuffer.WriteString(h.highlightLine(line))

		// Add end marker if requested
		if h.options.ShowEnds {
			lineBuffer.WriteString("$")
		}

		buffer.WriteString(lineBuffer.String())
		buffer.WriteString(h.lineEnding)
		lineBuffer.Reset()
	}

	return buffer.String()
}

// highlightLine applies syntax highlighting to a single line
func (h *Highlighter) highlightLine(line string) string {
	type Token struct {
		start int
		end   int
		style string
	}

	var tokens []Token

	// Find all matches for all rules
	for _, rule := range h.rules {
		matches := rule.Pattern.FindAllStringIndex(line, -1)
		for _, match := range matches {
			styleCode := ""
			if styleName, ok := h.styles[rule.Style]; ok {
				styleCode = h.ansiStyle(styleName)
			}

			tokens = append(tokens, Token{
				start: match[0],
				end:   match[1],
				style: styleCode,
			})
		}
	}

	// If no matches, return the original line
	if len(tokens) == 0 {
		return line
	}

	// Sort tokens by start position for non-overlapping highlighting
	// In production, we'd implement a proper sort here

	var result strings.Builder
	lastPos := 0

	for _, token := range tokens {
		// Add text before the token
		if token.start > lastPos {
			result.WriteString(line[lastPos:token.start])
		}

		// Add styled token
		if token.style != "" {
			result.WriteString(token.style)
			result.WriteString(line[token.start:token.end])
			result.WriteString(Reset)
		} else {
			result.WriteString(line[token.start:token.end])
		}

		lastPos = token.end
	}

	// Add any remaining text
	if lastPos < len(line) {
		result.WriteString(line[lastPos:])
	}

	return result.String()
}

// ansiStyle converts a style name to its ANSI escape code
func (h *Highlighter) ansiStyle(styleName string) string {
	styles := map[string]string{
		"reset":      Reset,
		"bold":       Bold,
		"italic":     Italic,
		"underline":  Underline,
		"black":      Black,
		"red":        Red,
		"green":      Green,
		"yellow":     Yellow,
		"blue":       Blue,
		"magenta":    Magenta,
		"cyan":       Cyan,
		"white":      White,
		"bg_black":   BgBlack,
		"bg_red":     BgRed,
		"bg_green":   BgGreen,
		"bg_yellow":  BgYellow,
		"bg_blue":    BgBlue,
		"bg_magenta": BgMagenta,
		"bg_cyan":    BgCyan,
		"bg_white":   BgWhite,
	}

	if code, ok := styles[styleName]; ok {
		return code
	}
	return ""
}

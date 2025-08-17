# Language Configuration Guide for highlight

This guide explains how to create and customize language configurations for the `highlight` tool.

## Table of Contents
- [Understanding the Configuration File](#understanding-the-configuration-file)
- [Adding a New Language](#adding-a-new-language)
- [Regex Pattern Writing Tips](#regex-pattern-writing-tips)
- [Available Styles](#available-styles)
- [Testing Your Configuration](#testing-your-configuration)
- [Optimizing Regex Patterns](#optimizing-regex-patterns)
- [Examples](#examples)

## Understanding the Configuration File

The `highlight` tool uses a JSON configuration file to define syntax highlighting rules for different programming languages. The default configuration is located at `/etc/highlight/config.json`.

### Configuration Structure

```json
{
  "languages": {
    "language-name": {
      "extensions": ["ext1", "ext2"],
      "rules": [
        {
          "name": "rule-name",
          "pattern": "regex-pattern",
          "style": "style-name"
        }
      ],
      "styles": {
        "style-name": "color-name"
      }
    }
  }
}
```

### Key Elements

1. **`language-name`**: The identifier for the language (e.g., "go", "python", "javascript")
2. **`extensions`**: File extensions associated with this language (e.g., ["go"], ["py"], ["js", "jsx"])
3. **`rules`**: Array of highlighting rules, each containing:
   - **`name`**: Descriptive name for the rule (e.g., "keywords", "strings", "comments")
   - **`pattern`**: Regular expression pattern to match code elements
   - **`style`**: Reference to a style defined in the styles section
4. **`styles`**: Map of style names to color names or ANSI color codes

## Adding a New Language

To add a new language, follow these steps:

1. **Identify the language elements** to highlight:
   - Keywords and reserved words
   - Strings and string literals
   - Comments (single-line and multi-line)
   - Numbers and numeric literals
   - Functions and method declarations
   - Special syntax elements specific to the language

2. **Create regular expressions** for each element type:
   - Use word boundaries (`\b`) for keywords
   - Create patterns that match the language syntax precisely
   - Test patterns against sample code

3. **Define styles** for each element type:
   - Choose appropriate colors for readability
   - Consider color contrast and accessibility

4. **Add the language definition** to the configuration file

## Regex Pattern Writing Tips

Regular expressions are powerful but can be tricky. Here are some tips:

### Common Patterns

- **Keywords**: `\\b(keyword1|keyword2|keyword3)\\b`
- **Strings**: `"[^"]*"` (double-quoted) or `'[^']*'` (single-quoted)
- **Numbers**: `\\b\\d+(\\.\\d+)?\\b` (integers and decimals)
- **Comments**: `//.*` (single-line) or `/\\*[\\s\\S]*?\\*/` (multi-line)
- **Function declarations**: `\\bfunction\\s+([A-Za-z0-9_]+)\\s*\\(`

### Important Considerations

- **Escape backslashes** in JSON: Write `\\` instead of `\` in patterns
- **Order matters**: Rules are applied in the order they appear
- **Overlapping matches**: Be careful of patterns that might overlap
- **Performance**: Complex patterns can slow down highlighting for large files

## Available Styles

The `highlight` tool supports these ANSI color and style codes:

### Text Styles
- `bold` - Bold text
- `italic` - Italic text
- `underline` - Underlined text
- `reset` - Reset all styles

### Colors
- `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`
- `brightblack`, `brightred`, `brightgreen`, `brightyellow`, `brightblue`, `brightmagenta`, `brightcyan`, `brightwhite`

## Testing Your Configuration

After adding a new language:

1. **Create a sample file** with the target language
2. **Run highlight**: `highlight --config your-config.json sample-file.ext`
3. **Verify highlighting** works as expected
4. **Iterate and refine** your patterns as needed

## Optimizing Regex Patterns

For better performance:

1. **Prefer specific patterns** over generic ones
2. **Use word boundaries** (`\b`) where appropriate
3. **Avoid excessive capturing groups**
4. **Be careful with backtracking** patterns like `.*`
5. **Test with large files** to ensure performance

## Examples

### Example 1: Ruby Language Configuration

```json
"ruby": {
  "extensions": ["rb"],
  "rules": [
    {
      "name": "keywords",
      "pattern": "\\b(begin|class|def|do|else|elsif|end|ensure|for|if|module|rescue|return|self|super|then|unless|until|when|while|yield)\\b",
      "style": "keyword"
    },
    {
      "name": "symbols",
      "pattern": ":[a-zA-Z_][a-zA-Z0-9_]*",
      "style": "symbol"
    },
    {
      "name": "strings",
      "pattern": "\"[^\"]*\"|'[^']*'",
      "style": "string"
    },
    {
      "name": "comments",
      "pattern": "#.*",
      "style": "comment"
    },
    {
      "name": "class_names",
      "pattern": "\\b[A-Z][a-zA-Z0-9_]*\\b",
      "style": "class"
    }
  ],
  "styles": {
    "keyword": "cyan",
    "symbol": "brightred",
    "string": "green",
    "comment": "brightblack",
    "class": "brightmagenta"
  }
}
```

### Example 2: Adding CSS Language Support

```json
"css": {
  "extensions": ["css"],
  "rules": [
    {
      "name": "selectors",
      "pattern": "[.#]?[a-zA-Z0-9_-]+\\s*[,{]",
      "style": "selector"
    },
    {
      "name": "properties",
      "pattern": "[a-zA-Z-]+\\s*:",
      "style": "property"
    },
    {
      "name": "values",
      "pattern": ":\\s*[^;]+;",
      "style": "value"
    },
    {
      "name": "comments",
      "pattern": "/\\*[\\s\\S]*?\\*/",
      "style": "comment"
    },
    {
      "name": "units",
      "pattern": "\\d+(\\.\\d+)?(px|em|rem|vh|vw|%|s|ms)",
      "style": "unit"
    }
  ],
  "styles": {
    "selector": "brightcyan",
    "property": "yellow",
    "value": "green",
    "comment": "brightblack",
    "unit": "brightmagenta"
  }
}

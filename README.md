# highlight - Lightweight Syntax Highlighter

A highly optimized, minimal, and scalable command-line syntax highlighting tool written in Go. This program is designed as a more powerful alternative to the `cat` command with syntax highlighting capabilities.

## Features

- **Highly efficient:** Uses a two-goroutine pipeline architecture for maximum performance
- **Extensible:** Supports custom syntax highlighting rules via external JSON configuration
- **Minimal dependencies:** Only uses Go standard library packages
- **Automatic language detection:** Based on file extensions
- **Line ending support:** Handles both LF and CRLF line endings
- **Support for stdin:** Can be used in command pipelines
- **Standard `cat` compatibility:** Supports common cat flags like `-n`, `-b`, `-s`, and `-E`
- **Multi-language support:** Includes built-in support for Go, Python, JavaScript, JSON, Markdown, XML/HTML, and SQL
- **Security-focused:** Uses low-level syscall operations for file I/O
- **Integrated paging:** Use `--less` flag to view large files with the `less` pager

## Installation

```bash
go install github.com/AmirMahdyJebreily/hili-cat/cmd@latest
```

## Usage

```bash
# Highlight a file with automatic language detection
highlight main.go

# Highlight multiple files
highlight file1.go file2.go

# Highlight stdin with a specific language
cat main.go | highlight --lang go

# Use a custom configuration file
highlight --config /path/to/config.json main.go

# Specify line ending format
highlight --line-ending crlf file.go

# Show line numbers (like cat -n)
highlight -n file.go

# Show only non-blank line numbers (like cat -b)
highlight -b file.go

# Squeeze repeated blank lines (like cat -s)
highlight -s file.go

# Show line endings with $ marker (like cat -E)
highlight -E file.go

# View highlighted file using less for pagination
highlight --less file.go
```

## Configuration

The default configuration is stored at `/etc/highlight/config.json`. A sample configuration is provided in the `config` directory.

You can create your own configuration file with custom syntax highlighting rules for different languages. The configuration file uses JSON format with the following structure:

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

Available styles include:
- Text styles: `bold`, `italic`, `underline`
- Colors: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`
- Bright colors: `brightblack`, `brightred`, `brightgreen`, `brightyellow`, `brightblue`, `brightmagenta`, `brightcyan`, `brightwhite`

## Architecture

The `highlight` tool uses a highly efficient two-goroutine pipeline design:

1. **Reader Goroutine**: Handles file I/O or stdin using low-level syscalls for maximum performance
2. **Highlighter Goroutine**: Processes the data and applies syntax highlighting using regex-based rules

This pipeline approach allows for efficient streaming of data, even with large files, minimizing memory usage while maintaining high performance.

### Data Flow:

```
┌─────────┐     ┌─────────────┐     ┌───────────┐
│ Input   │────▶│ Buffered    │────▶│ Output    │
│ Source  │     │ Channel     │     │ (Stdout)  │
└─────────┘     └─────────────┘     └───────────┘
    │                                    ▲
    │                                    │
    ▼                                    │
┌─────────┐                        ┌─────────────┐
│ Reader  │                        │ Highlighter │
│ Goroutine│                        │ Goroutine  │
└─────────┘                        └─────────────┘
```

## Options

- `--config`: Path to the configuration file (default: `/etc/highlight/config.json`)
- `--lang`: Language for syntax highlighting (required when reading from stdin)
- `--line-ending`: Line ending to use (`auto`, `lf`, `crlf`, default: `auto`)
- `-n, --number`: Number all output lines
- `-b, --number-nonblank`: Number non-blank output lines
- `-s, --squeeze-blank`: Suppress repeated empty output lines
- `-E, --show-ends`: Display $ at end of each line
- `--less, --pager`: Pipe output to `less -R` command for paged viewing
- `--help`: Show help message

## Performance Considerations

`highlight` is designed with performance as a primary goal:

- **Buffered I/O**: Uses efficient buffer sizes for optimal read/write performance
- **Regexp Optimization**: Precompiles regex patterns to minimize CPU usage
- **Memory Management**: Minimizes allocations to reduce GC overhead
- **Syscall Usage**: Direct syscall usage for file operations instead of higher-level abstractions
- **Channel Buffering**: Properly sized channels to prevent blocking in the pipeline

## Examples

### Basic Usage

```bash
# Highlight a Go file with automatic language detection
highlight main.go
```

### Piping from stdin

```bash
# Pipe the output of git diff into highlight for colored diff output
git diff | highlight --lang diff
```

### Use in scripts

```bash
# Use highlight in a script to show syntax-highlighted output
#!/bin/bash
echo "Displaying highlighted code:"
highlight --lang go main.go
```

## Running the Demo

The project includes a demo script that showcases the key features of the `highlight` tool:

```bash
# Make the demo script executable
chmod +x demo.sh

# Run the demo
./demo.sh
```

The demo will:
1. Build the highlight tool if needed
2. Show syntax highlighting for Go and JSON files
3. Demonstrate line numbering
4. Show how to use highlight with piped input
5. Display line endings
6. Demonstrate highlighting multiple files

## Contributing

Contributions are welcome! Here are some areas where help is needed:

1. Adding support for more programming languages
2. Performance optimizations
3. Additional formatting options
4. Bug fixes and test improvements

For detailed information on how to contribute, please read the [Contributing Guide](CONTRIBUTING.md).

## Language Configuration

`highlight` is designed to be easily extended with support for additional programming languages.
To learn how to add support for your favorite language, see the [Language Configuration Guide](CONFIG_GUIDE.md).

## Release Process

The project uses GitHub Actions to automate builds:
- Tag `pre-release0.0.1` and any version with `v*` prefix triggers automatic builds
- Linux builds are automatically created and attached to GitHub releases
- See the workflow file in `.github/workflows/linux-build.yml` for details

## License

See the [LICENSE](LICENSE) file for details.

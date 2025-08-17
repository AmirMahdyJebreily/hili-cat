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
```

## Configuration

The default configuration is stored at `/etc/highlight/config.json`. A sample configuration is provided in the `config` directory.

You can create your own configuration file with custom syntax highlighting rules for different languages.

## Options

- `--config`: Path to the configuration file (default: `/etc/highlight/config.json`)
- `--lang`: Language for syntax highlighting (required when reading from stdin)
- `--line-ending`: Line ending to use (`auto`, `lf`, `crlf`, default: `auto`)
- `-n, --number`: Number all output lines
- `-b, --number-nonblank`: Number non-blank output lines
- `-s, --squeeze-blank`: Suppress repeated empty output lines
- `-E, --show-ends`: Display $ at end of each line
- `--help`: Show help message

## License

See the [LICENSE](LICENSE) file for details.

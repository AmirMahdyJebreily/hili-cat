package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/AmirMahdyJebreily/hili-cat/internal/config"
	"github.com/AmirMahdyJebreily/hili-cat/internal/highlighter"
	fileio "github.com/AmirMahdyJebreily/hili-cat/internal/io" // Renamed to avoid conflict with standard io
)

// Default buffer sizes for performance optimization
const (
	defaultBufferSize = 4096
	channelBufferSize = 1000
)

// printUsage prints the program's usage information
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: hili-cat [options] [file...]\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  hili-cat file.go                    # Highlight a Go file\n")
	fmt.Fprintf(os.Stderr, "  cat file.json | hili-cat --lang json # Highlight JSON from stdin\n")
	fmt.Fprintf(os.Stderr, "  hili-cat --config /path/to/config.json file.py # Use custom config\n")
	fmt.Fprintf(os.Stderr, "  hili-cat --less large_file.go       # View highlighted file with pagination\n")
	fmt.Fprintf(os.Stderr, "\nNote: hili-cat is designed for Linux systems only and requires the 'less' command for pagination.\n")
}

func main() {
	// Parse command-line flags
	configPath := flag.String("config", config.DefaultConfigPath, "Path to the configuration file")
	lang := flag.String("lang", "", "Language for syntax highlighting (required when reading from stdin)")
	lineEnding := flag.String("line-ending", "auto", "Line ending to use (auto, lf, crlf)")
	numberLines := flag.Bool("n", false, "Number all output lines")
	numberNonBlank := flag.Bool("b", false, "Number non-blank output lines")
	squeezeBlank := flag.Bool("s", false, "Suppress repeated empty output lines")
	showEnds := flag.Bool("E", false, "Display $ at end of each line")
	useLess := flag.Bool("less", false, "Pipe output to 'less -R' command for paged viewing")
	help := flag.Bool("help", false, "Show help message")

	// Add long-form flags
	flag.BoolVar(numberLines, "number", *numberLines, "Number all output lines")
	flag.BoolVar(numberNonBlank, "number-nonblank", *numberNonBlank, "Number non-blank output lines")
	flag.BoolVar(squeezeBlank, "squeeze-blank", *squeezeBlank, "Suppress repeated empty output lines")
	flag.BoolVar(showEnds, "show-ends", *showEnds, "Display $ at end of each line")
	flag.BoolVar(useLess, "pager", *useLess, "Pipe output to 'less -R' command for paged viewing")

	flag.Parse()

	if *help {
		printUsage()
		return
	}

	// Ensure config file exists (create default if not)
	if err := config.EnsureExists(*configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	// Load the configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Initialize the reader
	reader := fileio.NewReader(defaultBufferSize)

	// Get the file arguments
	args := flag.Args()

	// Create highlighter options
	opts := highlighter.Options{
		NumberLines:    *numberLines,
		NumberNonBlank: *numberNonBlank,
		SqueezeBlank:   *squeezeBlank,
		ShowEnds:       *showEnds,
	}

	// Determine if we're reading from stdin or files
	if len(args) == 0 {
		processStdin(reader, cfg, *lang, *lineEnding, opts, *useLess)
	} else {
		processFiles(reader, cfg, args, *lang, *lineEnding, opts, *useLess)
	}
}

// processStdin handles input from standard input
func processStdin(reader *fileio.Reader, cfg config.Config, lang, lineEnding string, opts highlighter.Options, useLess bool) {
	// Reading from stdin
	if lang == "" {
		fmt.Fprintln(os.Stderr, "Error: --lang is required when reading from stdin")
		printUsage()
		os.Exit(1)
	}

	// Detect or set line ending
	var detectedLineEnding string
	if lineEnding == "auto" {
		// Try to detect from stdin
		bufReader := bufio.NewReader(os.Stdin)
		buf := make([]byte, 1024)
		n, _ := bufReader.Read(buf)
		if n > 0 {
			detectedLineEnding = fileio.DetectLineEnding(buf[:n])
		} else {
			detectedLineEnding = highlighter.LF // Default to LF if no data
		}

		// Need to reset stdin
		if n > 0 {
			// This is a simplified approach; in a real application,
			// we would need a more robust way to peek at stdin without consuming it
			fmt.Fprintf(os.Stderr, "Warning: Line ending detection consumed some input data\n")
			detectedLineEnding = highlighter.LF // Default to LF as fallback
		}
	} else if lineEnding == "crlf" {
		detectedLineEnding = highlighter.CRLF
	} else {
		detectedLineEnding = highlighter.LF
	}

	// Create highlighter
	h, err := highlighter.NewHighlighter(convertConfig(cfg), lang, detectedLineEnding, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Set up the pipeline
	dataCh := make(chan []byte, channelBufferSize)
	var wg sync.WaitGroup

	// Launch reader goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader.ProcessFile("", dataCh, nil)
		close(dataCh) // Important: close channel when reader is done
	}()

	// Launch processor goroutine
	wg.Add(1)
	go processOutput(dataCh, h, &wg, useLess)

	// Wait for both goroutines to complete
	wg.Wait()
}

// processFiles handles input from multiple files
func processFiles(reader *fileio.Reader, cfg config.Config, files []string, langOverride, lineEnding string, opts highlighter.Options, useLess bool) {
	for _, filePath := range files {
		// Determine language from file extension if not explicitly provided
		fileLang := langOverride
		if fileLang == "" {
			fileLang = config.DetectLanguage(cfg, filePath)
			if fileLang == "" {
				fmt.Fprintf(os.Stderr, "Error: Could not determine language for %s. Use --lang flag.\n", filePath)
				continue
			}
		}

		// Determine line ending
		var detectedLineEnding string
		if lineEnding == "auto" {
			// Open file to detect line ending
			file, err := fileio.OpenFile(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}

			buf := make([]byte, 1024)
			n, _ := file.Read(buf)
			file.Close()

			if n > 0 {
				detectedLineEnding = fileio.DetectLineEnding(buf[:n])
			} else {
				detectedLineEnding = highlighter.LF // Default to LF if file is empty
			}
		} else if lineEnding == "crlf" {
			detectedLineEnding = highlighter.CRLF
		} else {
			detectedLineEnding = highlighter.LF
		}

		// Create highlighter
		h, err := highlighter.NewHighlighter(convertConfig(cfg), fileLang, detectedLineEnding, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		// Set up the pipeline
		dataCh := make(chan []byte, channelBufferSize)
		var wg sync.WaitGroup

		// Launch reader goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			reader.ProcessFile(filePath, dataCh, nil)
			close(dataCh) // Important: close channel when reader is done
		}()

		// Launch processor goroutine
		wg.Add(1)
		go processOutput(dataCh, h, &wg, useLess)

		// Wait for both goroutines to complete
		wg.Wait()
	}
}

// processOutput handles the highlighting and output of data
func processOutput(dataCh <-chan []byte, h *highlighter.Highlighter, wg *sync.WaitGroup, useLess bool) {
	defer wg.Done()

	if useLess {
		// Start less command with -R flag to interpret ANSI color codes
		cmd := exec.Command("less", "-R")

		// Get pipes for stdin and stdout
		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to start less: %v\n", err)
			// Fallback to regular output if less fails
			for data := range dataCh {
				fmt.Print(h.ProcessContent(data))
			}
			return
		}

		// Set output to terminal
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Start the command
		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to start less: %v\n", err)
			// Fallback to regular output if less fails
			for data := range dataCh {
				fmt.Print(h.ProcessContent(data))
			}
			return
		}

		// Process data and write to less
		for data := range dataCh {
			processed := h.ProcessContent(data)
			_, err := io.WriteString(stdin, processed)
			if err != nil {
				break
			}
		}

		// Close stdin to signal EOF to less
		stdin.Close()

		// Wait for less to exit
		cmd.Wait()
	} else {
		// Regular output to stdout
		for data := range dataCh {
			fmt.Print(h.ProcessContent(data))
		}
	}
}

// Note: waitForCompletion function has been removed as it's no longer needed
// The channel closing is now handled directly in the reader goroutines

// convertConfig converts config.Config to highlighter.Config
func convertConfig(cfg config.Config) highlighter.Config {
	languages := make(map[string]highlighter.Language)

	for lang, language := range cfg.Languages {
		rules := make([]highlighter.HighlightRule, len(language.Rules))
		for i, rule := range language.Rules {
			rules[i] = highlighter.HighlightRule{
				Name:    rule.Name,
				Pattern: rule.Pattern,
				Style:   rule.Style,
			}
		}

		languages[lang] = highlighter.Language{
			Extensions: language.Extensions,
			Rules:      rules,
			Styles:     language.Styles,
		}
	}

	return highlighter.Config{
		Languages: languages,
	}
}

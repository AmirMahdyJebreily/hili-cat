#!/bin/bash

# Demo script for highlight tool
# This script demonstrates the capabilities of the highlight tool

# Check if highlight is installed
if [ ! -f "./build/highlight-linux-amd64" ] && [ ! -f "./build/highlight" ]; then
  echo "Building highlight..."
  ./build.sh
fi

HIGHLIGHT="./build/highlight"
if [ -f "./build/highlight-linux-amd64" ]; then
  HIGHLIGHT="./build/highlight-linux-amd64"
fi

# Create examples directory if it doesn't exist
mkdir -p examples

# Display header
echo -e "\n\033[1m=== highlight Tool Demo ===\033[0m\n"

# Demo 1: Basic Go file highlighting
echo -e "\033[1m1. Highlighting a Go file:\033[0m"
$HIGHLIGHT examples/demo.go
echo

# Demo 2: JSON file highlighting
echo -e "\033[1m2. Highlighting a JSON file:\033[0m"
$HIGHLIGHT examples/demo.json
echo

# Demo 3: Using line numbers
echo -e "\033[1m3. Go file with line numbers:\033[0m"
$HIGHLIGHT -n examples/demo.go
echo

# Demo 4: Piping from another command
echo -e "\033[1m4. Highlighting from stdin (piped input):\033[0m"
cat examples/demo.go | $HIGHLIGHT --lang go
echo

# Demo 5: Show line endings
echo -e "\033[1m5. Displaying line endings:\033[0m"
$HIGHLIGHT -E examples/demo.go | head -n 5
echo

# Demo 6: Multiple files at once
echo -e "\033[1m6. Highlighting multiple files:\033[0m"
$HIGHLIGHT examples/demo.go examples/demo.json
echo

echo -e "\033[1m=== Demo Complete ===\033[0m"

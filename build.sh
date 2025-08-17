#!/bin/bash

# Build script for highlight tool

# Function to display usage information
show_usage() {
  echo "Usage: $0 [options]"
  echo "Options:"
  echo "  --all     Build for all platforms (Linux, macOS, Windows)"
  echo "  --linux   Build for Linux"
  echo "  --darwin  Build for macOS"
  echo "  --windows Build for Windows"
  echo "  --arm     Build for ARM (Linux)"
  echo "  --help    Show this help message"
}

# Create output directory if it doesn't exist
mkdir -p build

# Default build for current platform
if [ $# -eq 0 ]; then
  echo "Building highlight for current platform..."
  go build -ldflags="-s -w" -trimpath -o build/highlight ./cmd/highlight
  echo "Binary saved to build/highlight"
  exit 0
fi

# Parse command-line arguments
while [ "$1" != "" ]; do
  case $1 in
    --all )
      echo "Building highlight for all platforms..."
      GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o build/highlight-linux-amd64 ./cmd/highlight
      GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o build/highlight-darwin-amd64 ./cmd/highlight
      GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o build/highlight-windows-amd64.exe ./cmd/highlight
      GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -trimpath -o build/highlight-linux-arm ./cmd/highlight
      echo "All binaries saved to build/ directory"
      ;;
    --linux )
      echo "Building highlight for Linux..."
      GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o build/highlight-linux-amd64 ./cmd/highlight
      echo "Binary saved to build/highlight-linux-amd64"
      ;;
    --darwin )
      echo "Building highlight for macOS..."
      GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o build/highlight-darwin-amd64 ./cmd/highlight
      echo "Binary saved to build/highlight-darwin-amd64"
      ;;
    --windows )
      echo "Building highlight for Windows..."
      GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o build/highlight-windows-amd64.exe ./cmd/highlight
      echo "Binary saved to build/highlight-windows-amd64.exe"
      ;;
    --arm )
      echo "Building highlight for ARM Linux..."
      GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -trimpath -o build/highlight-linux-arm ./cmd/highlight
      echo "Binary saved to build/highlight-linux-arm"
      ;;
    --help )
      show_usage
      exit 0
      ;;
    * )
      echo "Unknown option: $1"
      show_usage
      exit 1
      ;;
  esac
  shift
done

echo "Build complete!"

#!/bin/bash

# Build script for highlight tool

# Build options
BUILD_FLAGS="-ldflags=\"-s -w\" -trimpath"
OPTIMIZATION="-gcflags=\"all=-N -l -c=2\" -asmflags=\"all=-trimpath=/go\""

# Function to display usage information
show_usage() {
  echo "Usage: $0 [options]"
  echo "Options:"
  echo "  --all       Build for all platforms (Linux, macOS, Windows)"
  echo "  --linux     Build for Linux"
  echo "  --darwin    Build for macOS"
  echo "  --windows   Build for Windows"
  echo "  --arm       Build for ARM (Linux)"
  echo "  --debug     Build with debug symbols (no optimization)"
  echo "  --release   Build with maximum optimization (default)"
  echo "  --upx       Compress binary with UPX (if available)"
  echo "  --help      Show this help message"
}

# Create output directory if it doesn't exist
mkdir -p build

# Build parameters
DEBUG=0
COMPRESS=0

# Function to build for a specific platform
build_binary() {
  local os=$1
  local arch=$2
  local arm=$3
  local output=$4
  local flags="-ldflags=\"-s -w\""
  
  # If debug mode is enabled, use debug flags
  if [ $DEBUG -eq 1 ]; then
    flags=""
  fi
  
  # Build command
  BUILD_CMD="GOOS=$os GOARCH=$arch"
  
  # Add GOARM if specified
  if [ -n "$arm" ]; then
    BUILD_CMD="$BUILD_CMD GOARM=$arm"
  fi
  
  # Execute build
  echo "Building highlight for $os/$arch..."
  eval "$BUILD_CMD go build -trimpath $flags -o $output ./cmd/highlight"
  
  # Compress with UPX if requested
  if [ $COMPRESS -eq 1 ]; then
    if command -v upx &> /dev/null; then
      echo "Compressing $output with UPX..."
      upx -9 "$output"
    else
      echo "Warning: UPX not found, skipping compression"
    fi
  fi
  
  echo "Binary saved to $output"
}

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
      build_binary "linux" "amd64" "" "build/highlight-linux-amd64"
      build_binary "darwin" "amd64" "" "build/highlight-darwin-amd64"
      build_binary "windows" "amd64" "" "build/highlight-windows-amd64.exe"
      build_binary "linux" "arm" "7" "build/highlight-linux-arm"
      build_binary "linux" "arm64" "" "build/highlight-linux-arm64"
      echo "All binaries saved to build/ directory"
      ;;
    --linux )
      build_binary "linux" "amd64" "" "build/highlight-linux-amd64"
      ;;
    --darwin )
      build_binary "darwin" "amd64" "" "build/highlight-darwin-amd64"
      ;;
    --windows )
      build_binary "windows" "amd64" "" "build/highlight-windows-amd64.exe"
      ;;
    --arm )
      build_binary "linux" "arm" "7" "build/highlight-linux-arm"
      ;;
    --arm64 )
      build_binary "linux" "arm64" "" "build/highlight-linux-arm64"
      ;;
    --debug )
      DEBUG=1
      echo "Debug mode enabled"
      ;;
    --release )
      DEBUG=0
      echo "Release mode enabled (optimized build)"
      ;;
    --upx )
      COMPRESS=1
      echo "UPX compression enabled"
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

# Make the binary executable
chmod +x build/highlight*

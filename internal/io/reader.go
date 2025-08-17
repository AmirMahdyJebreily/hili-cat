package io

import (
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"
)

const defaultBufferSize = 4096

// Reader provides low-level I/O operations using syscalls
type Reader struct {
	bufSize int
}

// NewReader creates a new Reader instance
func NewReader(bufSize int) *Reader {
	if bufSize <= 0 {
		bufSize = defaultBufferSize
	}
	return &Reader{
		bufSize: bufSize,
	}
}

// OpenFile opens a file using syscall for enhanced security and performance
func OpenFile(path string) (*os.File, error) {
	fd, err := syscall.Open(path, syscall.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", path, err)
	}

	return os.NewFile(uintptr(fd), path), nil
}

// ProcessFile reads from a file or stdin and sends data to the channel
func (r *Reader) ProcessFile(filePath string, dataCh chan<- []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	var reader io.Reader

	if filePath == "" {
		// Reading from stdin
		reader = os.Stdin
	} else {
		// Reading from file using syscall for security
		file, err := OpenFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		defer file.Close()
		reader = file
	}

	// Use buffered reader for performance
	buf := make([]byte, r.bufSize)

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			// Create a copy of the buffer to avoid race conditions
			data := make([]byte, n)
			copy(data, buf[:n])
			dataCh <- data
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading: %v\n", err)
			break
		}
	}
}

// DetectLineEnding examines a byte slice to determine the line ending
func DetectLineEnding(data []byte) string {
	const (
		LF   = "\n"
		CRLF = "\r\n"
	)

	if len(data) == 0 {
		return LF // Default to LF if no data
	}

	for i := 0; i < len(data)-1; i++ {
		if data[i] == '\r' && data[i+1] == '\n' {
			return CRLF
		}
	}
	return LF
}

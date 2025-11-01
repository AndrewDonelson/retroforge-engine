package modulestate

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileSystemReader reads files from the filesystem
type FileSystemReader struct {
	basePath string
}

// NewFileSystemReader creates a file reader that reads from the filesystem
func NewFileSystemReader(basePath string) *FileSystemReader {
	return &FileSystemReader{
		basePath: basePath,
	}
}

// ReadFile reads a file from the filesystem
func (fsr *FileSystemReader) ReadFile(path string) ([]byte, error) {
	var fullPath string
	if fsr.basePath != "" {
		fullPath = filepath.Join(fsr.basePath, path)
	} else {
		fullPath = path
	}
	return os.ReadFile(fullPath)
}

// MapFileReader reads files from an in-memory map
type MapFileReader struct {
	files map[string][]byte
}

// NewMapFileReader creates a file reader that reads from a map
func NewMapFileReader(files map[string][]byte) *MapFileReader {
	return &MapFileReader{
		files: files,
	}
}

// ReadFile reads a file from the map
func (mfr *MapFileReader) ReadFile(path string) ([]byte, error) {
	// Try exact path first
	if content, ok := mfr.files[path]; ok {
		return content, nil
	}

	// Try with assets/ prefix
	assetsPath := filepath.Join("assets", path)
	if content, ok := mfr.files[assetsPath]; ok {
		return content, nil
	}

	return nil, fmt.Errorf("file not found: %s", path)
}

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// OpenFile opens a file at the specified path and returns a pointer to the file.
func OpenFile(filePath string) (*os.File, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	return file, nil
}

func ReadFile(filePath string) ([]byte, error) {
	file, err := OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return data, nil
}

func ParseJsonToMap(data []byte) (map[string]any, error) {
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

// WriteToFileFromReader writes the content from an io.Reader to a file in the specified directory.
func WriteToFileFromReader(dir, fileName string, reader io.Reader) (string, error) {
	if fileName == "" {
		return "", fmt.Errorf("file name cannot be empty")
	}

	if err := CreateDir(dir); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	filePath := filepath.Join(dir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}
	return filePath, file.Sync()
}

// DeleteFiles removes all files in the specified directory.
func DeleteFiles(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, file := range files {
		if err := os.RemoveAll(filepath.Join(dir, file.Name())); err != nil {
			return fmt.Errorf("failed to delete file %s: %w", file.Name(), err)
		}
	}
	return nil
}

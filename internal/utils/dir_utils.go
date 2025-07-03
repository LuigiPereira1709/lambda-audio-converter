package utils

import (
	"log"
	"os"
	"strings"
)

func GetWorkDir() string {
	workDir := os.Getenv("WORK_DIR")
	if workDir == "" {
		panic("WORK_DIR is not set")
	}

	if stats, err := os.Stat(workDir); err != nil || !stats.IsDir() {
		log.Printf("WORK_DIR doesn't exist, trying to create: %s", workDir)
		if err := CreateDir(workDir); err != nil {
			log.Printf("Failed to create work directory: %v", err)
			panic("Failed to create work directory: " + err.Error())
		}
	}

	return workDir
}

// CreateDir creates a directory if it does not exist, including all necessary parent directories.
func CreateDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// GetParentDir returns the parent directory of the given dir.
func GetParentDir(dir string) string {
	idx := strings.LastIndex(dir, "/")
	if idx == -1 {
		return ""
	}

	parentDir := dir[:idx]
	return parentDir
}

package converter

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pitanguinha.com/audio-converter/internal/utils"
)

var RequiredMetadataKeys = []string{"title", "year"}

// FFmpegCommand represents a command to be executed by FFmpeg.
type FFmpegCommand struct {
	GlobalOptions []string
	Inputs        []string
	Filter        []string
	Map           []string
	Codec         []string
	Metadata      []string
	Flags         []string
	Output        string
}

const (
	processedFileName = "processed_file"
)

// NewFFmpegCommand creates a new FFmpegCommand with default values.
func NewFFmpegCommand(inputsPaths, metadataMap map[string]string, requiredKeys []string) (*FFmpegCommand, error) {
	ffmpegBinPath := os.Getenv("FFMPEG_BIN_PATH")
	audioCodec := os.Getenv("AUDIO_CODEC")
	audioFormat := os.Getenv("AUDIO_FORMAT")

	if inputsPaths == nil || metadataMap == nil {
		return nil, errors.New("filesPaths and metadataMap cannot be nil")
	}

	absPaths, err := parseToAbsPaths(inputsPaths)
	if err != nil {
		return nil, fmt.Errorf("error parsing file paths: %v", err)
	}

	requiredKeys = append(requiredKeys, RequiredMetadataKeys...)
	if err := validateRequiredMetadataKeys(requiredKeys, metadataMap); err != nil {
		return nil, fmt.Errorf("missing required metadata: %v", err)
	}
	metadataArr := []string{
		"-metadata:s:v", "title=Album cover",
		"-metadata:s:v", "comment=Cover (front)",
		"-metadata", "title=" + metadataMap["title"],
		"-metadata", "year=" + metadataMap["year"],
	}

	outputPath := filepath.Join(utils.GetWorkDir(), processedFileName+"."+audioFormat)

	return &FFmpegCommand{
		GlobalOptions: []string{ffmpegBinPath, "-y", "-progress", "pipe:1", "-nostats"},
		Inputs:        []string{"-i", absPaths["content"], "-i", absPaths["thumbnail"]},
		Filter:        []string{"-vf", "scale=trunc(iw/2)*2:trunc(ih/2)*2"},
		Map:           []string{"-map", "0:a", "-map", "1:v"},
		Codec:         []string{"-c:a", audioCodec},
		Metadata:      metadataArr,
		Flags:         []string{"-movflags", "faststart"},
		Output:        outputPath,
	}, nil
}

// GetOutputFilePath returns the output file path of the FFmpeg command.
func (c *FFmpegCommand) GetOutputFilePath() string {
	return c.Output
}

// AddMetadataFromMap adds metadata to the FFmpeg command from a map of keys and values.
func (c *FFmpegCommand) AddMetadataFromMap(keys []string, metadataMap map[string]string) {
	for _, key := range keys {
		if value, ok := metadataMap[key]; ok && value != "" {
			c.Metadata = append(c.Metadata, "-metadata", fmt.Sprintf("%s=%s", key, value))
		}
	}
}

// BuildCommand constructs the FFmpeg command as a slice of strings.
func (c *FFmpegCommand) BuildCommand() []string {
	command := append(c.GlobalOptions, c.Inputs...)
	command = append(command, c.Filter...)
	command = append(command, c.Map...)
	command = append(command, c.Codec...)
	command = append(command, c.Metadata...)
	command = append(command, c.Flags...)
	command = append(command, c.Output)
	return command
}

// requiredMetadataKeys checks if all required metadata keys are present and non-empty.
func validateRequiredMetadataKeys(requiredKeys []string, metadataMap map[string]string) error {
	for _, key := range requiredKeys {
		if value, exists := metadataMap[key]; !exists || value == "" {
			return fmt.Errorf("required metadata key %s is missing or empty", key)
		}
	}
	return nil
}

// parseToAbsPath converts relative file paths to absolute paths.
func parseToAbsPaths(inputPaths map[string]string) (map[string]string, error) {
	absPaths := make(map[string]string)

	for key, path := range inputPaths {
		if path == "" {
			return nil, fmt.Errorf("path for %s is empty", key)
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path for %s: %v", key, err)
		}
		absPaths[key] = absPath
	}

	return absPaths, nil
}

// NormalizeFilename replaces invalid characters in a filename with underscores.
func NormalizeFilename(filename string) string {
	normalized := strings.ReplaceAll(filename, "/", "_")
	normalized = strings.ReplaceAll(normalized, "\\", "_")
	normalized = strings.ReplaceAll(normalized, ":", "_")
	normalized = strings.ReplaceAll(normalized, "*", "_")
	normalized = strings.ReplaceAll(normalized, "?", "_")
	normalized = strings.ReplaceAll(normalized, "\"", "_")
	normalized = strings.ReplaceAll(normalized, "<", "_")
	normalized = strings.ReplaceAll(normalized, ">", "_")
	normalized = strings.ReplaceAll(normalized, "|", "_")
	return normalized
}

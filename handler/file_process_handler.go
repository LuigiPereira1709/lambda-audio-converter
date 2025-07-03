package handler

import (
	"fmt"
	"log"

	"pitanguinha.com/audio-converter/internal/converter"
	"pitanguinha.com/audio-converter/internal/converter/music"
	"pitanguinha.com/audio-converter/internal/converter/podcast"
)

// ProcessAudioFile processes the files based on their type (music or podcast) and executes the FFmpeg command.
// Returns the details of the conversion process and any error encountered during the process.
func ProcessAudioFile(duration float64, filesPaths, metadataMap map[string]string) (*converter.FFmpegProgressDetails, error) {
	cmd, err := buildFFmpegCommand(filesPaths, metadataMap)
	if err != nil {
		return nil, fmt.Errorf("error building ffmpeg command: %w", err)
	}

	log.Println("FFmpeg command:", cmd)

	return converter.FFmpegExecutor(cmd, duration)
}

// buildFFmpegCommand constructs the FFmpeg command based on the type of media (music or podcast).
func buildFFmpegCommand(filesPaths, metadataMap map[string]string) ([]string, error) {
	var cmd []string
	var err error
	switch metadataMap["type"] {
	case "music":
		cmd, err = music.BuildCommand(filesPaths, metadataMap)
	case "podcast":
		cmd, err = podcast.BuildCommand(filesPaths, metadataMap)
	}

	if err != nil {
		return nil, err
	}

	return cmd, nil
}

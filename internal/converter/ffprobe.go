package converter

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"pitanguinha.com/audio-converter/internal/utils"
)

// GetDurationFromFile retrieves the duration of a media file using ffprobe.
func GetDurationFromFile(filePath string) (float64, error) {
	FFprobeBinPath := os.Getenv("FFPROBE_BIN_PATH")
	command := []string{FFprobeBinPath, "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := utils.ExecCommand(ctx, command...)
	if cmd == nil {
		return 0, fmt.Errorf("failed to create command for ffprobe")
	}

	output, err := utils.GetCommandOutput(cmd)
	if err != nil {
		return 0, err
	}

	output = output[:len(output)-1] // Remove the trailing newline character

	duration, err := strconv.ParseFloat(output, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return duration, nil
}

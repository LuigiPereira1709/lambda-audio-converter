package converter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"pitanguinha.com/audio-converter/internal/utils"
)

// FFmpegProgressDetails holds the details of the FFmpeg command execution progress.
type FFmpegProgressDetails struct {
	Duration          float64
	CurrentTime       float64
	CurrentLine       string
	Progress          float64
	TimeElapsed       string
	Finished          bool
	ProcessedFilePath string
}

const (
	ctxTimeOut  = 6 * time.Minute // 5 minutes timeout for FFmpeg command execution
	keyOutTime  = "out_time"
	keyProgress = "progress"
)

// FFmpegExecutor executes an FFmpeg command and tracks its progress.
func FFmpegExecutor(command []string, duration float64) (*FFmpegProgressDetails, error) {
	startTime := time.Now()
	details := newFFmpegProgressDetails(duration)

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	cmd := utils.ExecCommand(ctx, command...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return details, fmt.Errorf("error getting stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return details, fmt.Errorf("error starting ffmpeg command: %w", err)
	}

	utils.ScanStd(stdout, ffmpegProgressHandler(details))

	if err := cmd.Wait(); err != nil {
		return details, fmt.Errorf("error waiting for ffmpeg command: %w", err)
	}

	details.TimeElapsed = utils.FormatDuration(time.Since(startTime))
	details.ProcessedFilePath = command[len(command)-1] // INFO: The output file path always is the last argument in the command
	return details, nil
}

// newFFmpegProgressDetails creates a new instance of FFmpegProgressDetails with the specified duration.
func newFFmpegProgressDetails(duration float64) *FFmpegProgressDetails {
	return &FFmpegProgressDetails{
		Duration:    duration,
		CurrentTime: 0,
		CurrentLine: "",
		Progress:    0,
		TimeElapsed: "",
		Finished:    false,
	}
}

// String returns a string representation of the FFmpegProgressDetails.
func (t *FFmpegProgressDetails) String() string {
	return fmt.Sprintf("Progress: %.2f%%. Current Time: %s. Duration: %.2fs. Finished: %t. Elapsed Time: %s. Current Line: %s",
		t.Progress, utils.FormatSecondsToTime(t.CurrentTime), t.Duration, t.Finished, t.TimeElapsed, t.CurrentLine)
}

// ffmpegProgressHandler processes each line of output from the FFmpeg command to update the progress details.
func ffmpegProgressHandler(t *FFmpegProgressDetails) func(line string) {
	return func(line string) {
		if t.Finished {
			return
		}

		t.CurrentLine = line

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return
		}

		key, value := parts[0], parts[1]

		switch key {
		case keyOutTime:
			if t.Duration > 0 {
				seconds := utils.ParserTimeToSeconds(value)
				if seconds > t.CurrentTime {
					t.CurrentTime = seconds
				}

				progress := (t.CurrentTime / t.Duration) * 100

				if progress > t.Progress {
					t.Progress = progress
				}
			}
		case keyProgress:
			if value == "end" {
				t.Finished = true
			}
		}
	}
}

package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseTimeToSeconds converts a time string in the format "HH:MM:SS" to seconds.
func ParserTimeToSeconds(timeStr string) float64 {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0
	}
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	s, _ := strconv.ParseFloat(parts[2], 64)
	return float64(h*3600+m*60) + s
}

// FormatSecondsToTime converts a duration in seconds to a string in the format "HH:MM:SS".
func FormatSecondsToTime(seconds float64) string {
	h := int(seconds) / 3600
	m := (int(seconds) % 3600) / 60
	s := int(seconds) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// formatDuration formats a time.Duration into a string in the format "HH:MM:SS".
func FormatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

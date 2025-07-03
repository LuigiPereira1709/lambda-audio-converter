package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
)

// ExecCommand executes a command in the given context and returns the command, its stdout and stderr pipes, and any error encountered during execution.
func ExecCommand(ctx context.Context, command ...string) *exec.Cmd {
	if len(command) == 0 {
		return nil
	}
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Dir = GetWorkDir() // Set the working directory for the command

	return cmd
}

// Read lines from the std(stdout or stderr) of a command execution an call the callback function for each line.
func ScanStd(std io.ReadCloser, callback func(line string)) error {
	scanner := bufio.NewScanner(std)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			callback(line)
		}
	}

	readErr := scanner.Err()
	closeErr := std.Close()
	if readErr != nil {
		return fmt.Errorf("error reading std: %w", readErr)
	}
	if closeErr != nil {
		return fmt.Errorf("error closing std: %w", closeErr)
	}

	return nil
}

// GetCommandOutput executes a command and returns its combined stdout and stderr output as a string.
func GetCommandOutput(cmd *exec.Cmd) (string, error) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w", err)
	}
	return string(output), nil
}

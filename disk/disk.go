package cxsysinfodisk

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// GetDirectorySize returns the size of the directory at the given path in bytes.
func GetDirectorySize(path string) (uint64, error) {
	// Execute the du command with -sx flag
	cmd := exec.Command("du", "-sx", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}

	// Extract the size from the command output
	output := strings.Fields(out.String())
	if len(output) > 0 {
		// Convert the size from kilobytes to bytes
		sizeKB, err := strconv.ParseUint(output[0], 10, 64)
		if err != nil {
			return 0, err
		}
		sizeBytes := sizeKB * 1024
		return sizeBytes, nil
	}
	return 0, fmt.Errorf("failed to parse du command output")
}

package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// MoveToTrash moves a file or directory to a custom trash directory.
// If a file with the same name already exists in the trash, a timestamp is appended.
// Returns the destination path where the file was moved to.
func MoveToTrash(sourcePath string, trashPath string) (string, error) {
	// Get absolute path for the source
	absSourcePath, err := filepath.Abs(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Ensure trash directory exists
	if err := os.MkdirAll(trashPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create trash directory: %w", err)
	}

	// Get the base name of the file/directory
	baseName := filepath.Base(absSourcePath)
	destPath := filepath.Join(trashPath, baseName)

	// If a file with the same name already exists in trash, append timestamp
	if _, err := os.Stat(destPath); err == nil {
		ext := filepath.Ext(baseName)
		nameWithoutExt := baseName[:len(baseName)-len(ext)]
		timestamp := time.Now().Format("20060102-150405")
		destPath = filepath.Join(trashPath, fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext))
	}

	// Move the file to trash using SafeMove to support cross-filesystem moves
	if err := SafeMove(absSourcePath, destPath); err != nil {
		return "", fmt.Errorf("failed to move to trash: %w", err)
	}

	return destPath, nil
}

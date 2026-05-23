package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

func (rs sceneRoutes) Upload(w http.ResponseWriter, r *http.Request) {
	// Limit upload size to 10GB
	mgr := manager.GetInstance()
	// Use configured max upload size, default to 10GB
	maxUploadSize := mgr.Config.GetMaxUploadSize()
	if maxUploadSize <= 0 {
		maxUploadSize = 10 << 30
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		logger.Errorf("error parsing upload form: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		logger.Errorf("error getting uploaded file: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	cfg := mgr.Config

	// Determine destination path
	var destPath string
	stashPaths := cfg.GetStashPaths()
	if len(stashPaths) > 0 {
		destPath = filepath.Join(stashPaths[0].Path, filepath.Base(header.Filename))
	} else if genPath := cfg.GetGeneratedPath(); genPath != "" {
		destPath = filepath.Join(genPath, filepath.Base(header.Filename))
	} else {
		http.Error(w, "no stash paths or generated path configured", http.StatusBadRequest)
		return
	}

	destPath = resolveUniquePath(destPath)

	// Write the file to disk
	dst, err := os.Create(destPath)
	if err != nil {
		logger.Errorf("error creating file %s: %v", destPath, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	written, err := io.Copy(dst, file)
	if err != nil {
		dst.Close()
		logger.Errorf("error writing file %s: %v", destPath, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dst.Close()

	logger.Infof("uploaded file %s (%d bytes)", destPath, written)

	// Trigger scan for the uploaded file
	input := manager.ScanMetadataInput{
		Paths: []string{destPath},
	}

	jobID, err := mgr.Scan(r.Context(), input)
	if err != nil {
		logger.Errorf("error starting scan for %s: %v", destPath, err)
		// File was saved; return success with scan=failed so frontend can handle
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"path":    destPath,
			"size":    written,
			"scan":    "failed",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"job_id":  jobID,
		"path":    destPath,
		"size":    written,
	})
}

// resolveUniquePath returns the next available filename by appending
// a numeric suffix if the file already exists.
func resolveUniquePath(path string) string {
	if exists, _ := fsutil.FileExists(path); !exists {
		return path
	}

	ext := filepath.Ext(path)
	base := path[:len(path)-len(ext)]

	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s_%d%s", base, i, ext)
		if exists, _ := fsutil.FileExists(candidate); !exists {
			return candidate
		}
	}
}

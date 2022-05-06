package astikit

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MoveFile is a cancellable move of a local file to a local or remote location
func MoveFile(ctx context.Context, dst, src string, f CopyFileFunc) (err error) {
	// Copy
	if err = CopyFile(ctx, dst, src, f); err != nil {
		err = fmt.Errorf("astikit: copying file %s to %s failed: %w", src, dst, err)
		return
	}

	// Delete
	if err = os.Remove(src); err != nil {
		err = fmt.Errorf("astikit: removing %s failed: %w", src, err)
		return
	}
	return
}

// CopyFileFunc represents a CopyFile func
type CopyFileFunc func(ctx context.Context, dst string, srcStat os.FileInfo, srcFile *os.File) error

// CopyFile is a cancellable copy of a local file to a local or remote location
func CopyFile(ctx context.Context, dst, src string, f CopyFileFunc) (err error) {
	// Check context
	if err = ctx.Err(); err != nil {
		return
	}

	// Stat src
	var srcStat os.FileInfo
	if srcStat, err = os.Stat(src); err != nil {
		err = fmt.Errorf("astikit: stating %s failed: %w", src, err)
		return
	}

	// Src is a dir
	if srcStat.IsDir() {
		// Walk through the dir
		if err = filepath.Walk(src, func(path string, info os.FileInfo, errWalk error) (err error) {
			// Check error
			if errWalk != nil {
				err = errWalk
				return
			}

			// Do not process root
			if src == path {
				return
			}

			// Copy
			p := filepath.Join(dst, strings.TrimPrefix(path, filepath.Clean(src)))
			if err = CopyFile(ctx, p, path, f); err != nil {
				err = fmt.Errorf("astikit: copying %s to %s failed: %w", path, p, err)
				return
			}
			return nil
		}); err != nil {
			err = fmt.Errorf("astikit: walking through %s failed: %w", src, err)
			return
		}
		return
	}

	// Open src
	var srcFile *os.File
	if srcFile, err = os.Open(src); err != nil {
		err = fmt.Errorf("astikit: opening %s failed: %w", src, err)
		return
	}
	defer srcFile.Close()

	// Custom
	if err = f(ctx, dst, srcStat, srcFile); err != nil {
		err = fmt.Errorf("astikit: custom failed: %w", err)
		return
	}
	return
}

// LocalCopyFileFunc is the local CopyFileFunc that allows doing cross partition copies
func LocalCopyFileFunc(ctx context.Context, dst string, srcStat os.FileInfo, srcFile *os.File) (err error) {
	// Check context
	if err = ctx.Err(); err != nil {
		return
	}

	// Create the destination folder
	if err = os.MkdirAll(filepath.Dir(dst), DefaultDirMode); err != nil {
		err = fmt.Errorf("astikit: mkdirall %s failed: %w", filepath.Dir(dst), err)
		return
	}

	// Create the destination file
	var dstFile *os.File
	if dstFile, err = os.Create(dst); err != nil {
		err = fmt.Errorf("astikit: creating %s failed: %w", dst, err)
		return
	}
	defer dstFile.Close()

	// Chmod using os.chmod instead of file.Chmod
	if err = os.Chmod(dst, srcStat.Mode()); err != nil {
		err = fmt.Errorf("astikit: chmod %s %s failed, %w", dst, srcStat.Mode(), err)
		return
	}

	// Copy the content
	if _, err = Copy(ctx, dstFile, srcFile); err != nil {
		err = fmt.Errorf("astikit: copying content of %s to %s failed: %w", srcFile.Name(), dstFile.Name(), err)
		return
	}
	return
}

// SignalHandler represents a func that can handle a signal
type SignalHandler func(s os.Signal)

// TermSignalHandler returns a SignalHandler that is executed only on a term signal
func TermSignalHandler(f func()) SignalHandler {
	return func(s os.Signal) {
		if isTermSignal(s) {
			f()
		}
	}
}

// LoggerSignalHandler returns a SignalHandler that logs the signal
func LoggerSignalHandler(l SeverityLogger, ignoredSignals ...os.Signal) SignalHandler {
	ss := make(map[os.Signal]bool)
	for _, s := range ignoredSignals {
		ss[s] = true
	}
	return func(s os.Signal) {
		if _, ok := ss[s]; ok {
			return
		}
		l.Debugf("astikit: received signal %s", s)
	}
}

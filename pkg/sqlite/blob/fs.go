package blob

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

const (
	blobsDirDepth  int = 2
	blobsDirLength int = 2 // thumbDirDepth * thumbDirLength must be smaller than the length of checksum
)

type FSReader interface {
	Open(name string) (fs.ReadDirFile, error)
}

type FSWriter interface {
	Create(name string) (*os.File, error)
	MkdirAll(path string, perm fs.FileMode) error

	Remove(name string) error

	file.RenamerRemover
}

type FS interface {
	FSReader
	FSWriter
}

type FilesystemReader struct {
	path string
	fs   FSReader
}

func (s *FilesystemReader) checksumToPath(checksum string) string {
	return filepath.Join(s.path, fsutil.GetIntraDir(checksum, blobsDirDepth, blobsDirLength), checksum)
}

func (s *FilesystemReader) Read(ctx context.Context, checksum string) ([]byte, error) {
	if s.path == "" {
		return nil, fmt.Errorf("no path set")
	}

	fn := s.checksumToPath(checksum)
	f, err := s.fs.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", fn, err)
	}

	defer f.Close()

	return io.ReadAll(f)
}

type FilesystemStore struct {
	FilesystemReader
	deleter *file.Deleter
}

func NewFilesystemStore(path string, fs FS) *FilesystemStore {
	deleter := &file.Deleter{
		RenamerRemover: fs,
	}

	return &FilesystemStore{
		FilesystemReader: *NewReadonlyFilesystemStore(path, fs),
		deleter:          deleter,
	}
}

func NewReadonlyFilesystemStore(path string, fs FSReader) *FilesystemReader {
	return &FilesystemReader{
		path: path,
		fs:   fs,
	}
}

func (s *FilesystemStore) Write(ctx context.Context, checksum string, data []byte) error {
	fs, ok := s.fs.(FS)
	if !ok {
		return fmt.Errorf("internal error: fs is not an FS")
	}

	if s.path == "" {
		return fmt.Errorf("no path set")
	}

	fn := s.checksumToPath(checksum)

	// create the directory if it doesn't exist
	if err := fs.MkdirAll(filepath.Dir(fn), 0755); err != nil {
		return fmt.Errorf("creating directory %q: %w", filepath.Dir(fn), err)
	}

	logger.Debugf("Writing blob file %s", fn)
	out, err := fs.Create(fn)
	if err != nil {
		return fmt.Errorf("creating file %q: %w", fn, err)
	}

	r := bytes.NewReader(data)

	if _, err = io.Copy(out, r); err != nil {
		return fmt.Errorf("writing file %q: %w", fn, err)
	}

	return nil
}

func (s *FilesystemStore) Delete(ctx context.Context, checksum string) error {
	if s.path == "" {
		return fmt.Errorf("no path set")
	}

	s.deleter.RegisterHooks(ctx)

	fn := s.checksumToPath(checksum)

	if err := s.deleter.Files([]string{fn}); err != nil {
		return fmt.Errorf("deleting file %q: %w", fn, err)
	}

	return nil
}

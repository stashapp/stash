package astikit

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// internal shouldn't lead with a "/"
func zipInternalPath(p string) (external, internal string) {
	if items := strings.Split(p, ".zip"); len(items) > 1 {
		external = items[0] + ".zip"
		internal = strings.TrimPrefix(strings.Join(items[1:], ".zip"), string(os.PathSeparator))
		return
	}
	external = p
	return
}

// Zip zips a src into a dst
// Possible dst formats are:
//   - /path/to/zip.zip
//   - /path/to/zip.zip/root/path
func Zip(ctx context.Context, dst, src string) (err error) {
	// Get external/internal path
	externalPath, internalPath := zipInternalPath(dst)

	// Make sure the directory exists
	if err = os.MkdirAll(filepath.Dir(externalPath), DefaultDirMode); err != nil {
		return fmt.Errorf("astikit: mkdirall %s failed: %w", filepath.Dir(externalPath), err)
	}

	// Create destination file
	var dstFile *os.File
	if dstFile, err = os.Create(externalPath); err != nil {
		return fmt.Errorf("astikit: creating %s failed: %w", externalPath, err)
	}
	defer dstFile.Close()

	// Create zip writer
	var zw = zip.NewWriter(dstFile)
	defer zw.Close()

	// Walk
	if err = filepath.Walk(src, func(path string, info os.FileInfo, e error) (err error) {
		// Process error
		if e != nil {
			err = e
			return
		}

		// Init header
		var h *zip.FileHeader
		if h, err = zip.FileInfoHeader(info); err != nil {
			return fmt.Errorf("astikit: initializing zip header failed: %w", err)
		}

		// Set header info
		h.Name = filepath.Join(internalPath, strings.TrimPrefix(path, src))
		if info.IsDir() {
			h.Name += string(os.PathSeparator)
		} else {
			h.Method = zip.Deflate
		}

		// Create writer
		var w io.Writer
		if w, err = zw.CreateHeader(h); err != nil {
			return fmt.Errorf("astikit: creating zip header failed: %w", err)
		}

		// If path is dir, stop here
		if info.IsDir() {
			return
		}

		// Open path
		var walkFile *os.File
		if walkFile, err = os.Open(path); err != nil {
			return fmt.Errorf("astikit: opening %s failed: %w", path, err)
		}
		defer walkFile.Close()

		// Copy
		if _, err = Copy(ctx, w, walkFile); err != nil {
			return fmt.Errorf("astikit: copying failed: %w", err)
		}
		return
	}); err != nil {
		return fmt.Errorf("astikit: walking failed: %w", err)
	}
	return
}

// Unzip unzips a src into a dst
// Possible src formats are:
//   - /path/to/zip.zip
//   - /path/to/zip.zip/root/path
func Unzip(ctx context.Context, dst, src string) (err error) {
	// Get external/internal path
	externalPath, internalPath := zipInternalPath(src)

	// Make sure the destination exists
	if err = os.MkdirAll(dst, DefaultDirMode); err != nil {
		return fmt.Errorf("astikit: mkdirall %s failed: %w", dst, err)
	}

	// Open overall reader
	var r *zip.ReadCloser
	if r, err = zip.OpenReader(externalPath); err != nil {
		return fmt.Errorf("astikit: opening overall zip reader on %s failed: %w", externalPath, err)
	}
	defer r.Close()

	// Loop through files to determine their type
	var dirs, files, symlinks = make(map[string]*zip.File), make(map[string]*zip.File), make(map[string]*zip.File)
	for _, f := range r.File {
		// Validate internal path
		if internalPath != "" && !strings.HasPrefix(f.Name, internalPath) {
			continue
		}
		var p = filepath.Join(dst, strings.TrimPrefix(f.Name, internalPath))

		// Check file type
		if f.FileInfo().Mode()&os.ModeSymlink != 0 {
			symlinks[p] = f
		} else if f.FileInfo().IsDir() {
			dirs[p] = f
		} else {
			files[p] = f
		}
	}

	// Invalid internal path
	if internalPath != "" && len(dirs) == 0 && len(files) == 0 && len(symlinks) == 0 {
		return fmt.Errorf("astikit: content in archive does not match specified internal path %s", internalPath)
	}

	// Create dirs
	for p, f := range dirs {
		if err = os.MkdirAll(p, f.FileInfo().Mode().Perm()); err != nil {
			return fmt.Errorf("astikit: mkdirall %s failed: %w", p, err)
		}
	}

	// Create files
	for p, f := range files {
		if err = createZipFile(ctx, f, p); err != nil {
			return fmt.Errorf("astikit: creating zip file into %s failed: %w", p, err)
		}
	}

	// Create symlinks
	for p, f := range symlinks {
		if err = createZipSymlink(f, p); err != nil {
			return fmt.Errorf("astikit: creating zip symlink into %s failed: %w", p, err)
		}
	}
	return
}

func createZipFile(ctx context.Context, f *zip.File, p string) (err error) {
	// Open file reader
	var fr io.ReadCloser
	if fr, err = f.Open(); err != nil {
		return fmt.Errorf("astikit: opening zip reader on file %s failed: %w", f.Name, err)
	}
	defer fr.Close()

	// Since dirs don't always come up we make sure the directory of the file exists with default
	// file mode
	if err = os.MkdirAll(filepath.Dir(p), DefaultDirMode); err != nil {
		return fmt.Errorf("astikit: mkdirall %s failed: %w", filepath.Dir(p), err)
	}

	// Open the file
	var fl *os.File
	if fl, err = os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode().Perm()); err != nil {
		return fmt.Errorf("astikit: opening file %s failed: %w", p, err)
	}
	defer fl.Close()

	// Copy
	if _, err = Copy(ctx, fl, fr); err != nil {
		return fmt.Errorf("astikit: copying %s into %s failed: %w", f.Name, p, err)
	}
	return
}

func createZipSymlink(f *zip.File, p string) (err error) {
	// Open file reader
	var fr io.ReadCloser
	if fr, err = f.Open(); err != nil {
		return fmt.Errorf("astikit: opening zip reader on file %s failed: %w", f.Name, err)
	}
	defer fr.Close()

	// If file is a symlink we retrieve the target path that is in the content of the file
	var b []byte
	if b, err = ioutil.ReadAll(fr); err != nil {
		return fmt.Errorf("astikit: ioutil.Readall on %s failed: %w", f.Name, err)
	}

	// Create the symlink
	if err = os.Symlink(string(b), p); err != nil {
		return fmt.Errorf("astikit: creating symlink from %s to %s failed: %w", string(b), p, err)
	}
	return
}

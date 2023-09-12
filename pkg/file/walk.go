package file

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/stashapp/stash/pkg/models"
)

// Modified from github.com/facebookgo/symwalk

// BSD License

// For symwalk software

// Copyright (c) 2015, Facebook, Inc. All rights reserved.

// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:

//  * Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.

//  * Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.

//  * Neither the name Facebook nor the names of its contributors may be used to
//    endorse or promote products derived from this software without specific
//    prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
// ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// symwalkFunc calls the provided WalkFn for regular files.
// However, when it encounters a symbolic link, it resolves the link fully using the
// filepath.EvalSymlinks function and recursively calls symwalk.Walk on the resolved path.
// This ensures that unlink filepath.Walk, traversal does not stop at symbolic links.
//
// Note that symwalk.Walk does not terminate if there are any non-terminating loops in
// the file structure.
func walkSym(f models.FS, filename string, linkDirname string, walkFn fs.WalkDirFunc) error {
	symWalkFunc := func(path string, info fs.DirEntry, err error) error {

		if fname, err := filepath.Rel(filename, path); err == nil {
			path = filepath.Join(linkDirname, fname)
		} else {
			return err
		}

		if err == nil && info.Type()&os.ModeSymlink == os.ModeSymlink {
			finalPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				// don't bail out if symlink is invalid
				return walkFn(path, info, err)
			}
			info, err := f.Lstat(finalPath)
			if err != nil {
				return walkFn(path, &statDirEntry{
					info: info,
				}, err)
			}
			if info.IsDir() {
				return walkSym(f, finalPath, path, walkFn)
			}
		}

		return walkFn(path, info, err)
	}
	return fsWalk(f, filename, symWalkFunc)
}

// symWalk extends filepath.Walk to also follow symlinks
func symWalk(fs models.FS, path string, walkFn fs.WalkDirFunc) error {
	return walkSym(fs, path, path, walkFn)
}

type statDirEntry struct {
	info fs.FileInfo
}

func (d *statDirEntry) Name() string               { return d.info.Name() }
func (d *statDirEntry) IsDir() bool                { return d.info.IsDir() }
func (d *statDirEntry) Type() fs.FileMode          { return d.info.Mode().Type() }
func (d *statDirEntry) Info() (fs.FileInfo, error) { return d.info, nil }

func fsWalk(f models.FS, root string, fn fs.WalkDirFunc) error {
	info, err := f.Lstat(root)
	if err != nil {
		err = fn(root, nil, err)
	} else {
		err = walkDir(f, root, &statDirEntry{info}, fn)
	}
	if errors.Is(err, fs.SkipDir) {
		return nil
	}
	return err
}

func walkDir(f models.FS, path string, d fs.DirEntry, walkDirFn fs.WalkDirFunc) error {
	if err := walkDirFn(path, d, nil); err != nil || !d.IsDir() {
		if errors.Is(err, fs.SkipDir) && d.IsDir() {
			// Successfully skipped directory.
			err = nil
		}
		return err
	}

	dirs, err := readDir(f, path)
	if err != nil {
		// Second call, to report ReadDir error.
		err = walkDirFn(path, d, err)
		if err != nil {
			return err
		}
	}

	for _, d1 := range dirs {
		name := d1.Name()
		// Prevent infinite loops; this can happen with certain FS implementations (e.g. ZipFS).
		if name == "" || name == "." {
			continue
		}
		path1 := filepath.Join(path, name)
		if err := walkDir(f, path1, d1, walkDirFn); err != nil {
			if errors.Is(err, fs.SkipDir) {
				break
			}
			return err
		}
	}
	return nil
}

// readDir reads the directory named by dirname and returns
// a sorted list of directory entries.
func readDir(fs models.FS, dirname string) ([]fs.DirEntry, error) {
	f, err := fs.Open(dirname)
	if err != nil {
		return nil, err
	}
	dirs, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	return dirs, nil
}

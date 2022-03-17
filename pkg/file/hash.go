package file

import (
	"io"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/hash/oshash"
)

type FSHasher struct{}

func (h *FSHasher) OSHash(src io.ReadSeeker, size int64) (string, error) {
	return oshash.FromReader(src, size)
}

func (h *FSHasher) MD5(src io.Reader) (string, error) {
	return md5.FromReader(src)
}

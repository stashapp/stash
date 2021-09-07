package file

import (
	"io"

	"github.com/stashapp/stash/pkg/utils"
)

type FSHasher struct{}

func (h *FSHasher) OSHash(src io.ReadSeeker, size int64) (string, error) {
	return utils.OSHashFromReader(src, size)
}

func (h *FSHasher) MD5(src io.Reader) (string, error) {
	return utils.MD5FromReader(src)
}

package file

import "github.com/stashapp/stash/pkg/utils"

type FSHasher struct{}

func (h *FSHasher) OSHash(path string) (string, error) {
	return utils.OSHashFromFilePath(path)
}

func (h *FSHasher) MD5(path string) (string, error) {
	return utils.MD5FromFilePath(path)
}

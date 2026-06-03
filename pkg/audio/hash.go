package audio

import (
	"github.com/stashapp/stash/pkg/models"
)

// GetHash returns the hash of the file, based on the hash algorithm provided. If
// hash algorithm is MD5, then Checksum is returned. Otherwise, then panic.
func GetHash(f models.File, hashAlgorithm models.HashAlgorithm) string {
	switch hashAlgorithm {
	case models.HashAlgorithmMd5:
		return f.Base().Fingerprints.GetString(models.FingerprintTypeMD5)
	default:
		panic("unknown hash algorithm")
	}
}

package scene

import (
	"github.com/stashapp/stash/pkg/models"
)

// GetHash returns the hash of the file, based on the hash algorithm provided. If
// hash algorithm is MD5, then Checksum is returned. Otherwise, OSHash is returned.
func GetHash(f models.File, hashAlgorithm models.HashAlgorithm) string {
	switch hashAlgorithm {
	case models.HashAlgorithmMd5:
		return f.Base().Fingerprints.GetString(models.FingerprintTypeMD5)
	case models.HashAlgorithmOshash:
		return f.Base().Fingerprints.GetString(models.FingerprintTypeOshash)
	default:
		panic("unknown hash algorithm")
	}
}

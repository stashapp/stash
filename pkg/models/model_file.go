package models

import "time"

type File struct {
	Checksum    string    `db:"checksum" json:"checksum"`
	OSHash      string    `db:"oshash" json:"oshash"`
	Path        string    `db:"path" json:"path"`
	Size        string    `db:"size" json:"size"`
	FileModTime time.Time `db:"file_mod_time" json:"file_mod_time"`
}

// GetHash returns the hash of the scene, based on the hash algorithm provided. If
// hash algorithm is MD5, then Checksum is returned. Otherwise, OSHash is returned.
func (s File) GetHash(hashAlgorithm HashAlgorithm) string {
	switch hashAlgorithm {
	case HashAlgorithmMd5:
		return s.Checksum
	case HashAlgorithmOshash:
		return s.OSHash
	default:
		panic("unknown hash algorithm")
	}
}

func (s File) Equal(o File) bool {
	return s.Path == o.Path && s.Checksum == o.Checksum && s.OSHash == o.OSHash && s.Size == o.Size && s.FileModTime.Equal(o.FileModTime)
}

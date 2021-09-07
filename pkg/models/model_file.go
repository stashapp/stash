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
	var ret string
	if hashAlgorithm == HashAlgorithmMd5 {
		ret = s.Checksum
	} else if hashAlgorithm == HashAlgorithmOshash {
		ret = s.OSHash
	} else {
		panic("unknown hash algorithm")
	}

	return ret
}

func (s File) Equal(o File) bool {
	return s.Path == o.Path && s.Checksum == o.Checksum && s.OSHash == o.OSHash && s.Size == o.Size && s.FileModTime == o.FileModTime
}

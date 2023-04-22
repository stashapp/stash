package sqlite

import "github.com/corona10/goimagehash"

func phashDistanceFn(phash1 int64, phash2 int64) (int64, error) {
	hash1 := goimagehash.NewImageHash(uint64(phash1), goimagehash.PHash)
	hash2 := goimagehash.NewImageHash(uint64(phash2), goimagehash.PHash)
	distance, _ := hash1.Distance(hash2)
	return int64(distance), nil
}

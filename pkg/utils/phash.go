package utils

import (
	"math"
	"strconv"

	"github.com/corona10/goimagehash"
	"github.com/stashapp/stash/pkg/sliceutil"
)

type Phash struct {
	ID        int    `db:"id"`
	Hash      int64  `db:"phash"`
	Duration  float64 `db:"duration"`
	Neighbors []int
	Bucket    int
}

func FindDuplicates(hashes []*Phash, distance int, durationDiff float64) [][]int {
	for i, subject := range hashes {
		subjectHash := goimagehash.NewImageHash(uint64(subject.Hash), goimagehash.PHash)
		for j, neighbor := range hashes {
			if i != j && subject.ID != neighbor.ID {
				neighbourDurationDistance := 0.
				if subject.Duration > 0 && neighbor.Duration > 0 {
					neighbourDurationDistance = math.Abs(subject.Duration - neighbor.Duration)
				}
				if (neighbourDurationDistance <= durationDiff) || (durationDiff < 0) {
					neighborHash := goimagehash.NewImageHash(uint64(neighbor.Hash), goimagehash.PHash)
					neighborDistance, _ := subjectHash.Distance(neighborHash)
					if neighborDistance <= distance {
						subject.Neighbors = append(subject.Neighbors, j)
					}
				}
			}
		}
	}

	var buckets [][]int
	for _, subject := range hashes {
		if len(subject.Neighbors) > 0 && subject.Bucket == -1 {
			bucket := len(buckets)
			ids := []int{subject.ID}
			subject.Bucket = bucket
			findNeighbors(bucket, subject.Neighbors, hashes, &ids)

			if len(ids) > 1 {
				buckets = append(buckets, ids)
			}
		}
	}

	return buckets
}

func findNeighbors(bucket int, neighbors []int, hashes []*Phash, ids *[]int) {
	for _, id := range neighbors {
		hash := hashes[id]
		if hash.Bucket == -1 {
			hash.Bucket = bucket
			*ids = sliceutil.AppendUnique(*ids, hash.ID)
			findNeighbors(bucket, hash.Neighbors, hashes, ids)
		}
	}
}

func PhashToString(phash int64) string {
	return strconv.FormatUint(uint64(phash), 16)
}

func StringToPhash(s string) (int64, error) {
	ret, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0, err
	}

	return int64(ret), nil
}

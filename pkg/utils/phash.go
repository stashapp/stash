package utils

import (
	"math"
	"math/bits"
	"runtime"
	"strconv"
	"sync"

	"github.com/stashapp/stash/pkg/sliceutil"
)

type Phash struct {
	ID        int     `db:"id"`
	Hash      int64   `db:"phash"`
	Duration  float64 `db:"duration"`
	Neighbors []int
	Bucket    int
}

func FindDuplicates(hashes []*Phash, distance int, durationDiff float64) [][]int {
	// Pre-calculate hash values to avoid allocations and method calls in the inner loop
	uintHashes := make([]uint64, len(hashes))
	for i, h := range hashes {
		uintHashes[i] = uint64(h.Hash)
	}

	numHashes := len(hashes)
	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Distribute work among workers
	for w := 0; w < numWorkers; w++ {
		go func(workerID int) {
			defer wg.Done()
			for i := workerID; i < numHashes; i += numWorkers {
				subject := hashes[i]
				subjectHash := uintHashes[i]

				for j := 0; j < numHashes; j++ {
					if i == j {
						continue
					}
					neighbor := hashes[j]
					if subject.ID == neighbor.ID {
						continue
					}

					// Check duration if applicable (for scenes)
					if durationDiff >= 0 {
						if subject.Duration > 0 && neighbor.Duration > 0 {
							if math.Abs(subject.Duration-neighbor.Duration) > durationDiff {
								continue
							}
						}
					}

					neighborHash := uintHashes[j]
					// Hamming distance using native bit counting
					if bits.OnesCount64(subjectHash^neighborHash) <= distance {
						subject.Neighbors = append(subject.Neighbors, j)
					}
				}
			}
		}(w)
	}

	wg.Wait()

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

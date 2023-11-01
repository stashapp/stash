package utils

import (
	"math"
	"strconv"

	"github.com/corona10/goimagehash"
	"github.com/stashapp/stash/pkg/sliceutil"
)

type Phash struct {
	SceneID   int     `db:"id"`
	Hash      int64   `db:"phash"`
	Duration  float64 `db:"duration"`
	Neighbors []int
	Bucket    int
}

func FindDuplicates(hashes []*Phash, distance int, durationDiff float64) [][]int {
	for i, scene := range hashes {
		sceneHash := goimagehash.NewImageHash(uint64(scene.Hash), goimagehash.PHash)
		for j, neighbor := range hashes {
			if i != j && scene.SceneID != neighbor.SceneID {
				neighbourDurationDistance := 0.
				if scene.Duration > 0 && neighbor.Duration > 0 {
					neighbourDurationDistance = math.Abs(scene.Duration - neighbor.Duration)
				}
				if (neighbourDurationDistance <= durationDiff) || (durationDiff < 0) {
					neighborHash := goimagehash.NewImageHash(uint64(neighbor.Hash), goimagehash.PHash)
					neighborDistance, _ := sceneHash.Distance(neighborHash)
					if neighborDistance <= distance {
						scene.Neighbors = append(scene.Neighbors, j)
					}
				}
			}
		}
	}

	var buckets [][]int
	for _, scene := range hashes {
		if len(scene.Neighbors) > 0 && scene.Bucket == -1 {
			bucket := len(buckets)
			scenes := []int{scene.SceneID}
			scene.Bucket = bucket
			findNeighbors(bucket, scene.Neighbors, hashes, &scenes)

			if len(scenes) > 1 {
				buckets = append(buckets, scenes)
			}
		}
	}

	return buckets
}

func findNeighbors(bucket int, neighbors []int, hashes []*Phash, scenes *[]int) {
	for _, id := range neighbors {
		hash := hashes[id]
		if hash.Bucket == -1 {
			hash.Bucket = bucket
			*scenes = sliceutil.AppendUnique(*scenes, hash.SceneID)
			findNeighbors(bucket, hash.Neighbors, hashes, scenes)
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

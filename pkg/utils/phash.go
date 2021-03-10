package utils

import (
	"github.com/corona10/goimagehash"
)

type Phash struct {
	SceneID   int   `db:"id"`
	Hash      int64 `db:"phash"`
	Neighbors []int
	Bucket    int
}

func FindDuplicates(hashes []*Phash, distance int) [][]int {
	for i, scene := range hashes {
		sceneHash := goimagehash.NewImageHash(uint64(scene.Hash), goimagehash.PHash)
		for j, neighbor := range hashes {
			if i != j {
				neighborHash := goimagehash.NewImageHash(uint64(neighbor.Hash), goimagehash.PHash)
				neighborDistance, _ := sceneHash.Distance(neighborHash)
				if neighborDistance <= distance {
					scene.Neighbors = append(scene.Neighbors, j)
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
			buckets = append(buckets, scenes)
		}
	}

	return buckets
}

func findNeighbors(bucket int, neighbors []int, hashes []*Phash, scenes *[]int) {
	for _, id := range neighbors {
		hash := hashes[id]
		if hash.Bucket == -1 {
			hash.Bucket = bucket
			*scenes = append(*scenes, hash.SceneID)
			findNeighbors(bucket, hash.Neighbors, hashes, scenes)
		}
	}
}

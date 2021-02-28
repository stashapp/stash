package main

import (
	"fmt"
	"github.com/corona10/goimagehash"
	"image/jpeg"
	"os"
)

func main() {
	file1, _ := os.Open("sample1.jpg")
	file2, _ := os.Open("sample2.jpg")
	defer file1.Close()
	defer file2.Close()

	img1, _ := jpeg.Decode(file1)
	img2, _ := jpeg.Decode(file2)
	hash1, _ := goimagehash.AverageHash(img1)
	hash2, _ := goimagehash.AverageHash(img2)
	distance, _ := hash1.Distance(hash2)
	fmt.Printf("Distance between images: %v\n", distance)

	hash1, _ = goimagehash.DifferenceHash(img1)
	hash2, _ = goimagehash.DifferenceHash(img2)
	distance, _ = hash1.Distance(hash2)
	fmt.Printf("Distance between images: %v\n", distance)

	hash1, _ = goimagehash.PerceptionHash(img1)
	hash2, _ = goimagehash.PerceptionHash(img2)
	distance, _ = hash1.Distance(hash2)
	fmt.Printf("Distance between images: %v\n", distance)
	fmt.Println(hash1.ToString())
	fmt.Println(hash2.ToString())
	fmt.Println(hash1.Bits())
	fmt.Println(hash2.Bits())
}

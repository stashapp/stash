package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image/jpeg"
	"os"

	"github.com/corona10/goimagehash"
)

func main() {
	file1, _ := os.Open("sample1.jpg")
	file2, _ := os.Open("sample2.jpg")
	defer file1.Close()
	defer file2.Close()
	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	img1, _ := jpeg.Decode(file1)
	img2, _ := jpeg.Decode(file2)
	width, height := 16, 16
	hash1, _ := goimagehash.ExtPerceptionHash(img1, width, height)
	hash2, _ := goimagehash.ExtPerceptionHash(img2, width, height)
	hash1024, _ := goimagehash.ExtAverageHash(img2, 32, 32)
	distance, _ := hash1.Distance(hash2)
	fmt.Printf("Distance between images: %v\n", distance)
	err := hash1.Dump(foo)
	if err != nil {
		fmt.Println(err)
	}
	foo.Flush()
	bar := bufio.NewReader(&b)
	hash3, err := goimagehash.LoadExtImageHash(bar)
	if err != nil {
		fmt.Println(err)
	}
	distance, err = hash1.Distance(hash1024)
	fmt.Println(err)
	distance, _ = hash1.Distance(hash3)
	fmt.Printf("Distance between hash1 and hash3: %v\n", distance)
	distance, _ = hash2.Distance(hash3)
	fmt.Printf("Distance between hash2 and hash3: %v\n", distance)
	fmt.Println(hash1.ToString())
	fmt.Println(hash2.ToString())
	fmt.Println(hash3.ToString())
	fmt.Println(hash1.Bits())
	fmt.Println(hash2.Bits())
	fmt.Println(hash3.Bits())
}

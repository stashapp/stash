![GitHub Action](https://github.com/corona10/goimagehash/workflows/goimagehash%20workflow/badge.svg)
[![GoDoc](https://godoc.org/github.com/corona10/goimagehash?status.svg)](https://godoc.org/github.com/corona10/goimagehash)
[![Go Report Card](https://goreportcard.com/badge/github.com/corona10/goimagehash)](https://goreportcard.com/report/github.com/corona10/goimagehash)

# goimagehash
> Inspired by [imagehash](https://github.com/JohannesBuchner/imagehash)

A image hashing library written in Go. ImageHash supports:
* [Average hashing](http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html)
* [Difference hashing](http://www.hackerfactor.com/blog/index.php?/archives/529-Kind-of-Like-That.html)
* [Perception hashing](http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html)
* [Wavelet hashing](https://fullstackml.com/wavelet-image-hash-in-python-3504fdd282b5) [TODO]

## Installation
```
go get github.com/corona10/goimagehash
```
## Special thanks to
* [Haeun Kim](https://github.com/haeungun/)

## Usage

``` Go
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
        width, height := 8, 8
        hash3, _ = goimagehash.ExtAverageHash(img1, width, height)
        hash4, _ = goimagehash.ExtAverageHash(img2, width, height)
        distance, _ = hash3.Distance(hash4)
        fmt.Printf("Distance between images: %v\n", distance)
        fmt.Printf("hash3 bit size: %v\n", hash3.Bits())
        fmt.Printf("hash4 bit size: %v\n", hash4.Bits())

        var b bytes.Buffer
        foo := bufio.NewWriter(&b)
        _ = hash4.Dump(foo)
        foo.Flush()
        bar := bufio.NewReader(&b)
        hash5, _ := goimagehash.LoadExtImageHash(bar)
}
```

## Release Note
### v1.0.3
- Add workflow for GithubAction
- Fix typo on the GoDoc for LoadImageHash

### v1.0.2
- go.mod is now used for install goimagehash

### v1.0.1
- Perception/ExtPerception hash creation times are reduced

### v1.0.0
**IMPORTANT**
goimagehash v1.0.0 does not have compatible with the before version for future features

- More flexible extended hash APIs are provided ([ExtAverageHash](https://godoc.org/github.com/corona10/goimagehash#ExtAverageHash), [ExtPerceptionHash](https://godoc.org/github.com/corona10/goimagehash#ExtPerceptionHash), [ExtDifferenceHash](https://godoc.org/github.com/corona10/goimagehash#ExtDifferenceHash))
- New serialization APIs are provided([ImageHash.Dump](https://godoc.org/github.com/corona10/goimagehash#ImageHash.Dump), [ExtImageHash.Dump](https://godoc.org/github.com/corona10/goimagehash#ExtImageHash.Dump))
- [ExtImageHashFromString](https://godoc.org/github.com/corona10/goimagehash#ExtImageHashFromString), [ImageHashFromString](https://godoc.org/github.com/corona10/goimagehash#ImageHashFromString) is deprecated and will be removed
- New deserialization APIs are provided([LoadImageHash](https://godoc.org/github.com/corona10/goimagehash#LoadImageHash), [LoadExtImageHash](https://godoc.org/github.com/corona10/goimagehash#LoadExtImageHash))
- Bits APIs are provided to measure actual bit size of hash

### v0.3.0
- Support DifferenceHashExtend.
- Support AverageHashExtend.
- Support PerceptionHashExtend by @TokyoWolFrog.

### v0.2.0
- Perception Hash is updated.
- Fix a critical bug of finding median value.

### v0.1.0
- Support Average hashing
- Support Difference hashing
- Support Perception hashing
- Use bits.OnesCount64 for computing Hamming distance by @dominikh
- Support hex serialization methods to ImageHash by @brunoro

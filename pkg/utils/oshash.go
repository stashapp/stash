package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const chunkSize int64 = 64 * 1024

func oshash(size int64, head []byte, tail []byte) (string, error) {
	// put the head and tail together
	buf := append(head, tail...)

	// convert bytes into uint64
	ints := make([]uint64, len(buf)/8)
	reader := bytes.NewReader(buf)
	err := binary.Read(reader, binary.LittleEndian, &ints)
	if err != nil {
		return "", err
	}

	// sum the integers
	var sum uint64
	for _, v := range ints {
		sum += v
	}

	// add the filesize
	sum += uint64(size)

	// output as hex
	return fmt.Sprintf("%016x", sum), nil
}

func OSHashFromReader(src io.ReadSeeker, fileSize int64) (string, error) {
	if fileSize == 0 {
		return "", nil
	}

	fileChunkSize := chunkSize
	if fileSize < fileChunkSize {
		fileChunkSize = fileSize
	}

	head := make([]byte, fileChunkSize)
	tail := make([]byte, fileChunkSize)

	// read the head of the file into the start of the buffer
	_, err := src.Read(head)
	if err != nil {
		return "", err
	}

	// seek to the end of the file - the chunk size
	_, err = src.Seek(-fileChunkSize, 2)
	if err != nil {
		return "", err
	}

	// read the tail of the file
	_, err = src.Read(tail)
	if err != nil {
		return "", err
	}

	return oshash(fileSize, head, tail)
}

// OSHashFromFilePath calculates the hash using the same algorithm that
// OpenSubtitles.org uses.
//
// Calculation is as follows:
// size + 64 bit checksum of the first and last 64k bytes of the file.
func OSHashFromFilePath(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return "", err
	}

	fileSize := fi.Size()

	return OSHashFromReader(f, fileSize)
}

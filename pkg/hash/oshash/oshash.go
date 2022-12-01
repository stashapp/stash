// Package oshash implements the algorithm that OpenSubtitles.org uses to generate unique hashes.
//
// Calculation is as follows:
// size + 64 bit checksum of the first and last 64k bytes of the file.
package oshash

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

const chunkSize int64 = 64 * 1024

var ErrOsHashLen = errors.New("buffer is not a multiple of 8")

func sumBytes(buf []byte) (uint64, error) {
	if len(buf)%8 != 0 {
		return 0, ErrOsHashLen
	}

	sz := len(buf) / 8
	var sum uint64
	for j := 0; j < sz; j++ {
		sum += binary.LittleEndian.Uint64(buf[8*j : 8*(j+1)])
	}

	return sum, nil
}

func oshash(size int64, head []byte, tail []byte) (string, error) {
	headSum, err := sumBytes(head)
	if err != nil {
		return "", fmt.Errorf("oshash head: %w", err)
	}
	tailSum, err := sumBytes(tail)
	if err != nil {
		return "", fmt.Errorf("oshash tail: %w", err)
	}

	// Compute the sum of the head, tail and file size
	result := headSum + tailSum + uint64(size)
	// output as hex
	return fmt.Sprintf("%016x", result), nil
}

// FromReader calculates the hash reading from src.
func FromReader(src io.ReadSeeker, fileSize int64) (string, error) {
	if fileSize <= 8 {
		return "", fmt.Errorf("cannot calculate oshash where size < 8 (%d)", fileSize)
	}

	fileChunkSize := chunkSize
	if fileSize < fileChunkSize {
		// Must be a multiple of 8.
		fileChunkSize = (fileSize / 8) * 8
	}

	head := make([]byte, fileChunkSize)
	tail := make([]byte, fileChunkSize)

	// read the head of the file into the start of the buffer
	_, err := src.Read(head)
	if err != nil {
		return "", err
	}

	// seek to the end of the file - the chunk size
	_, err = src.Seek(-fileChunkSize, io.SeekEnd)
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

// Is the equivalent of opening filePath, and calling FromReader with the data and file size.
func FromFilePath(filePath string) (string, error) {
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

	return FromReader(f, fileSize)
}

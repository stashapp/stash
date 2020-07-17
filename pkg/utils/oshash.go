package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

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

	fileSize := int64(fi.Size())

	if fileSize == 0 {
		return "", nil
	}

	const chunkSize = 64 * 1024
	fileChunkSize := int64(chunkSize)
	if fileSize < fileChunkSize {
		fileChunkSize = fileSize
	}

	head := make([]byte, fileChunkSize)
	tail := make([]byte, fileChunkSize)

	// read the head of the file into the start of the buffer
	_, err = f.Read(head)
	if err != nil {
		return "", err
	}

	// seek to the end of the file - the chunk size
	_, err = f.Seek(-fileChunkSize, 2)
	if err != nil {
		return "", err
	}

	// read the tail of the file
	_, err = f.Read(tail)
	if err != nil {
		return "", err
	}

	// put the head and tail together
	buf := append(head, tail...)

	// convert bytes into uint64
	ints := make([]uint64, len(buf)/8)
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &ints)
	if err != nil {
		return "", err
	}

	// sum the integers
	var sum uint64
	for _, v := range ints {
		sum += v
	}

	// add the filesize
	sum += uint64(fileSize)

	// output as hex
	return fmt.Sprintf("%016x", sum), nil
}

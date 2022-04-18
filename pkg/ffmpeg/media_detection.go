package ffmpeg

import (
	"bytes"
	"os"
)

// detect file format from magic file number
// https://github.com/lex-r/filetype/blob/73c10ad714e3b8ecf5cd1564c882ed6d440d5c2d/matchers/video.go

func mkv(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x1A && buf[1] == 0x45 &&
		buf[2] == 0xDF && buf[3] == 0xA3 &&
		containsMatroskaSignature(buf, []byte{'m', 'a', 't', 'r', 'o', 's', 'k', 'a'})
}

func webm(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x1A && buf[1] == 0x45 &&
		buf[2] == 0xDF && buf[3] == 0xA3 &&
		containsMatroskaSignature(buf, []byte{'w', 'e', 'b', 'm'})
}

func containsMatroskaSignature(buf, subType []byte) bool {
	limit := 4096
	if len(buf) < limit {
		limit = len(buf)
	}

	index := bytes.Index(buf[:limit], subType)
	if index < 3 {
		return false
	}

	return buf[index-3] == 0x42 && buf[index-2] == 0x82
}

// magicContainer returns the container type of a file path.
// Returns the zero-value on errors or no-match. Implements mkv or
// webm only, as ffprobe can't distinguish between them and not all
// browsers support mkv
func magicContainer(filePath string) (Container, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	buf := make([]byte, 4096)
	_, err = file.Read(buf)
	if err != nil {
		return "", err
	}

	if webm(buf) {
		return Webm, nil
	}
	if mkv(buf) {
		return Matroska, nil
	}
	return "", nil
}

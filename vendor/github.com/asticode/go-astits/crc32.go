package astits

const (
	crc32Polynomial = uint32(0xffffffff)
)

func computeCRC32(bs []byte) uint32 {
	return updateCRC32(crc32Polynomial, bs)
}

// Based on VLC implementation using a static CRC table (1kb additional memory on start, without
// reallocations): https://github.com/videolan/vlc/blob/master/modules/mux/mpeg/ps.c
func updateCRC32(crc32 uint32, bs []byte) uint32 {
	for _, b := range bs {
		crc32 = (crc32 << 8) ^ tableCRC32[((crc32>>24)^uint32(b))&0xff]
	}
	return crc32
}

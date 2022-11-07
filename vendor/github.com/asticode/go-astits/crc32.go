package astits

const (
	crc32Polynomial = uint32(0xffffffff)
)

// computeCRC32 computes a CRC32
// https://stackoverflow.com/questions/35034042/how-to-calculate-crc32-in-psi-si-packet
func computeCRC32(bs []byte) uint32 {
	return updateCRC32(crc32Polynomial, bs)
}

func updateCRC32(crc32 uint32, bs []byte) uint32 {
	for _, b := range bs {
		for i := 0; i < 8; i++ {
			if (crc32 >= uint32(0x80000000)) != (b >= uint8(0x80)) {
				crc32 = (crc32 << 1) ^ 0x04C11DB7
			} else {
				crc32 = crc32 << 1
			}
			b <<= 1
		}
	}
	return crc32
}

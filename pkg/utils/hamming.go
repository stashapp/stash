// Copied from https://github.com/rivo/duplo/blob/master/hamming.go
package utils

const (
	m1  = 0x5555555555555555 //binary: 0101...
	m2  = 0x3333333333333333 //binary: 00110011..
	m4  = 0x0f0f0f0f0f0f0f0f //binary:  4 zeros,  4 ones ...
	m8  = 0x00ff00ff00ff00ff //binary:  8 zeros,  8 ones ...
	m16 = 0x0000ffff0000ffff //binary: 16 zeros, 16 ones ...
	m32 = 0x00000000ffffffff //binary: 32 zeros, 32 ones
	hff = 0xffffffffffffffff //binary: all ones
	h01 = 0x0101010101010101 //the sum of 256 to the power of 0,1,2,3...
)

// hammingDistance calculates the hamming distance between two 64-bit values.
// The implementation is based on the code found on:
// http://en.wikipedia.org/wiki/Hamming_weight#Efficient_implementation
func HammingDistance(left, right uint64) int {
	x := left ^ right
	x -= (x >> 1) & m1             //put count of each 2 bits into those 2 bits
	x = (x & m2) + ((x >> 2) & m2) //put count of each 4 bits into those 4 bits
	x = (x + (x >> 4)) & m4        //put count of each 8 bits into those 8 bits
	return int((x * h01) >> 56)    //returns left 8 bits of x + (x<<8) + (x<<16) + (x<<24) + ...
}

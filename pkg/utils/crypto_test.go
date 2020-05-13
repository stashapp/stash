package utils

import (
	"fmt"
	"strings"
	"testing"
)

var md5s = []string{
	"00000000000000000000000000000000",  //00
	"0af63ce3c99162e9df23a997f62621c5",  //01
	"1be40aeeaacaf9a55f44441143ac4270",  //02
	"8c2cb8a46335fdec761626357b012eb7",  //03
	"9d3e8ea9006e66a0f671cbb3ce8b4d02",  //04 xor(1,2,3)
	"af63ce3c99162e9df23a997f62621c5",   //05 same as 1 different length
	"00a5ea27bf5cdbc6a11de570b9088afa",  //06
	"0a5ea27bf5cdbc6a11de570b9088afa",   //07 6 with different length
	"a5ea27bf5cdbc6a11de570b9088afa",    //08 6 with different length
	"1112360d635b9b4c8067ed86b58a63b5",  //09 xor(1,2)
	"adscx",                             //10 not valid hex
	"ad3f2",                             //11
	"1112360d635b9b4c8067ed86b580b047",  //12 xor(9,11)
	"1112360d635b9b4c8067ed86b580b0475", //13 not valid md5  >32 length
	"3d5dff",                            //14
	"00000000000000000000000000378e0d",  //15 xor(14,11)

}

var returnErrors = []string{
	"",                   //00
	"invalid byte",       //01
	"is not a valid MD5", //02
}

var xorTests = []struct {
	md5sIndex   []int
	resultIndex int
	errorIndex  int
}{ //tests
	{[]int{1, 1}, 0, 0},    //01. xor with self equals 0
	{[]int{1, 5}, 0, 0},    //02
	{[]int{6, 7}, 0, 0},    //03
	{[]int{6, 8}, 0, 0},    //04
	{[]int{7, 8}, 0, 0},    //05
	{[]int{1, 2}, 9, 0},    //06
	{[]int{1, 0}, 1, 0},    //07 xor with 0 equals self
	{[]int{1, 2, 1}, 2, 0}, //08 xor order doesn't matter
	{[]int{1, 1, 2}, 2, 0}, //09
	{[]int{1, 2, 3}, 4, 0}, //10
	{[]int{1, 3, 2}, 4, 0}, //11
	{[]int{2, 1, 3}, 4, 0}, //12
	{[]int{2, 3, 1}, 4, 0}, //13
	{[]int{3, 1, 2}, 4, 0}, //14
	{[]int{3, 2, 1}, 4, 0}, //15
	{[]int{9, 11}, 12, 0},  //16 xor diff length
	{[]int{10, 11}, 12, 1}, //17 invalid hex error (11 not a  hex)
	{[]int{13, 11}, 12, 2}, //18 not valid md5 error (13 not a MD5)
	{[]int{14, 11}, 15, 0}, //19

}

func TestXorMD5Strings(t *testing.T) {
	for i, test := range xorTests {
		var checksums []string
		for _, n := range test.md5sIndex {
			checksums = append(checksums, md5s[n])
		}
		result, err := XorMD5Strings(checksums)
		if err != nil {
			wanted := returnErrors[test.errorIndex]
			if wanted == "" || !strings.Contains(err.Error(), wanted) {
				t.Error(fmt.Errorf("Test %02d: %s", i+1, err))
			}
		} else {
			if result != md5s[test.resultIndex] {
				t.Error(fmt.Errorf("Test %02d: Was expecting %s, found %s", i+1, md5s[test.resultIndex], result))
			}
		}
	}
}

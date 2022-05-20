package astisub

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/bits"
	"sort"
	"strings"
	"time"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astits"
)

// Errors
var (
	ErrNoValidTeletextPID = errors.New("astisub: no valid teletext PID")
)

type teletextCharset [96][]byte

type teletextNationalSubset [13][]byte

// Chapter: 15.2 | Page: 109 | Link: http://www.etsi.org/deliver/etsi_i_ets/300700_300799/300706/01_60/ets_300706e01p.pdf
// It is indexed by triplet1 then by national option subset code
var teletextCharsets = map[uint8]map[uint8]struct {
	g0       *teletextCharset
	g2       *teletextCharset
	national *teletextNationalSubset
}{
	0: {
		0: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetEnglish},
		1: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetFrench},
		2: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetSwedishFinnishHungarian},
		3: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetCzechSlovak},
		4: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetGerman},
		5: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetPortugueseSpanish},
		6: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetItalian},
		7: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin},
	},
	1: {
		0: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetPolish},
		1: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetFrench},
		2: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetSwedishFinnishHungarian},
		3: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetCzechSlovak},
		4: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetGerman},
		5: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin},
		6: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetItalian},
		7: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin},
	},
	2: {
		0: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetEnglish},
		1: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetFrench},
		2: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetSwedishFinnishHungarian},
		3: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetCzechSlovak},
		4: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetGerman},
		5: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetPortugueseSpanish},
		6: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetItalian},
		7: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin},
	},
	3: {
		0: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin},
		1: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin},
		2: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin},
		3: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin},
		4: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin},
		5: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetSerbianCroatianSlovenian},
		6: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin},
		7: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetRomanian},
	},
	4: {
		0: {g0: teletextCharsetG0CyrillicOption1, g2: teletextCharsetG2Cyrillic},
		1: {g0: teletextCharsetG0CyrillicOption2, g2: teletextCharsetG2Cyrillic},
		2: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetEstonian},
		3: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetCzechSlovak},
		4: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetGerman},
		5: {g0: teletextCharsetG0CyrillicOption3, g2: teletextCharsetG2Cyrillic},
		6: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetLettishLithuanian},
	},
	6: {
		3: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Latin, national: teletextNationalSubsetTurkish},
		7: {g0: teletextCharsetG0Greek, g2: teletextCharsetG2Greek},
	},
	8: {
		0: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Arabic, national: teletextNationalSubsetEnglish},
		1: {g0: teletextCharsetG0Latin, g2: teletextCharsetG2Arabic, national: teletextNationalSubsetFrench},
		7: {g0: teletextCharsetG0Arabic, g2: teletextCharsetG2Arabic},
	},
	10: {
		5: {g0: teletextCharsetG0Hebrew, g2: teletextCharsetG2Arabic},
		7: {g0: teletextCharsetG0Arabic, g2: teletextCharsetG2Arabic},
	},
}

// Teletext G0 charsets
var (
	teletextCharsetG0CyrillicOption1 = &teletextCharset{
		[]byte{0x20}, []byte{0x21}, []byte{0x22}, []byte{0x23}, []byte{0x24}, []byte{0x25}, []byte{0xd1, 0x8b},
		[]byte{0x27}, []byte{0x28}, []byte{0x29}, []byte{0x2a}, []byte{0x2b}, []byte{0x2c}, []byte{0x2d},
		[]byte{0x2e}, []byte{0x2f}, []byte{0x30}, []byte{0x31}, []byte{0xe3, 0x88, 0x80}, []byte{0x33}, []byte{0x34},
		[]byte{0x35}, []byte{0x36}, []byte{0x37}, []byte{0x38}, []byte{0x39}, []byte{0x3a}, []byte{0x3b},
		[]byte{0x3c}, []byte{0x3d}, []byte{0x3e}, []byte{0x3f}, []byte{0xd0, 0xa7}, []byte{0xd0, 0x90},
		[]byte{0xd0, 0x91}, []byte{0xd0, 0xa6}, []byte{0xd0, 0x94}, []byte{0xd0, 0x95}, []byte{0xd0, 0xa4},
		[]byte{0xd0, 0x93}, []byte{0xd0, 0xa5}, []byte{0xd0, 0x98}, []byte{0xd0, 0x88}, []byte{0xd0, 0x9a},
		[]byte{0xd0, 0x9b}, []byte{0xd0, 0x9c}, []byte{0xd0, 0x9d}, []byte{0xd0, 0x9e}, []byte{0xd0, 0x9f},
		[]byte{0xd0, 0x8c}, []byte{0xd0, 0xa0}, []byte{0xd0, 0xa1}, []byte{0xd0, 0xa2}, []byte{0xd0, 0xa3},
		[]byte{0xd0, 0x92}, []byte{0xd0, 0x83}, []byte{0xd0, 0x89}, []byte{0xd0, 0x8a}, []byte{0xd0, 0x97},
		[]byte{0xd0, 0x8b}, []byte{0xd0, 0x96}, []byte{0xd0, 0x82}, []byte{0xd0, 0xa8}, []byte{0xd0, 0x8f},
		[]byte{0xd1, 0x87}, []byte{0xd0, 0xb0}, []byte{0xd0, 0xb1}, []byte{0xd1, 0x86}, []byte{0xd0, 0xb4},
		[]byte{0xd0, 0xb5}, []byte{0xd1, 0x84}, []byte{0xd0, 0xb3}, []byte{0xd1, 0x85}, []byte{0xd0, 0xb8},
		[]byte{0xd0, 0xa8}, []byte{0xd0, 0xba}, []byte{0xd0, 0xbb}, []byte{0xd0, 0xbc}, []byte{0xd0, 0xbd},
		[]byte{0xd0, 0xbe}, []byte{0xd0, 0xbf}, []byte{0xd0, 0xac}, []byte{0xd1, 0x80}, []byte{0xd1, 0x81},
		[]byte{0xd1, 0x82}, []byte{0xd1, 0x83}, []byte{0xd0, 0xb2}, []byte{0xd0, 0xa3}, []byte{0xd0, 0xa9},
		[]byte{0xd0, 0xaa}, []byte{0xd0, 0xb7}, []byte{0xd0, 0xab}, []byte{0xd0, 0xb6}, []byte{0xd0, 0xa2},
		[]byte{0xd1, 0x88}, []byte{0xd0, 0xaf},
	}
	teletextCharsetG0CyrillicOption2 = &teletextCharset{
		[]byte{0x20}, []byte{0x21}, []byte{0x22}, []byte{0x23}, []byte{0x24}, []byte{0x25}, []byte{0xd1, 0x8b},
		[]byte{0x27}, []byte{0x28}, []byte{0x29}, []byte{0x2a}, []byte{0x2b}, []byte{0x2c}, []byte{0x2d},
		[]byte{0x2e}, []byte{0x2f}, []byte{0x30}, []byte{0x31}, []byte{0x32}, []byte{0x33}, []byte{0x34},
		[]byte{0x35}, []byte{0x36}, []byte{0x37}, []byte{0x38}, []byte{0x39}, []byte{0x3a}, []byte{0x3b},
		[]byte{0x3c}, []byte{0x3d}, []byte{0x3e}, []byte{0x3f}, []byte{0xd0, 0xae}, []byte{0xd0, 0x90},
		[]byte{0xd0, 0x91}, []byte{0xd0, 0xa6}, []byte{0xd0, 0x94}, []byte{0xd0, 0x95}, []byte{0xd0, 0xa4},
		[]byte{0xd0, 0x93}, []byte{0xd0, 0xa5}, []byte{0xd0, 0x98}, []byte{0xd0, 0x99}, []byte{0xd0, 0x9a},
		[]byte{0xd0, 0x9b}, []byte{0xd0, 0x9c}, []byte{0xd0, 0x9d}, []byte{0xd0, 0x9e}, []byte{0xd0, 0x9f},
		[]byte{0xd0, 0xaf}, []byte{0xd0, 0xa0}, []byte{0xd0, 0xa1}, []byte{0xd0, 0xa2}, []byte{0xd0, 0xa3},
		[]byte{0xd0, 0x96}, []byte{0xd0, 0x92}, []byte{0xd0, 0xac}, []byte{0xd0, 0xaa}, []byte{0xd0, 0x97},
		[]byte{0xd0, 0xa8}, []byte{0xd0, 0xad}, []byte{0xd0, 0xa9}, []byte{0xd0, 0xa7}, []byte{0xd0, 0xab},
		[]byte{0xd1, 0x8e}, []byte{0xd0, 0xb0}, []byte{0xd0, 0xb1}, []byte{0xd1, 0x86}, []byte{0xd0, 0xb4},
		[]byte{0xd0, 0xb5}, []byte{0xd1, 0x84}, []byte{0xd0, 0xb3}, []byte{0xd1, 0x85}, []byte{0xd0, 0xb8},
		[]byte{0xd0, 0xb9}, []byte{0xd0, 0xba}, []byte{0xd0, 0xbb}, []byte{0xd0, 0xbc}, []byte{0xd0, 0xbd},
		[]byte{0xd0, 0xbe}, []byte{0xd0, 0xbf}, []byte{0xd1, 0x8f}, []byte{0xd1, 0x80}, []byte{0xd1, 0x81},
		[]byte{0xd1, 0x82}, []byte{0xd1, 0x83}, []byte{0xd0, 0xb6}, []byte{0xd0, 0xb2}, []byte{0xd1, 0x8c},
		[]byte{0xd1, 0x8a}, []byte{0xd0, 0xb7}, []byte{0xd1, 0x88}, []byte{0xd1, 0x8d}, []byte{0xd1, 0x89},
		[]byte{0xd1, 0x87}, []byte{0xd1, 0x8b},
	}
	teletextCharsetG0CyrillicOption3 = &teletextCharset{
		[]byte{0x20}, []byte{0x21}, []byte{0x22}, []byte{0x23}, []byte{0x24}, []byte{0x25}, []byte{0xc3, 0xaf},
		[]byte{0x27}, []byte{0x28}, []byte{0x29}, []byte{0x2a}, []byte{0x2b}, []byte{0x2c}, []byte{0x2d},
		[]byte{0x2e}, []byte{0x2f}, []byte{0x30}, []byte{0x31}, []byte{0x32}, []byte{0x33}, []byte{0x34},
		[]byte{0x35}, []byte{0x36}, []byte{0x37}, []byte{0x38}, []byte{0x39}, []byte{0x3a}, []byte{0x3b},
		[]byte{0x3c}, []byte{0x3d}, []byte{0x3e}, []byte{0x3f}, []byte{0xd0, 0xae}, []byte{0xd0, 0x90},
		[]byte{0xd0, 0x91}, []byte{0xd0, 0xa6}, []byte{0xd0, 0x94}, []byte{0xd0, 0x95}, []byte{0xd0, 0xa4},
		[]byte{0xd0, 0x93}, []byte{0xd0, 0xa5}, []byte{0xd0, 0x98}, []byte{0xd0, 0x99}, []byte{0xd0, 0x9a},
		[]byte{0xd0, 0x9b}, []byte{0xd0, 0x9c}, []byte{0xd0, 0x9d}, []byte{0xd0, 0x9e}, []byte{0xd0, 0x9f},
		[]byte{0xd0, 0xaf}, []byte{0xd0, 0xa0}, []byte{0xd0, 0xa1}, []byte{0xd0, 0xa2}, []byte{0xd0, 0xa3},
		[]byte{0xd0, 0x96}, []byte{0xd0, 0x92}, []byte{0xd0, 0xac}, []byte{0x49}, []byte{0xd0, 0x97},
		[]byte{0xd0, 0xa8}, []byte{0xd0, 0xad}, []byte{0xd0, 0xa9}, []byte{0xd0, 0xa7}, []byte{0xc3, 0x8f},
		[]byte{0xd1, 0x8e}, []byte{0xd0, 0xb0}, []byte{0xd0, 0xb1}, []byte{0xd1, 0x86}, []byte{0xd0, 0xb4},
		[]byte{0xd0, 0xb5}, []byte{0xd1, 0x84}, []byte{0xd0, 0xb3}, []byte{0xd1, 0x85}, []byte{0xd0, 0xb8},
		[]byte{0xd0, 0xb9}, []byte{0xd0, 0xba}, []byte{0xd0, 0xbb}, []byte{0xd0, 0xbc}, []byte{0xd0, 0xbd},
		[]byte{0xd0, 0xbe}, []byte{0xd0, 0xbf}, []byte{0xd1, 0x8f}, []byte{0xd1, 0x80}, []byte{0xd1, 0x81},
		[]byte{0xd1, 0x82}, []byte{0xd1, 0x83}, []byte{0xd0, 0xb6}, []byte{0xd0, 0xb2}, []byte{0xd1, 0x8c},
		[]byte{0x69}, []byte{0xd0, 0xb7}, []byte{0xd1, 0x88}, []byte{0xd1, 0x8d}, []byte{0xd1, 0x89},
		[]byte{0xd1, 0x87}, []byte{0xc3, 0xbf},
	}
	teletextCharsetG0Greek = &teletextCharset{
		[]byte{0x20}, []byte{0x21}, []byte{0x22}, []byte{0x23}, []byte{0x24}, []byte{0x25}, []byte{0x26},
		[]byte{0x27}, []byte{0x28}, []byte{0x29}, []byte{0x2a}, []byte{0x2b}, []byte{0x2c}, []byte{0x2d},
		[]byte{0x2e}, []byte{0x2f}, []byte{0x30}, []byte{0x31}, []byte{0x32}, []byte{0x33}, []byte{0x34},
		[]byte{0x35}, []byte{0x36}, []byte{0x37}, []byte{0x38}, []byte{0x39}, []byte{0x3a}, []byte{0x3b},
		[]byte{0x3c}, []byte{0x3d}, []byte{0x3e}, []byte{0x3f}, []byte{0xce, 0x90}, []byte{0xce, 0x91},
		[]byte{0xce, 0x92}, []byte{0xce, 0x93}, []byte{0xce, 0x94}, []byte{0xce, 0x95}, []byte{0xce, 0x96},
		[]byte{0xce, 0x97}, []byte{0xce, 0x98}, []byte{0xce, 0x99}, []byte{0xce, 0x9a}, []byte{0xce, 0x9b},
		[]byte{0xce, 0x9c}, []byte{0xce, 0x9d}, []byte{0xce, 0x9e}, []byte{0xce, 0x9f}, []byte{0xce, 0xa0},
		[]byte{0xce, 0xa1}, []byte{0xce, 0xa2}, []byte{0xce, 0xa3}, []byte{0xce, 0xa4}, []byte{0xce, 0xa5},
		[]byte{0xce, 0xa6}, []byte{0xce, 0xa7}, []byte{0xce, 0xa8}, []byte{0xce, 0xa9}, []byte{0xce, 0xaa},
		[]byte{0xce, 0xab}, []byte{0xce, 0xac}, []byte{0xce, 0xad}, []byte{0xce, 0xae}, []byte{0xce, 0xaf},
		[]byte{0xce, 0xb0}, []byte{0xce, 0xb1}, []byte{0xce, 0xb2}, []byte{0xce, 0xb3}, []byte{0xce, 0xb4},
		[]byte{0xce, 0xb5}, []byte{0xce, 0xb6}, []byte{0xce, 0xb7}, []byte{0xce, 0xb8}, []byte{0xce, 0xb9},
		[]byte{0xce, 0xba}, []byte{0xce, 0xbb}, []byte{0xce, 0xbc}, []byte{0xce, 0xbd}, []byte{0xce, 0xbe},
		[]byte{0xce, 0xbf}, []byte{0xcf, 0x80}, []byte{0xcf, 0x81}, []byte{0xcf, 0x82}, []byte{0xcf, 0x83},
		[]byte{0xcf, 0x84}, []byte{0xcf, 0x85}, []byte{0xcf, 0x86}, []byte{0xcf, 0x87}, []byte{0xcf, 0x88},
		[]byte{0xcf, 0x89}, []byte{0xcf, 0x8a}, []byte{0xcf, 0x8b}, []byte{0xcf, 0x8c}, []byte{0xcf, 0x8d},
		[]byte{0xcf, 0x8e}, []byte{0xcf, 0x8f},
	}
	teletextCharsetG0Latin = &teletextCharset{
		[]byte{0x20}, []byte{0x21}, []byte{0x22}, []byte{0xc2, 0xa3}, []byte{0x24}, []byte{0x25}, []byte{0x26},
		[]byte{0x27}, []byte{0x28}, []byte{0x29}, []byte{0x2a}, []byte{0x2b}, []byte{0x2c}, []byte{0x2d},
		[]byte{0x2e}, []byte{0x2f}, []byte{0x30}, []byte{0x31}, []byte{0x32}, []byte{0x33}, []byte{0x34},
		[]byte{0x35}, []byte{0x36}, []byte{0x37}, []byte{0x38}, []byte{0x39}, []byte{0x3a}, []byte{0x3b},
		[]byte{0x3c}, []byte{0x3d}, []byte{0x3e}, []byte{0x3f}, []byte{0x40}, []byte{0x41}, []byte{0x42},
		[]byte{0x43}, []byte{0x44}, []byte{0x45}, []byte{0x46}, []byte{0x47}, []byte{0x48}, []byte{0x49},
		[]byte{0x4a}, []byte{0x4b}, []byte{0x4c}, []byte{0x4d}, []byte{0x4e}, []byte{0x4f}, []byte{0x50},
		[]byte{0x51}, []byte{0x52}, []byte{0x53}, []byte{0x54}, []byte{0x55}, []byte{0x56}, []byte{0x57},
		[]byte{0x58}, []byte{0x59}, []byte{0x5a}, []byte{0xc2, 0xab}, []byte{0xc2, 0xbd}, []byte{0xc2, 0xbb},
		[]byte{0x5e}, []byte{0x23}, []byte{0x2d}, []byte{0x61}, []byte{0x62}, []byte{0x63}, []byte{0x64},
		[]byte{0x65}, []byte{0x66}, []byte{0x67}, []byte{0x68}, []byte{0x69}, []byte{0x6a}, []byte{0x6b},
		[]byte{0x6c}, []byte{0x6d}, []byte{0x6e}, []byte{0x6f}, []byte{0x70}, []byte{0x71}, []byte{0x72},
		[]byte{0x73}, []byte{0x74}, []byte{0x75}, []byte{0x76}, []byte{0x77}, []byte{0x78}, []byte{0x79},
		[]byte{0x7a}, []byte{0xc2, 0xbc}, []byte{0xc2, 0xa6}, []byte{0xc2, 0xbe}, []byte{0xc3, 0xb7}, []byte{0x7f},
	}
	// TODO Add
	teletextCharsetG0Arabic = teletextCharsetG0Latin
	teletextCharsetG0Hebrew = teletextCharsetG0Latin
)

// Teletext G2 charsets
var (
	teletextCharsetG2Latin = &teletextCharset{
		[]byte{0x20}, []byte{0xc2, 0xa1}, []byte{0xc2, 0xa2}, []byte{0xc2, 0xa3}, []byte{0x24},
		[]byte{0xc2, 0xa5}, []byte{0x23}, []byte{0xc2, 0xa7}, []byte{0xc2, 0xa4}, []byte{0xe2, 0x80, 0x98},
		[]byte{0xe2, 0x80, 0x9c}, []byte{0xc2, 0xab}, []byte{0xe2, 0x86, 0x90}, []byte{0xe2, 0x86, 0x91},
		[]byte{0xe2, 0x86, 0x92}, []byte{0xe2, 0x86, 0x93}, []byte{0xc2, 0xb0}, []byte{0xc2, 0xb1},
		[]byte{0xc2, 0xb2}, []byte{0xc2, 0xb3}, []byte{0xc3, 0x97}, []byte{0xc2, 0xb5}, []byte{0xc2, 0xb6},
		[]byte{0xc2, 0xb7}, []byte{0xc3, 0xb7}, []byte{0xe2, 0x80, 0x99}, []byte{0xe2, 0x80, 0x9d},
		[]byte{0xc2, 0xbb}, []byte{0xc2, 0xbc}, []byte{0xc2, 0xbd}, []byte{0xc2, 0xbe}, []byte{0xc2, 0xbf},
		[]byte{0x20}, []byte{0xcc, 0x80}, []byte{0xcc, 0x81}, []byte{0xcc, 0x82}, []byte{0xcc, 0x83},
		[]byte{0xcc, 0x84}, []byte{0xcc, 0x86}, []byte{0xcc, 0x87}, []byte{0xcc, 0x88}, []byte{0x00},
		[]byte{0xcc, 0x8a}, []byte{0xcc, 0xa7}, []byte{0x5f}, []byte{0xcc, 0x8b}, []byte{0xcc, 0xa8},
		[]byte{0xcc, 0x8c}, []byte{0xe2, 0x80, 0x95}, []byte{0xc2, 0xb9}, []byte{0xc2, 0xae}, []byte{0xc2, 0xa9},
		[]byte{0xe2, 0x84, 0xa2}, []byte{0xe2, 0x99, 0xaa}, []byte{0xe2, 0x82, 0xac}, []byte{0xe2, 0x80, 0xb0},
		[]byte{0xce, 0xb1}, []byte{0x00}, []byte{0x00}, []byte{0x00}, []byte{0xe2, 0x85, 0x9b},
		[]byte{0xe2, 0x85, 0x9c}, []byte{0xe2, 0x85, 0x9d}, []byte{0xe2, 0x85, 0x9e}, []byte{0xce, 0xa9},
		[]byte{0xc3, 0x86}, []byte{0xc4, 0x90}, []byte{0xc2, 0xaa}, []byte{0xc4, 0xa6}, []byte{0x00},
		[]byte{0xc4, 0xb2}, []byte{0xc4, 0xbf}, []byte{0xc5, 0x81}, []byte{0xc3, 0x98}, []byte{0xc5, 0x92},
		[]byte{0xc2, 0xba}, []byte{0xc3, 0x9e}, []byte{0xc5, 0xa6}, []byte{0xc5, 0x8a}, []byte{0xc5, 0x89},
		[]byte{0xc4, 0xb8}, []byte{0xc3, 0xa6}, []byte{0xc4, 0x91}, []byte{0xc3, 0xb0}, []byte{0xc4, 0xa7},
		[]byte{0xc4, 0xb1}, []byte{0xc4, 0xb3}, []byte{0xc5, 0x80}, []byte{0xc5, 0x82}, []byte{0xc3, 0xb8},
		[]byte{0xc5, 0x93}, []byte{0xc3, 0x9f}, []byte{0xc3, 0xbe}, []byte{0xc5, 0xa7}, []byte{0xc5, 0x8b},
		[]byte{0x20},
	}
	// TODO Add
	teletextCharsetG2Arabic   = teletextCharsetG2Latin
	teletextCharsetG2Cyrillic = teletextCharsetG2Latin
	teletextCharsetG2Greek    = teletextCharsetG2Latin
)

var teletextNationalSubsetCharactersPositionInG0 = [13]uint8{0x03, 0x04, 0x20, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f, 0x40, 0x5b, 0x5c, 0x5d, 0x5e}

// Teletext national subsets
var (
	teletextNationalSubsetCzechSlovak = &teletextNationalSubset{
		[]byte{0x23}, []byte{0xc5, 0xaf}, []byte{0xc4, 0x8d}, []byte{0xc5, 0xa5}, []byte{0xc5, 0xbe},
		[]byte{0xc3, 0xbd}, []byte{0xc3, 0xad}, []byte{0xc5, 0x99}, []byte{0xc3, 0xa9}, []byte{0xc3, 0xa1},
		[]byte{0xc4, 0x9b}, []byte{0xc3, 0xba}, []byte{0xc5, 0xa1},
	}
	teletextNationalSubsetEnglish = &teletextNationalSubset{
		[]byte{0xc2, 0xa3}, []byte{0x24}, []byte{0x40}, []byte{0xc2, 0xab}, []byte{0xc2, 0xbd}, []byte{0xc2, 0xbb},
		[]byte{0x5e}, []byte{0x23}, []byte{0x2d}, []byte{0xc2, 0xbc}, []byte{0xc2, 0xa6}, []byte{0xc2, 0xbe},
		[]byte{0xc3, 0xb7},
	}
	teletextNationalSubsetEstonian = &teletextNationalSubset{
		[]byte{0x23}, []byte{0xc3, 0xb5}, []byte{0xc5, 0xa0}, []byte{0xc3, 0x84}, []byte{0xc3, 0x96},
		[]byte{0xc5, 0xbe}, []byte{0xc3, 0x9c}, []byte{0xc3, 0x95}, []byte{0xc5, 0xa1}, []byte{0xc3, 0xa4},
		[]byte{0xc3, 0xb6}, []byte{0xc5, 0xbe}, []byte{0xc3, 0xbc},
	}
	teletextNationalSubsetFrench = &teletextNationalSubset{
		[]byte{0xc3, 0xa9}, []byte{0xc3, 0xaf}, []byte{0xc3, 0xa0}, []byte{0xc3, 0xab}, []byte{0xc3, 0xaa},
		[]byte{0xc3, 0xb9}, []byte{0xc3, 0xae}, []byte{0x23}, []byte{0xc3, 0xa8}, []byte{0xc3, 0xa2},
		[]byte{0xc3, 0xb4}, []byte{0xc3, 0xbb}, []byte{0xc3, 0xa7},
	}
	teletextNationalSubsetGerman = &teletextNationalSubset{
		[]byte{0x23}, []byte{0x24}, []byte{0xc2, 0xa7}, []byte{0xc3, 0x84}, []byte{0xc3, 0x96}, []byte{0xc3, 0x9c},
		[]byte{0x5e}, []byte{0x5f}, []byte{0xc2, 0xb0}, []byte{0xc3, 0xa4}, []byte{0xc3, 0xb6}, []byte{0xc3, 0xbc},
		[]byte{0xc3, 0x9f},
	}
	teletextNationalSubsetItalian = &teletextNationalSubset{
		[]byte{0xc2, 0xa3}, []byte{0x24}, []byte{0xc3, 0xa9}, []byte{0xc2, 0xb0}, []byte{0xc3, 0xa7},
		[]byte{0xc2, 0xbb}, []byte{0x5e}, []byte{0x23}, []byte{0xc3, 0xb9}, []byte{0xc3, 0xa0}, []byte{0xc3, 0xb2},
		[]byte{0xc3, 0xa8}, []byte{0xc3, 0xac},
	}
	teletextNationalSubsetLettishLithuanian = &teletextNationalSubset{
		[]byte{0x23}, []byte{0x24}, []byte{0xc5, 0xa0}, []byte{0xc4, 0x97}, []byte{0xc4, 0x99}, []byte{0xc5, 0xbd},
		[]byte{0xc4, 0x8d}, []byte{0xc5, 0xab}, []byte{0xc5, 0xa1}, []byte{0xc4, 0x85}, []byte{0xc5, 0xb3},
		[]byte{0xc5, 0xbe}, []byte{0xc4, 0xaf},
	}
	teletextNationalSubsetPolish = &teletextNationalSubset{
		[]byte{0x23}, []byte{0xc5, 0x84}, []byte{0xc4, 0x85}, []byte{0xc5, 0xbb}, []byte{0xc5, 0x9a},
		[]byte{0xc5, 0x81}, []byte{0xc4, 0x87}, []byte{0xc3, 0xb3}, []byte{0xc4, 0x99}, []byte{0xc5, 0xbc},
		[]byte{0xc5, 0x9b}, []byte{0xc5, 0x82}, []byte{0xc5, 0xba},
	}
	teletextNationalSubsetPortugueseSpanish = &teletextNationalSubset{
		[]byte{0xc3, 0xa7}, []byte{0x24}, []byte{0xc2, 0xa1}, []byte{0xc3, 0xa1}, []byte{0xc3, 0xa9},
		[]byte{0xc3, 0xad}, []byte{0xc3, 0xb3}, []byte{0xc3, 0xba}, []byte{0xc2, 0xbf}, []byte{0xc3, 0xbc},
		[]byte{0xc3, 0xb1}, []byte{0xc3, 0xa8}, []byte{0xc3, 0xa0},
	}
	teletextNationalSubsetRomanian = &teletextNationalSubset{
		[]byte{0x23}, []byte{0xc2, 0xa4}, []byte{0xc5, 0xa2}, []byte{0xc3, 0x82}, []byte{0xc5, 0x9e},
		[]byte{0xc4, 0x82}, []byte{0xc3, 0x8e}, []byte{0xc4, 0xb1}, []byte{0xc5, 0xa3}, []byte{0xc3, 0xa2},
		[]byte{0xc5, 0x9f}, []byte{0xc4, 0x83}, []byte{0xc3, 0xae},
	}
	teletextNationalSubsetSerbianCroatianSlovenian = &teletextNationalSubset{
		[]byte{0x23}, []byte{0xc3, 0x8b}, []byte{0xc4, 0x8c}, []byte{0xc4, 0x86}, []byte{0xc5, 0xbd},
		[]byte{0xc4, 0x90}, []byte{0xc5, 0xa0}, []byte{0xc3, 0xab}, []byte{0xc4, 0x8d}, []byte{0xc4, 0x87},
		[]byte{0xc5, 0xbe}, []byte{0xc4, 0x91}, []byte{0xc5, 0xa1},
	}
	teletextNationalSubsetSwedishFinnishHungarian = &teletextNationalSubset{
		[]byte{0x23}, []byte{0xc2, 0xa4}, []byte{0xc3, 0x89}, []byte{0xc3, 0x84}, []byte{0xc3, 0x96},
		[]byte{0xc3, 0x85}, []byte{0xc3, 0x9c}, []byte{0x5f}, []byte{0xc3, 0xa9}, []byte{0xc3, 0xa4},
		[]byte{0xc3, 0xb6}, []byte{0xc3, 0xa5}, []byte{0xc3, 0xbc},
	}
	teletextNationalSubsetTurkish = &teletextNationalSubset{
		[]byte{0x54}, []byte{0xc4, 0x9f}, []byte{0xc4, 0xb0}, []byte{0xc5, 0x9e}, []byte{0xc3, 0x96},
		[]byte{0xc3, 0x87}, []byte{0xc3, 0x9c}, []byte{0xc4, 0x9e}, []byte{0xc4, 0xb1}, []byte{0xc5, 0x9f},
		[]byte{0xc3, 0xb6}, []byte{0xc3, 0xa7}, []byte{0xc3, 0xbc},
	}
)

// Teletext PES data types
const (
	teletextPESDataTypeEBU     = "EBU"
	teletextPESDataTypeUnknown = "unknown"
)

func teletextPESDataType(dataIdentifier uint8) string {
	switch {
	case dataIdentifier >= 0x10 && dataIdentifier <= 0x1f:
		return teletextPESDataTypeEBU
	}
	return teletextPESDataTypeUnknown
}

// Teletext PES data unit ids
const (
	teletextPESDataUnitIDEBUNonSubtitleData = 0x2
	teletextPESDataUnitIDEBUSubtitleData    = 0x3
	teletextPESDataUnitIDStuffing           = 0xff
)

// TeletextOptions represents teletext options
type TeletextOptions struct {
	Page int
	PID  int
}

// ReadFromTeletext parses a teletext content
// http://www.etsi.org/deliver/etsi_en/300400_300499/300472/01.03.01_60/en_300472v010301p.pdf
// http://www.etsi.org/deliver/etsi_i_ets/300700_300799/300706/01_60/ets_300706e01p.pdf
// TODO Update README
// TODO Add tests
func ReadFromTeletext(r io.Reader, o TeletextOptions) (s *Subtitles, err error) {
	// Init
	s = &Subtitles{}
	var dmx = astits.NewDemuxer(context.Background(), r)

	// Get the teletext PID
	var pid uint16
	if pid, err = teletextPID(dmx, o); err != nil {
		if err != ErrNoValidTeletextPID {
			err = fmt.Errorf("astisub: getting teletext PID failed: %w", err)
		}
		return
	}

	// Create character decoder
	cd := newTeletextCharacterDecoder()

	// Create page buffer
	b := newTeletextPageBuffer(o.Page, cd)

	// Loop in data
	var firstTime, lastTime time.Time
	var d *astits.DemuxerData
	var ps []*teletextPage
	for {
		// Fetch next data
		if d, err = dmx.NextData(); err != nil {
			if err == astits.ErrNoMorePackets {
				err = nil
				break
			}
			err = fmt.Errorf("astisub: fetching next data failed: %w", err)
			return
		}

		// We only parse PES data
		if d.PES == nil {
			continue
		}

		// This data is not of interest to us
		if d.PID != pid || d.PES.Header.StreamID != astits.StreamIDPrivateStream1 {
			continue
		}

		// Get time
		t := teletextDataTime(d)
		if t.IsZero() {
			continue
		}

		// First and last time
		if firstTime.IsZero() || firstTime.After(t) {
			firstTime = t
		}
		if lastTime.IsZero() || lastTime.Before(t) {
			lastTime = t
		}

		// Append pages
		ps = append(ps, b.process(d.PES, t)...)
	}

	// Dump buffer
	ps = append(ps, b.dump(lastTime)...)

	// Parse pages
	for _, p := range ps {
		p.parse(s, cd, firstTime)
	}
	return
}

// TODO Add tests
func teletextDataTime(d *astits.DemuxerData) time.Time {
	if d.PES.Header != nil && d.PES.Header.OptionalHeader != nil && d.PES.Header.OptionalHeader.PTS != nil {
		return d.PES.Header.OptionalHeader.PTS.Time()
	} else if d.FirstPacket != nil && d.FirstPacket.AdaptationField != nil && d.FirstPacket.AdaptationField.PCR != nil {
		return d.FirstPacket.AdaptationField.PCR.Time()
	}
	return time.Time{}
}

// If the PID teletext option is not indicated, it will walk through the ts data until it reaches a PMT packet to
// detect the first valid teletext PID
// TODO Add tests
func teletextPID(dmx *astits.Demuxer, o TeletextOptions) (pid uint16, err error) {
	// PID is in the options
	if o.PID > 0 {
		pid = uint16(o.PID)
		return
	}

	// Loop in data
	var d *astits.DemuxerData
	for {
		// Fetch next data
		if d, err = dmx.NextData(); err != nil {
			if err == astits.ErrNoMorePackets {
				err = ErrNoValidTeletextPID
				return
			}
			err = fmt.Errorf("astisub: fetching next data failed: %w", err)
			return
		}

		// PMT data
		if d.PMT != nil {
			// Retrieve valid teletext PIDs
			var pids []uint16
			for _, s := range d.PMT.ElementaryStreams {
				for _, dsc := range s.ElementaryStreamDescriptors {
					if dsc.Tag == astits.DescriptorTagTeletext || dsc.Tag == astits.DescriptorTagVBITeletext {
						pids = append(pids, s.ElementaryPID)
					}
				}
			}

			// No valid teletext PIDs
			if len(pids) == 0 {
				err = ErrNoValidTeletextPID
				return
			}

			// Set pid
			pid = pids[0]
			log.Printf("astisub: no teletext pid specified, using pid %d", pid)

			// Rewind
			if _, err = dmx.Rewind(); err != nil {
				err = fmt.Errorf("astisub: rewinding failed: %w", err)
				return
			}
			return
		}
	}
}

type teletextPageBuffer struct {
	cd             *teletextCharacterDecoder
	currentPage    *teletextPage
	donePages      []*teletextPage
	magazineNumber uint8
	pageNumber     int
	receiving      bool
}

func newTeletextPageBuffer(page int, cd *teletextCharacterDecoder) *teletextPageBuffer {
	return &teletextPageBuffer{
		cd:             cd,
		magazineNumber: uint8(page / 100),
		pageNumber:     page % 100,
	}
}

// TODO Add tests
func (b *teletextPageBuffer) dump(lastTime time.Time) (ps []*teletextPage) {
	if b.currentPage != nil {
		b.currentPage.end = lastTime
		ps = []*teletextPage{b.currentPage}
	}
	return
}

// TODO Add tests
func (b *teletextPageBuffer) process(d *astits.PESData, t time.Time) (ps []*teletextPage) {
	// Data identifier
	var offset int
	dataIdentifier := uint8(d.Data[offset])
	offset += 1

	// Check data type
	if teletextPESDataType(dataIdentifier) != teletextPESDataTypeEBU {
		return
	}

	// Loop through data units
	for offset < len(d.Data) {
		// ID
		id := uint8(d.Data[offset])
		offset += 1

		// Length
		length := uint8(d.Data[offset])
		offset += 1

		// Offset end
		offsetEnd := offset + int(length)
		if offsetEnd > len(d.Data) {
			break
		}

		// Parse data unit
		b.parseDataUnit(d.Data[offset:offsetEnd], id, t)

		// Seek to end of data unit
		offset = offsetEnd
	}

	// Dump buffer
	ps = b.donePages
	b.donePages = []*teletextPage(nil)
	return ps
}

// TODO Add tests
func (b *teletextPageBuffer) parseDataUnit(i []byte, id uint8, t time.Time) {
	// Check id
	if id != teletextPESDataUnitIDEBUSubtitleData {
		return
	}

	// Field parity: i[0]&0x20 > 0
	// Line offset: uint8(i[0] & 0x1f)
	// Framing code
	framingCode := uint8(i[1])

	// Check framing code
	if framingCode != 0xe4 {
		return
	}

	// Magazine number and packet number
	h1, ok := astikit.ByteHamming84Decode(i[2])
	if !ok {
		return
	}
	h2, ok := astikit.ByteHamming84Decode(i[3])
	if !ok {
		return
	}
	h := h2<<4 | h1
	magazineNumber := h & 0x7
	if magazineNumber == 0 {
		magazineNumber = 8
	}
	packetNumber := h >> 3

	// Parse packet
	b.parsePacket(i[4:], magazineNumber, packetNumber, t)
}

// TODO Add tests
func (b *teletextPageBuffer) parsePacket(i []byte, magazineNumber, packetNumber uint8, t time.Time) {
	if packetNumber == 0 {
		b.parsePacketHeader(i, magazineNumber, t)
	} else if b.receiving && magazineNumber == b.magazineNumber && (packetNumber >= 1 && packetNumber <= 25) {
		b.parsePacketData(i, packetNumber)
	} else {
		// Designation code
		designationCode, ok := astikit.ByteHamming84Decode(i[0])
		if !ok {
			return
		}

		// Parse packet
		if b.receiving && magazineNumber == b.magazineNumber && packetNumber == 26 {
			// TODO Implement
		} else if b.receiving && magazineNumber == b.magazineNumber && packetNumber == 28 {
			b.parsePacket28And29(i[1:], packetNumber, designationCode)
		} else if magazineNumber == b.magazineNumber && packetNumber == 29 {
			b.parsePacket28And29(i[1:], packetNumber, designationCode)
		} else if magazineNumber == 8 && packetNumber == 30 {
			b.parsePacket30(i, designationCode)
		}
	}
}

// TODO Add tests
func (b *teletextPageBuffer) parsePacketHeader(i []byte, magazineNumber uint8, t time.Time) (transmissionDone bool) {
	// Page number units
	pageNumberUnits, ok := astikit.ByteHamming84Decode(i[0])
	if !ok {
		return
	}

	// Page number tens
	pageNumberTens, ok := astikit.ByteHamming84Decode(i[1])
	if !ok {
		return
	}
	pageNumber := int(pageNumberTens)*10 + int(pageNumberUnits)

	// 0xff is a reserved page number value
	if pageNumberTens == 0xf && pageNumberUnits == 0xf {
		return
	}

	// Update magazine and page number
	if b.magazineNumber == 0 && b.pageNumber == 0 {
		// C6
		controlBits, ok := astikit.ByteHamming84Decode(i[5])
		if !ok {
			return
		}
		subtitleFlag := controlBits&0x8 > 0

		// This is a subtitle page
		if subtitleFlag {
			b.magazineNumber = magazineNumber
			b.pageNumber = pageNumber
			log.Printf("astisub: no teletext page specified, using page %d%.2d", b.magazineNumber, b.pageNumber)
		}
	}

	// C11 --> C14
	controlBits, ok := astikit.ByteHamming84Decode(i[7])
	if !ok {
		return
	}
	magazineSerial := controlBits&0x1 > 0
	charsetCode := controlBits >> 1

	// Page transmission is done
	if b.receiving && ((magazineSerial && pageNumber != b.pageNumber) ||
		(!magazineSerial && pageNumber != b.pageNumber && magazineNumber == b.magazineNumber)) {
		b.receiving = false
		return
	}

	// Invalid magazine or page number
	if pageNumber != b.pageNumber || magazineNumber != b.magazineNumber {
		return
	}

	// Now that we know when the previous page ends we can add it to the done slice
	if b.currentPage != nil {
		b.currentPage.end = t
		b.donePages = append(b.donePages, b.currentPage)
	}

	// Reset
	b.receiving = true
	b.currentPage = newTeletextPage(charsetCode, t)
	return
}

// TODO Add tests
func (b *teletextPageBuffer) parsePacketData(i []byte, packetNumber uint8) {
	// Make sure the map is initialized
	if _, ok := b.currentPage.data[packetNumber]; !ok {
		b.currentPage.data[packetNumber] = make([]byte, 40)
	}

	// Loop through input
	b.currentPage.rows = append(b.currentPage.rows, int(packetNumber))
	for idx := uint8(0); idx < 40; idx++ {
		v, ok := astikit.ByteParity(bits.Reverse8(i[idx]))
		if !ok {
			v = 0
		}
		b.currentPage.data[packetNumber][idx] = v
	}
}

// TODO Add tests
func (b *teletextPageBuffer) parsePacket28And29(i []byte, packetNumber, designationCode uint8) {
	// Invalid designation code
	if designationCode != 0 && designationCode != 4 {
		return
	}

	// Triplet 1
	// TODO triplet1 should be the results of hamming 24/18 decoding
	triplet1 := uint32(i[2])<<16 | uint32(i[1])<<8 | uint32(i[0])

	// We only process x/28 format 1
	if packetNumber == 28 && triplet1&0xf > 0 {
		return
	}

	// Update character decoder
	if packetNumber == 28 {
		b.cd.setTripletX28(triplet1)
	} else {
		b.cd.setTripletM29(triplet1)
	}
}

// TODO Add tests
func (b *teletextPageBuffer) parsePacket30(i []byte, designationCode uint8) {
	// Switch on designation code to determine format
	switch designationCode {
	case 0, 1:
		b.parsePacket30Format1(i)
	case 2, 3:
		b.parsePacket30Format2(i)
	}
}

func (b *teletextPageBuffer) parsePacket30Format1(i []byte) {
	// TODO Implement

}

func (b *teletextPageBuffer) parsePacket30Format2(i []byte) {
	// TODO Implement
}

type teletextCharacterDecoder struct {
	c                   teletextCharset
	lastPageCharsetCode *uint8
	tripletM29          *uint32
	tripletX28          *uint32
}

func newTeletextCharacterDecoder() *teletextCharacterDecoder {
	return &teletextCharacterDecoder{}
}

// TODO Add tests
func (d *teletextCharacterDecoder) setTripletM29(i uint32) {
	if *d.tripletM29 != i {
		d.tripletM29 = astikit.UInt32Ptr(i)
		d.updateCharset(d.lastPageCharsetCode, true)
	}
}

// TODO Add tests
func (d *teletextCharacterDecoder) setTripletX28(i uint32) {
	if *d.tripletX28 != i {
		d.tripletX28 = astikit.UInt32Ptr(i)
		d.updateCharset(d.lastPageCharsetCode, true)
	}
}

// TODO Add tests
func (d *teletextCharacterDecoder) decode(i byte) []byte {
	if i < 0x20 {
		return []byte{}
	}
	return d.c[i-0x20]
}

// TODO Add tests
func (d *teletextCharacterDecoder) updateCharset(pageCharsetCode *uint8, force bool) {
	// Charset is up to date
	if d.lastPageCharsetCode != nil && *pageCharsetCode == *d.lastPageCharsetCode && !force {
		return
	}
	d.lastPageCharsetCode = pageCharsetCode

	// Get triplet
	var triplet uint32
	if d.tripletX28 != nil {
		triplet = *d.tripletX28
	} else if d.tripletM29 != nil {
		triplet = *d.tripletM29
	}

	// Get charsets
	d.c = *teletextCharsetG0Latin
	var nationalOptionSubset *teletextNationalSubset
	if v1, ok := teletextCharsets[uint8((triplet&0x3f80)>>10)]; ok {
		if v2, ok := v1[*pageCharsetCode]; ok {
			d.c = *v2.g0
			nationalOptionSubset = v2.national
		}
	}

	// Update g0 with national option subset
	if nationalOptionSubset != nil {
		for k, v := range nationalOptionSubset {
			d.c[teletextNationalSubsetCharactersPositionInG0[k]] = v
		}
	}
}

type teletextPage struct {
	charsetCode uint8
	data        map[uint8][]byte
	end         time.Time
	rows        []int
	start       time.Time
}

func newTeletextPage(charsetCode uint8, start time.Time) *teletextPage {
	return &teletextPage{
		charsetCode: charsetCode,
		data:        make(map[uint8][]byte),
		start:       start,
	}
}

func (p *teletextPage) parse(s *Subtitles, d *teletextCharacterDecoder, firstTime time.Time) {
	// Update charset
	d.updateCharset(astikit.UInt8Ptr(p.charsetCode), false)

	// No data
	if len(p.data) == 0 {
		return
	}

	// Order rows
	sort.Ints(p.rows)

	// Create item
	i := &Item{
		EndAt:   p.end.Sub(firstTime),
		StartAt: p.start.Sub(firstTime),
	}

	// Loop through rows
	for _, idxRow := range p.rows {
		parseTeletextRow(i, d, nil, p.data[uint8(idxRow)])
	}

	// Append item
	s.Items = append(s.Items, i)
}

type decoder interface {
	decode(i byte) []byte
}

type styler interface {
	hasBeenSet() bool
	hasChanged(s *StyleAttributes) bool
	parseSpacingAttribute(i byte)
	propagateStyleAttributes(s *StyleAttributes)
	update(sa *StyleAttributes)
}

func parseTeletextRow(i *Item, d decoder, fs func() styler, row []byte) {
	// Loop through columns
	var l = Line{}
	var li = LineItem{InlineStyle: &StyleAttributes{}}
	var started bool
	var s styler
	for _, v := range row {
		// Create specific styler
		if fs != nil {
			s = fs()
		}

		// Get spacing attributes
		var color *Color
		var doubleHeight, doubleSize, doubleWidth *bool
		switch v {
		case 0x0:
			color = ColorBlack
		case 0x1:
			color = ColorRed
		case 0x2:
			color = ColorGreen
		case 0x3:
			color = ColorYellow
		case 0x4:
			color = ColorBlue
		case 0x5:
			color = ColorMagenta
		case 0x6:
			color = ColorCyan
		case 0x7:
			color = ColorWhite
		case 0xa:
			started = false
		case 0xb:
			started = true
		case 0xc:
			doubleHeight = astikit.BoolPtr(false)
			doubleSize = astikit.BoolPtr(false)
			doubleWidth = astikit.BoolPtr(false)
		case 0xd:
			doubleHeight = astikit.BoolPtr(true)
		case 0xe:
			doubleWidth = astikit.BoolPtr(true)
		case 0xf:
			doubleSize = astikit.BoolPtr(true)
		default:
			if s != nil {
				s.parseSpacingAttribute(v)
			}
		}

		// Style has been set
		if color != nil || doubleHeight != nil || doubleSize != nil || doubleWidth != nil || (s != nil && s.hasBeenSet()) {
			// Style has changed
			if color != li.InlineStyle.TeletextColor || doubleHeight != li.InlineStyle.TeletextDoubleHeight ||
				doubleSize != li.InlineStyle.TeletextDoubleSize || doubleWidth != li.InlineStyle.TeletextDoubleWidth ||
				(s != nil && s.hasChanged(li.InlineStyle)) {
				// Line has started
				if started {
					// Append line item
					appendTeletextLineItem(&l, li, s)

					// Create new line item
					sa := &StyleAttributes{}
					*sa = *li.InlineStyle
					li = LineItem{InlineStyle: sa}
				}

				// Update style attributes
				if color != nil && color != li.InlineStyle.TeletextColor {
					li.InlineStyle.TeletextColor = color
				}
				if doubleHeight != nil && doubleHeight != li.InlineStyle.TeletextDoubleHeight {
					li.InlineStyle.TeletextDoubleHeight = doubleHeight
				}
				if doubleSize != nil && doubleSize != li.InlineStyle.TeletextDoubleSize {
					li.InlineStyle.TeletextDoubleSize = doubleSize
				}
				if doubleWidth != nil && doubleWidth != li.InlineStyle.TeletextDoubleWidth {
					li.InlineStyle.TeletextDoubleWidth = doubleWidth
				}
				if s != nil {
					s.update(li.InlineStyle)
				}
			}
		} else if started {
			// Append text
			li.Text += string(d.decode(v))
		}
	}

	// Append line item
	appendTeletextLineItem(&l, li, s)

	// Append line
	if len(l.Items) > 0 {
		i.Lines = append(i.Lines, l)
	}
}

func appendTeletextLineItem(l *Line, li LineItem, s styler) {
	// There's some text
	if len(strings.TrimSpace(li.Text)) > 0 {
		// Make sure inline style exists
		if li.InlineStyle == nil {
			li.InlineStyle = &StyleAttributes{}
		}

		// Get number of spaces before
		li.InlineStyle.TeletextSpacesBefore = astikit.IntPtr(0)
		for _, c := range li.Text {
			if c == ' ' {
				*li.InlineStyle.TeletextSpacesBefore++
			} else {
				break
			}
		}

		// Get number of spaces after
		li.InlineStyle.TeletextSpacesAfter = astikit.IntPtr(0)
		for idx := len(li.Text) - 1; idx >= 0; idx-- {
			if li.Text[idx] == ' ' {
				*li.InlineStyle.TeletextSpacesAfter++
			} else {
				break
			}
		}

		// Propagate style attributes
		li.InlineStyle.propagateTeletextAttributes()
		if s != nil {
			s.propagateStyleAttributes(li.InlineStyle)
		}

		// Append line item
		li.Text = strings.TrimSpace(li.Text)
		l.Items = append(l.Items, li)
	}
}

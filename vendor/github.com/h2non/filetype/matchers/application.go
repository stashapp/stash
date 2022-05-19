package matchers

var (
	TypeWasm = newType("wasm", "application/wasm")
	TypeDex  = newType("dex", "application/vnd.android.dex")
	TypeDey  = newType("dey", "application/vnd.android.dey")
)

var Application = Map{
	TypeWasm: Wasm,
	TypeDex:  Dex,
	TypeDey:  Dey,
}

// Wasm detects a Web Assembly 1.0 filetype.
func Wasm(buf []byte) bool {
	// WASM has starts with `\0asm`, followed by the version.
	// http://webassembly.github.io/spec/core/binary/modules.html#binary-magic
	return len(buf) >= 8 &&
		buf[0] == 0x00 && buf[1] == 0x61 &&
		buf[2] == 0x73 && buf[3] == 0x6D &&
		buf[4] == 0x01 && buf[5] == 0x00 &&
		buf[6] == 0x00 && buf[7] == 0x00
}

// Dex detects dalvik executable(DEX)
func Dex(buf []byte) bool {
	// https://source.android.com/devices/tech/dalvik/dex-format#dex-file-magic
	return len(buf) > 36 &&
		// magic
		buf[0] == 0x64 && buf[1] == 0x65 && buf[2] == 0x78 && buf[3] == 0x0A &&
		// file sise
		buf[36] == 0x70
}

// Dey Optimized Dalvik Executable(ODEX)
func Dey(buf []byte) bool {
	return len(buf) > 100 &&
		// dey magic
		buf[0] == 0x64 && buf[1] == 0x65 && buf[2] == 0x79 && buf[3] == 0x0A &&
		// dex
		Dex(buf[40:100])
}

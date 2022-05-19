package matchers

import "encoding/binary"

const (
	ZstdMagicSkippableStart = 0x184D2A50
	ZstdMagicSkippableMask  = 0xFFFFFFF0
)

var (
	TypeEpub   = newType("epub", "application/epub+zip")
	TypeZip    = newType("zip", "application/zip")
	TypeTar    = newType("tar", "application/x-tar")
	TypeRar    = newType("rar", "application/vnd.rar")
	TypeGz     = newType("gz", "application/gzip")
	TypeBz2    = newType("bz2", "application/x-bzip2")
	Type7z     = newType("7z", "application/x-7z-compressed")
	TypeXz     = newType("xz", "application/x-xz")
	TypeZstd   = newType("zst", "application/zstd")
	TypePdf    = newType("pdf", "application/pdf")
	TypeExe    = newType("exe", "application/vnd.microsoft.portable-executable")
	TypeSwf    = newType("swf", "application/x-shockwave-flash")
	TypeRtf    = newType("rtf", "application/rtf")
	TypeEot    = newType("eot", "application/octet-stream")
	TypePs     = newType("ps", "application/postscript")
	TypeSqlite = newType("sqlite", "application/vnd.sqlite3")
	TypeNes    = newType("nes", "application/x-nintendo-nes-rom")
	TypeCrx    = newType("crx", "application/x-google-chrome-extension")
	TypeCab    = newType("cab", "application/vnd.ms-cab-compressed")
	TypeDeb    = newType("deb", "application/vnd.debian.binary-package")
	TypeAr     = newType("ar", "application/x-unix-archive")
	TypeZ      = newType("Z", "application/x-compress")
	TypeLz     = newType("lz", "application/x-lzip")
	TypeRpm    = newType("rpm", "application/x-rpm")
	TypeElf    = newType("elf", "application/x-executable")
	TypeDcm    = newType("dcm", "application/dicom")
	TypeIso    = newType("iso", "application/x-iso9660-image")
	TypeMachO  = newType("macho", "application/x-mach-binary") // Mach-O binaries have no common extension.
)

var Archive = Map{
	TypeEpub:   bytePrefixMatcher(epubMagic),
	TypeZip:    Zip,
	TypeTar:    Tar,
	TypeRar:    Rar,
	TypeGz:     bytePrefixMatcher(gzMagic),
	TypeBz2:    bytePrefixMatcher(bz2Magic),
	Type7z:     bytePrefixMatcher(sevenzMagic),
	TypeXz:     bytePrefixMatcher(xzMagic),
	TypeZstd:   Zst,
	TypePdf:    bytePrefixMatcher(pdfMagic),
	TypeExe:    bytePrefixMatcher(exeMagic),
	TypeSwf:    Swf,
	TypeRtf:    bytePrefixMatcher(rtfMagic),
	TypeEot:    Eot,
	TypePs:     bytePrefixMatcher(psMagic),
	TypeSqlite: bytePrefixMatcher(sqliteMagic),
	TypeNes:    bytePrefixMatcher(nesMagic),
	TypeCrx:    bytePrefixMatcher(crxMagic),
	TypeCab:    Cab,
	TypeDeb:    bytePrefixMatcher(debMagic),
	TypeAr:     bytePrefixMatcher(arMagic),
	TypeZ:      Z,
	TypeLz:     bytePrefixMatcher(lzMagic),
	TypeRpm:    Rpm,
	TypeElf:    Elf,
	TypeDcm:    Dcm,
	TypeIso:    Iso,
	TypeMachO:  MachO,
}

var (
	epubMagic = []byte{
		0x50, 0x4B, 0x03, 0x04, 0x6D, 0x69, 0x6D, 0x65,
		0x74, 0x79, 0x70, 0x65, 0x61, 0x70, 0x70, 0x6C,
		0x69, 0x63, 0x61, 0x74, 0x69, 0x6F, 0x6E, 0x2F,
		0x65, 0x70, 0x75, 0x62, 0x2B, 0x7A, 0x69, 0x70,
	}
	gzMagic     = []byte{0x1F, 0x8B, 0x08}
	bz2Magic    = []byte{0x42, 0x5A, 0x68}
	sevenzMagic = []byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C}
	pdfMagic    = []byte{0x25, 0x50, 0x44, 0x46}
	exeMagic    = []byte{0x4D, 0x5A}
	rtfMagic    = []byte{0x7B, 0x5C, 0x72, 0x74, 0x66}
	nesMagic    = []byte{0x4E, 0x45, 0x53, 0x1A}
	crxMagic    = []byte{0x43, 0x72, 0x32, 0x34}
	psMagic     = []byte{0x25, 0x21}
	xzMagic     = []byte{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00}
	sqliteMagic = []byte{0x53, 0x51, 0x4C, 0x69}
	debMagic    = []byte{
		0x21, 0x3C, 0x61, 0x72, 0x63, 0x68, 0x3E, 0x0A,
		0x64, 0x65, 0x62, 0x69, 0x61, 0x6E, 0x2D, 0x62,
		0x69, 0x6E, 0x61, 0x72, 0x79,
	}
	arMagic   = []byte{0x21, 0x3C, 0x61, 0x72, 0x63, 0x68, 0x3E}
	zstdMagic = []byte{0x28, 0xB5, 0x2F, 0xFD}
	lzMagic   = []byte{0x4C, 0x5A, 0x49, 0x50}
)

func bytePrefixMatcher(magicPattern []byte) Matcher {
	return func(data []byte) bool {
		return compareBytes(data, magicPattern, 0)
	}
}

func Zip(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x50 && buf[1] == 0x4B &&
		(buf[2] == 0x3 || buf[2] == 0x5 || buf[2] == 0x7) &&
		(buf[3] == 0x4 || buf[3] == 0x6 || buf[3] == 0x8)
}

func Tar(buf []byte) bool {
	return len(buf) > 261 &&
		buf[257] == 0x75 && buf[258] == 0x73 &&
		buf[259] == 0x74 && buf[260] == 0x61 &&
		buf[261] == 0x72
}

func Rar(buf []byte) bool {
	return len(buf) > 6 &&
		buf[0] == 0x52 && buf[1] == 0x61 && buf[2] == 0x72 &&
		buf[3] == 0x21 && buf[4] == 0x1A && buf[5] == 0x7 &&
		(buf[6] == 0x0 || buf[6] == 0x1)
}

func Swf(buf []byte) bool {
	return len(buf) > 2 &&
		(buf[0] == 0x43 || buf[0] == 0x46) &&
		buf[1] == 0x57 && buf[2] == 0x53
}

func Cab(buf []byte) bool {
	return len(buf) > 3 &&
		((buf[0] == 0x4D && buf[1] == 0x53 && buf[2] == 0x43 && buf[3] == 0x46) ||
			(buf[0] == 0x49 && buf[1] == 0x53 && buf[2] == 0x63 && buf[3] == 0x28))
}

func Eot(buf []byte) bool {
	return len(buf) > 35 &&
		buf[34] == 0x4C && buf[35] == 0x50 &&
		((buf[8] == 0x02 && buf[9] == 0x00 &&
			buf[10] == 0x01) || (buf[8] == 0x01 &&
			buf[9] == 0x00 && buf[10] == 0x00) ||
			(buf[8] == 0x02 && buf[9] == 0x00 &&
				buf[10] == 0x02))
}

func Z(buf []byte) bool {
	return len(buf) > 1 &&
		((buf[0] == 0x1F && buf[1] == 0xA0) ||
			(buf[0] == 0x1F && buf[1] == 0x9D))
}

func Rpm(buf []byte) bool {
	return len(buf) > 96 &&
		buf[0] == 0xED && buf[1] == 0xAB &&
		buf[2] == 0xEE && buf[3] == 0xDB
}

func Elf(buf []byte) bool {
	return len(buf) > 52 &&
		buf[0] == 0x7F && buf[1] == 0x45 &&
		buf[2] == 0x4C && buf[3] == 0x46
}

func Dcm(buf []byte) bool {
	return len(buf) > 131 &&
		buf[128] == 0x44 && buf[129] == 0x49 &&
		buf[130] == 0x43 && buf[131] == 0x4D
}

func Iso(buf []byte) bool {
	return len(buf) > 32773 &&
		buf[32769] == 0x43 && buf[32770] == 0x44 &&
		buf[32771] == 0x30 && buf[32772] == 0x30 &&
		buf[32773] == 0x31
}

func MachO(buf []byte) bool {
	return len(buf) > 3 && ((buf[0] == 0xFE && buf[1] == 0xED && buf[2] == 0xFA && buf[3] == 0xCF) ||
		(buf[0] == 0xFE && buf[1] == 0xED && buf[2] == 0xFA && buf[3] == 0xCE) ||
		(buf[0] == 0xBE && buf[1] == 0xBA && buf[2] == 0xFE && buf[3] == 0xCA) ||
		// Big endian versions below here...
		(buf[0] == 0xCF && buf[1] == 0xFA && buf[2] == 0xED && buf[3] == 0xFE) ||
		(buf[0] == 0xCE && buf[1] == 0xFA && buf[2] == 0xED && buf[3] == 0xFE) ||
		(buf[0] == 0xCA && buf[1] == 0xFE && buf[2] == 0xBA && buf[3] == 0xBE))
}

// Zstandard compressed data is made of one or more frames.
// There are two frame formats defined by Zstandard: Zstandard frames and Skippable frames.
// See more details from https://tools.ietf.org/id/draft-kucherawy-dispatch-zstd-00.html#rfc.section.2
func Zst(buf []byte) bool {
  if compareBytes(buf, zstdMagic, 0) {
    return true
  } else {
		// skippable frames
    if len(buf) < 8 {
      return false
    }
    if binary.LittleEndian.Uint32(buf[:4]) & ZstdMagicSkippableMask == ZstdMagicSkippableStart {
      userDataLength := binary.LittleEndian.Uint32(buf[4:8])
      if len(buf) < 8 + int(userDataLength) {
        return false
      }
      nextFrame := buf[8+userDataLength:]
      return Zst(nextFrame)
    }
    return false
  }
}

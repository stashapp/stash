package w32

// combase!_SECTION_IMAGE_INFORMATION
//    +0x000 TransferAddress  : Ptr64 Void
//    +0x008 ZeroBits         : Uint4B
//    +0x010 MaximumStackSize : Uint8B
//    +0x018 CommittedStackSize : Uint8B
//    +0x020 SubSystemType    : Uint4B
//    +0x024 SubSystemMinorVersion : Uint2B
//    +0x026 SubSystemMajorVersion : Uint2B
//    +0x024 SubSystemVersion : Uint4B
//    +0x028 MajorOperatingSystemVersion : Uint2B
//    +0x02a MinorOperatingSystemVersion : Uint2B
//    +0x028 OperatingSystemVersion : Uint4B
//    +0x02c ImageCharacteristics : Uint2B
//    +0x02e DllCharacteristics : Uint2B
//    +0x030 Machine          : Uint2B
//    +0x032 ImageContainsCode : UChar
//    +0x033 ImageFlags       : UChar
//    +0x033 ComPlusNativeReady : Pos 0, 1 Bit
//    +0x033 ComPlusILOnly    : Pos 1, 1 Bit
//    +0x033 ImageDynamicallyRelocated : Pos 2, 1 Bit
//    +0x033 ImageMappedFlat  : Pos 3, 1 Bit
//    +0x033 BaseBelow4gb     : Pos 4, 1 Bit
//    +0x033 ComPlusPrefer32bit : Pos 5, 1 Bit
//    +0x033 Reserved         : Pos 6, 2 Bits
//    +0x034 LoaderFlags      : Uint4B
//    +0x038 ImageFileSize    : Uint4B
//    +0x03c CheckSum         : Uint4B
type SECTION_IMAGE_INFORMATION struct {
	TransferAddress             uintptr
	ZeroBits                    uint32
	MaximumStackSize            uint64
	CommittedStackSize          uint64
	SubSystemType               uint32
	SubSystemMinorVersion       uint16
	SubSystemMajorVersion       uint16
	SubSystemVersion            uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	OperatingSystemVersion      uint32
	ImageCharacteristics        uint16
	DllCharacteristics          uint16
	Machine                     uint16
	ImageContainsCode           uint8
	ImageFlags                  uint8
	ComPlusFlags                uint8
	LoaderFlags                 uint32
	ImageFileSize               uint32
	CheckSum                    uint32
}

func (si *SECTION_IMAGE_INFORMATION) ComPlusNativeReady() bool {
	return (si.ComPlusFlags & (1 << 0)) == 1
}

func (si *SECTION_IMAGE_INFORMATION) ComPlusILOnly() bool {
	return (si.ComPlusFlags & (1 << 1)) == 1
}

func (si *SECTION_IMAGE_INFORMATION) ImageDynamicallyRelocated() bool {
	return (si.ComPlusFlags & (1 << 2)) == 1
}

func (si *SECTION_IMAGE_INFORMATION) ImageMappedFlat() bool {
	return (si.ComPlusFlags & (1 << 3)) == 1
}

func (si *SECTION_IMAGE_INFORMATION) BaseBelow4gb() bool {
	return (si.ComPlusFlags & (1 << 4)) == 1
}

func (si *SECTION_IMAGE_INFORMATION) ComPlusPrefer32bit() bool {
	return (si.ComPlusFlags & (1 << 5)) == 1
}

// combase!_RTL_USER_PROCESS_INFORMATION
//    +0x000 Length           : Uint4B
//    +0x008 Process          : Ptr64 Void
//    +0x010 Thread           : Ptr64 Void
//    +0x018 ClientId         : _CLIENT_ID
//    +0x028 ImageInformation : _SECTION_IMAGE_INFORMATION
type RTL_USER_PROCESS_INFORMATION struct {
	Length           uint32
	Process          HANDLE
	Thread           HANDLE
	ClientId         CLIENT_ID
	ImageInformation SECTION_IMAGE_INFORMATION
}

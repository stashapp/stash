package w32

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa374931(v=vs.85).aspx
type ACL struct {
	AclRevision byte
	Sbz1        byte
	AclSize     uint16
	AceCount    uint16
	Sbz2        uint16
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa379561(v=vs.85).aspx

type SECURITY_DESCRIPTOR_CONTROL uint16

type SECURITY_DESCRIPTOR struct {
	Revision byte
	Sbz1     byte
	Control  SECURITY_DESCRIPTOR_CONTROL
	Owner    uintptr
	Group    uintptr
	Sacl     *ACL
	Dacl     *ACL
}

type SID_IDENTIFIER_AUTHORITY struct {
	Value [6]byte
}

// typedef struct _SID // 4 elements, 0xC bytes (sizeof)
// {
// /*0x000*/     UINT8        Revision;
// /*0x001*/     UINT8        SubAuthorityCount;
// /*0x002*/     struct _SID_IDENTIFIER_AUTHORITY IdentifierAuthority; // 1 elements, 0x6 bytes (sizeof)
// /*0x008*/     ULONG32      SubAuthority[1];
// }SID, *PSID;
type SID struct {
	Revision            byte
	SubAuthorityCount   byte
	IdentifierAuthority SID_IDENTIFIER_AUTHORITY
	SubAuthority        uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa363646.aspx
type EVENTLOGRECORD struct {
	Length              uint32
	Reserved            uint32
	RecordNumber        uint32
	TimeGenerated       uint32
	TimeWritten         uint32
	EventID             uint32
	EventType           uint16
	NumStrings          uint16
	EventCategory       uint16
	ReservedFlags       uint16
	ClosingRecordNumber uint32
	StringOffset        uint32
	UserSidLength       uint32
	UserSidOffset       uint32
	DataLength          uint32
	DataOffset          uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms685996.aspx
type SERVICE_STATUS struct {
	DwServiceType             uint32
	DwCurrentState            uint32
	DwControlsAccepted        uint32
	DwWin32ExitCode           uint32
	DwServiceSpecificExitCode uint32
	DwCheckPoint              uint32
	DwWaitHint                uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa364160(v=vs.85).aspx
type WNODE_HEADER struct {
	BufferSize        uint32
	ProviderId        uint32
	HistoricalContext uint64
	KernelHandle      HANDLE
	Guid              GUID
	ClientContext     uint32
	Flags             uint32
}

// These partially compensate for the anonymous unions we removed, but there
// are no setters.
func (w WNODE_HEADER) TimeStamp() uint64 {
	// TODO: Cast to the stupid LARGE_INTEGER struct which is, itself, nasty
	// and union-y
	return uint64(w.KernelHandle)
}

func (w WNODE_HEADER) Version() uint32 {
	return uint32(w.HistoricalContext >> 32)
}

func (w WNODE_HEADER) Linkage() uint32 {
	return uint32(w.HistoricalContext)
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa363784(v=vs.85).aspx
type EVENT_TRACE_PROPERTIES struct {
	Wnode               WNODE_HEADER
	BufferSize          uint32
	MinimumBuffers      uint32
	MaximumBuffers      uint32
	MaximumFileSize     uint32
	LogFileMode         uint32
	FlushTimer          uint32
	EnableFlags         uint32
	AgeLimit            int32
	NumberOfBuffers     uint32
	FreeBuffers         uint32
	EventsLost          uint32
	BuffersWritten      uint32
	LogBuffersLost      uint32
	RealTimeBuffersLost uint32
	LoggerThreadId      HANDLE
	LogFileNameOffset   uint32
	LoggerNameOffset    uint32
}

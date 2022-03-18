package w32

// typedef struct _PROCESS_INFORMATION {
//   HANDLE hProcess;
//   HANDLE hThread;
//   DWORD dwProcessId;
//   DWORD dwThreadId;
// } PROCESS_INFORMATION, *PPROCESS_INFORMATION, *LPPROCESS_INFORMATION;

type PROCESS_INFORMATION struct {
	Process   HANDLE
	Thread    HANDLE
	ProcessId uint32
	ThreadId  uint32
}

// typedef struct _STARTUPINFOW {
//   DWORD cb;
//   LPWSTR lpReserved;
//   LPWSTR lpDesktop;
//   LPWSTR lpTitle;
//   DWORD dwX;
//   DWORD dwY;
//   DWORD dwXSize;
//   DWORD dwYSize;
//   DWORD dwXCountChars;
//   DWORD dwYCountChars;
//   DWORD dwFillAttribute;
//   DWORD dwFlags;
//   WORD wShowWindow;
//   WORD cbReserved2;
//   LPBYTE lpReserved2;
//   HANDLE hStdInput;
//   HANDLE hStdOutput;
//   HANDLE hStdError;
// } STARTUPINFOW, *LPSTARTUPINFOW;

type STARTUPINFOW struct {
	cb            uint32
	_             *uint16
	Desktop       *uint16
	Title         *uint16
	X             uint32
	Y             uint32
	XSize         uint32
	YSize         uint32
	XCountChars   uint32
	YCountChars   uint32
	FillAttribute uint32
	Flags         uint32
	ShowWindow    uint16
	_             uint16
	_             *uint8
	StdInput      HANDLE
	StdOutput     HANDLE
	StdError      HANDLE
}

// combase!_SECURITY_ATTRIBUTES
//    +0x000 nLength          : Uint4B
//    +0x008 lpSecurityDescriptor : Ptr64 Void
//    +0x010 bInheritHandle   : Int4B

type SECURITY_ATTRIBUTES struct {
	Length             uint32
	SecurityDescriptor uintptr
	InheritHandle      BOOL
}

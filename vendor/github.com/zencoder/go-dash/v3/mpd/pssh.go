package mpd

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func MakePSSHBox(systemID, payload []byte) ([]byte, error) {
	if len(systemID) != 16 {
		return nil, fmt.Errorf("SystemID must be 16 bytes, was: %d", len(systemID))
	}

	psshBuf := &bytes.Buffer{}
	size := uint32(12 + 16 + 4 + len(payload)) // 3 uint32s, systemID, "pssh" string and payload
	if err := binary.Write(psshBuf, binary.BigEndian, size); err != nil {
		return nil, err
	}

	if err := binary.Write(psshBuf, binary.BigEndian, []byte("pssh")); err != nil {
		return nil, err
	}

	if err := binary.Write(psshBuf, binary.BigEndian, uint32(0)); err != nil {
		return nil, err
	}

	if _, err := psshBuf.Write(systemID); err != nil {
		return nil, err
	}

	if err := binary.Write(psshBuf, binary.BigEndian, uint32(len(payload))); err != nil {
		return nil, err
	}

	if _, err := psshBuf.Write(payload); err != nil {
		return nil, err
	}

	return psshBuf.Bytes(), nil
}

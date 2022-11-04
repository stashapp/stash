package astikit

func BoolToUInt32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

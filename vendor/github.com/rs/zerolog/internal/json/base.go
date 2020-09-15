package json

type Encoder struct{}

// AppendKey appends a new key to the output JSON.
func (e Encoder) AppendKey(dst []byte, key string) []byte {
	if len(dst) > 1 && dst[len(dst)-1] != '{' {
		dst = append(dst, ',')
	}
	dst = e.AppendString(dst, key)
	return append(dst, ':')
}
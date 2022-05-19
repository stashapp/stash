package otherpkg

type (
	Scalar string
	Map    map[string]string
	Slice  []string
)

type Struct struct {
	Name Scalar
	Desc *Scalar
}

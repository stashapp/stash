package astikit

import (
	"os"
	"strings"
)

// FlagCmd retrieves the command from the input Args
func FlagCmd() (o string) {
	if len(os.Args) >= 2 && os.Args[1][0] != '-' {
		o = os.Args[1]
		os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
	}
	return
}

// FlagStrings represents a flag that can be set several times and
// stores unique string values
type FlagStrings struct {
	Map   map[string]bool
	Slice *[]string
}

// NewFlagStrings creates a new FlagStrings
func NewFlagStrings() FlagStrings {
	return FlagStrings{
		Map:   make(map[string]bool),
		Slice: &[]string{},
	}
}

// String implements the flag.Value interface
func (f FlagStrings) String() string {
	if f.Slice == nil {
		return ""
	}
	return strings.Join(*f.Slice, ",")
}

// Set implements the flag.Value interface
func (f FlagStrings) Set(i string) error {
	if _, ok := f.Map[i]; ok {
		return nil
	}
	f.Map[i] = true
	*f.Slice = append(*f.Slice, i)
	return nil
}

package models

import (
	"fmt"
	"io"
	"strconv"
)

type ImportMissingRefEnum string

const (
	ImportMissingRefEnumIgnore ImportMissingRefEnum = "IGNORE"
	ImportMissingRefEnumFail   ImportMissingRefEnum = "FAIL"
	ImportMissingRefEnumCreate ImportMissingRefEnum = "CREATE"
)

var AllImportMissingRefEnum = []ImportMissingRefEnum{
	ImportMissingRefEnumIgnore,
	ImportMissingRefEnumFail,
	ImportMissingRefEnumCreate,
}

func (e ImportMissingRefEnum) IsValid() bool {
	switch e {
	case ImportMissingRefEnumIgnore, ImportMissingRefEnumFail, ImportMissingRefEnumCreate:
		return true
	}
	return false
}

func (e ImportMissingRefEnum) String() string {
	return string(e)
}

func (e *ImportMissingRefEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ImportMissingRefEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ImportMissingRefEnum", str)
	}
	return nil
}

func (e ImportMissingRefEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

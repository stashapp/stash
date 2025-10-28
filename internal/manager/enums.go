package manager

import (
	"fmt"
	"io"
	"strconv"
)

type SystemStatusEnum string

const (
	SystemStatusEnumSetup          SystemStatusEnum = "SETUP"
	SystemStatusEnumNeedsMigration SystemStatusEnum = "NEEDS_MIGRATION"
	SystemStatusEnumOk             SystemStatusEnum = "OK"
)

var AllSystemStatusEnum = []SystemStatusEnum{
	SystemStatusEnumSetup,
	SystemStatusEnumNeedsMigration,
	SystemStatusEnumOk,
}

func (e SystemStatusEnum) IsValid() bool {
	switch e {
	case SystemStatusEnumSetup, SystemStatusEnumNeedsMigration, SystemStatusEnumOk:
		return true
	}
	return false
}

func (e SystemStatusEnum) String() string {
	return string(e)
}

func (e *SystemStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SystemStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SystemStatusEnum", str)
	}
	return nil
}

func (e SystemStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

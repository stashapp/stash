package plugin

import (
	"fmt"
	"io"
	"strconv"
)

type PluginSettingTypeEnum string

const (
	PluginSettingTypeEnumString  PluginSettingTypeEnum = "STRING"
	PluginSettingTypeEnumNumber  PluginSettingTypeEnum = "NUMBER"
	PluginSettingTypeEnumBoolean PluginSettingTypeEnum = "BOOLEAN"
)

var AllPluginSettingTypeEnum = []PluginSettingTypeEnum{
	PluginSettingTypeEnumString,
	PluginSettingTypeEnumNumber,
	PluginSettingTypeEnumBoolean,
}

func (e PluginSettingTypeEnum) IsValid() bool {
	switch e {
	case PluginSettingTypeEnumString, PluginSettingTypeEnumNumber, PluginSettingTypeEnumBoolean:
		return true
	}
	return false
}

func (e PluginSettingTypeEnum) String() string {
	return string(e)
}

func (e *PluginSettingTypeEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PluginSettingTypeEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PluginSettingTypeEnum", str)
	}
	return nil
}

func (e PluginSettingTypeEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

package utils

import "reflect"

// NotNilFields returns the matching tag values of fields from an object that are not nil.
// Panics if the provided object is not a struct.
func NotNilFields(subject interface{}, tag string) []string {
	value := reflect.ValueOf(subject)
	structType := value.Type()

	if structType.Kind() != reflect.Struct {
		panic("subject must be struct")
	}

	var ret []string

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		kind := field.Type().Kind()
		if (kind == reflect.Ptr || kind == reflect.Slice) && !field.IsNil() {
			tagValue := structType.Field(i).Tag.Get(tag)
			if tagValue != "" {
				ret = append(ret, tagValue)
			}
		}
	}

	return ret
}

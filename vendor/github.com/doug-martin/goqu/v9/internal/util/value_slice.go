package util

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type ValueSlice []reflect.Value

func (vs ValueSlice) Len() int           { return len(vs) }
func (vs ValueSlice) Less(i, j int) bool { return vs[i].String() < vs[j].String() }
func (vs ValueSlice) Swap(i, j int)      { vs[i], vs[j] = vs[j], vs[i] }

func (vs ValueSlice) Equal(other ValueSlice) bool {
	sort.Sort(other)
	for i, key := range vs {
		if other[i].String() != key.String() {
			return false
		}
	}
	return true
}

func (vs ValueSlice) String() string {
	vals := make([]string, vs.Len())
	for i, key := range vs {
		vals[i] = fmt.Sprintf(`"%s"`, key.String())
	}
	sort.Strings(vals)
	return fmt.Sprintf("[%s]", strings.Join(vals, ","))
}

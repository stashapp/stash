package codegen

import (
	"fmt"
	"go/types"
	"strings"
)

func findGoNamedType(def types.Type) (*types.Named, error) {
	if def == nil {
		return nil, nil
	}

	namedType, ok := def.(*types.Named)
	if !ok {
		return nil, fmt.Errorf("expected %s to be a named type, instead found %T\n", def.String(), def)
	}

	return namedType, nil
}

func findGoInterface(def types.Type) (*types.Interface, error) {
	if def == nil {
		return nil, nil
	}
	namedType, err := findGoNamedType(def)
	if err != nil {
		return nil, err
	}
	if namedType == nil {
		return nil, nil
	}

	underlying, ok := namedType.Underlying().(*types.Interface)
	if !ok {
		return nil, fmt.Errorf("expected %s to be a named interface, instead found %s", def.String(), namedType.String())
	}

	return underlying, nil
}

func equalFieldName(source, target string) bool {
	source = strings.ReplaceAll(source, "_", "")
	source = strings.ReplaceAll(source, ",omitempty", "")
	target = strings.ReplaceAll(target, "_", "")
	return strings.EqualFold(source, target)
}

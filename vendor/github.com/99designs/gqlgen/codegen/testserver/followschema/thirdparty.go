package followschema

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

type ThirdParty struct {
	str string
}

func MarshalThirdParty(tp ThirdParty) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, err := io.WriteString(w, strconv.Quote(tp.str))
		if err != nil {
			panic(err)
		}
	})
}

func UnmarshalThirdParty(input interface{}) (ThirdParty, error) {
	switch input := input.(type) {
	case string:
		return ThirdParty{str: input}, nil
	default:
		return ThirdParty{}, fmt.Errorf("unknown type for input: %s", input)
	}
}

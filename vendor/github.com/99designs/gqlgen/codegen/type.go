package codegen

import (
	"fmt"

	"github.com/99designs/gqlgen/codegen/config"
)

func (b *builder) buildTypes() map[string]*config.TypeReference {
	ret := map[string]*config.TypeReference{}
	var key string
	var existing *config.TypeReference
	var found bool

	for _, ref := range b.Binder.References {
		for ref != nil {
			key = ref.UniquenessKey()
			if existing, found = ret[key]; found {
				// Simplistic check of content which is obviously different.
				existingGQL := fmt.Sprintf("%v", existing.GQL)
				newGQL := fmt.Sprintf("%v", ref.GQL)
				if existingGQL != newGQL {
					panic(fmt.Sprintf("non-unique key \"%s\", trying to replace %s with %s", key, existingGQL, newGQL))
				}
			}
			ret[key] = ref

			ref = ref.Elem()
		}
	}
	return ret
}

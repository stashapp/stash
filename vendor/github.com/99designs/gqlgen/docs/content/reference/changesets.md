---
linkTitle: Changesets
title: Using maps as changesets
description: Falling back to map[string]interface{} to allow for presence checks.
menu: { main: { parent: 'reference', weight: 10 } }
---

Occasionally you need to distinguish presence from nil (undefined vs null). In gqlgen we do this using maps:


```graphql
type Mutation {
	updateUser(id: ID!, changes: UserChanges!): User
}

type UserChanges {
	name: String
	email: String
}
```

Then in config set the backing type to `map[string]interface{}`
```yaml
models:
  UserChanges:
    model: "map[string]interface{}"
```

After running go generate you should end up with a resolver that looks like this:
```go
func (r *mutationResolver) UpdateUser(ctx context.Context, id int, changes map[string]interface{}) (*User, error) {
	u := fetchFromDb(id)
	/// apply the changes
	saveToDb(u)
	return u, nil
}
```

We often use the mapstructure library to directly apply these changesets directly to the object using reflection:
```go

func ApplyChanges(changes map[string]interface{}, to interface{}) error {
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: true,
		TagName:     "json",
		Result:      to,
		ZeroFields:  true,
		// This is needed to get mapstructure to call the gqlgen unmarshaler func for custom scalars (eg Date)
		DecodeHook: func(a reflect.Type, b reflect.Type, v interface{}) (interface{}, error) {
			if reflect.PtrTo(b).Implements(reflect.TypeOf((*graphql.Unmarshaler)(nil)).Elem()) {
				resultType := reflect.New(b)
				result := resultType.MethodByName("UnmarshalGQL").Call([]reflect.Value{reflect.ValueOf(v)})
				err, _ := result[0].Interface().(error)
				return resultType.Elem().Interface(), err
			}

			return v, nil
		},
	})

	if err != nil {
		return err
	}

	return dec.Decode(changes)
}
```

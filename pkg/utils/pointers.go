package utils

func FirstNotNil[T any](a ...*T) *T {
	for _, v := range a {
		if v != nil {
			return v
		}
	}
	return nil
}

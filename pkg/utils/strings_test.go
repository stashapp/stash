package utils

import "fmt"

func ExampleStrFormat() {
	fmt.Println(StrFormat("{foo} bar {baz}", StrFormatMap{
		"foo": "bar",
		"baz": "abc",
	}))
	// Output:
	// bar bar abc
}

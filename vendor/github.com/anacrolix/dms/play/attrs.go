//go:build ignore
// +build ignore

package main

import (
	"encoding/xml"
	"fmt"
)

type Meh struct {
	XMLName xml.Name
	Size
	// ChildCount *uint `xml:"childCount,attr"`
}

func main() {
	size := uint64(137)
	data, err := xml.Marshal(Meh{
		Size: &xml.Attr{
			Name:  xml.Name{Local: "size"},
			Value: fmt.Sprint(size),
		},
	})
	fmt.Println(string(data), err)
}

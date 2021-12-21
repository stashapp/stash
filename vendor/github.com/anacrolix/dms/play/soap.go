//go:build ignore
// +build ignore

package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/anacrolix/dms/soap"
)

type Browse struct {
	ObjectID       string
	BrowseFlag     string
	Filter         string
	StartingIndex  int
	RequestedCount int
}

type GetSortCapabilitiesResponse struct {
	XMLName  xml.Name `xml:"urn:schemas-upnp-org:service:ContentDirectory:1 GetSortCapabilitiesResponse"`
	SortCaps string
}

func main() {
	raw, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	var env soap.Envelope
	if err := xml.Unmarshal(raw, &env); err != nil {
		panic(err)
	}
	fmt.Println(env)
	var browse Browse
	err = xml.Unmarshal([]byte(env.Body.Action), &browse)
	if err != nil {
		panic(err)
	}
	fmt.Println(browse)
	raw, err = xml.MarshalIndent(
		GetSortCapabilitiesResponse{
			SortCaps: "dc:title",
		},
		"", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(raw))
}

//go:build ignore
// +build ignore

package main

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/anacrolix/dms/upnp"
)

func main() {
	scpd := upnp.SCPD{
		SpecVersion: upnp.SpecVersion{Major: 1, Minor: 0},
		ActionList: []upnp.Action{
			{
				Name: "Browse",
				Arguments: []upnp.Argument{
					{Name: "ObjectID", Direction: "in", RelatedStateVar: "A_ARG_TYPE_ObjectID"},
				},
			},
		},
		ServiceStateTable: []upnp.StateVariable{
			{
				SendEvents: "no", Name: "A_ARG_TYPE_ObjectID", DataType: "string",
				AllowedValues: &[]string{"hi", "there"},
			},
			{
				SendEvents: "yes",
				Name:       "loltype",
			},
		},
	}
	xml, err := xml.MarshalIndent(scpd, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(string(xml))
}

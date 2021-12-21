package upnp

import (
	"encoding/xml"
	"testing"
)

// Visually verify that property sets are marshalled correctly.
func TestMarshalPropertySet(t *testing.T) {
	b, err := xml.MarshalIndent(&PropertySet{
		Properties: []Property{
			{
				Variable: Variable{
					XMLName: xml.Name{
						Local: "SystemUpdateID",
					},
					Value: "0",
				},
			},
			{
				Variable: Variable{
					XMLName: xml.Name{
						Local: "answerToTheUniverse",
					},
					Value: "42",
				},
			},
		},
		Space: "urn:schemas-upnp-org:event-1-0",
	}, "", "  ")
	t.Log("\n" + string(b))
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseCallbackURLs(t *testing.T) {
	urls := ParseCallbackURLs("<http://hello><http://path>     <http://world>")
	if len(urls) != 3 {
		t.Fatal(len(urls))
	}
}

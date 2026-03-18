package migrations

import (
	"encoding/json"
	"testing"
)

func TestConvertObjectFilter(t *testing.T) {
	migrator := schema86Migrator{}
	input := `{
		"tags": {
			"modifier": "INCLUDES",
			"value": {
				"depth": 0,
				"excluded": [
					{
						"id": "27",
						"label": "JAV Actress"
					}
				],
				"items": [
					{
						"id": "28",
						"label": "xyz"
					}
				]
			}
		}
	}`

	expected := `{"tags":{"depth":0,"excludes":["27"],"modifier":"INCLUDES","value":["28"]}}`

	output, err := migrator.convertObjectFilter([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var outMap, expMap map[string]interface{}
	if err := json.Unmarshal(output, &outMap); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if err := json.Unmarshal([]byte(expected), &expMap); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	outJSON, err := json.Marshal(outMap)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	expJSON, err := json.Marshal(expMap)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	if string(outJSON) != string(expJSON) {
		t.Errorf("expected %s, got %s", string(expJSON), string(outJSON))
	}
}

func TestConvertObjectFilterPrimitive(t *testing.T) {
	migrator := schema86Migrator{}
	input := `{
		"galleries": {
			"modifier": "INCLUDES",
			"value": {
				"excluded": [],
				"items": [
					{
						"id": "1",
						"label": "gallery 1"
					}
				]
			}
		},
		"has_markers": {
			"modifier": "EQUALS",
			"value": "true"
		}
	}`

	expected := `{"galleries":{"excludes":[],"modifier":"INCLUDES","value":["1"]},"has_markers":{"modifier":"EQUALS","value":"true"}}`

	output, err := migrator.convertObjectFilter([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var outMap, expMap map[string]interface{}
	if err := json.Unmarshal(output, &outMap); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if err := json.Unmarshal([]byte(expected), &expMap); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	outJSON, err := json.Marshal(outMap)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	expJSON, err := json.Marshal(expMap)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	if string(outJSON) != string(expJSON) {
		t.Errorf("expected %s, got %s", string(expJSON), string(outJSON))
	}
}

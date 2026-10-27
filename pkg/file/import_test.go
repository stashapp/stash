package file

import (
	"context"
	gojson "encoding/json"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stretchr/testify/assert"
)

func TestImporterPreImport(t *testing.T) {
	type testCase struct {
		name     string
		input    jsonschema.DirEntry
		expected *models.BaseFile
		err      bool
	}

	testCases := []testCase{
		{
			name: "phash as hex string",
			input: &jsonschema.BaseFile{
				Fingerprints: []jsonschema.Fingerprint{
					{
						Type:        models.FingerprintTypePhash,
						Fingerprint: gojson.RawMessage(`"0123456789abcdef"`),
					},
				},
			},
			expected: &models.BaseFile{
				Fingerprints: []models.Fingerprint{
					{
						Type:        models.FingerprintTypePhash,
						Fingerprint: int64(0x0123456789abcdef),
					},
				},
			},
		},
		{
			name: "phash as legacy int64",
			input: &jsonschema.BaseFile{
				Fingerprints: []jsonschema.Fingerprint{
					{
						Type:        models.FingerprintTypePhash,
						Fingerprint: gojson.RawMessage(`1234567890123456789`),
					},
				},
			},
			expected: &models.BaseFile{
				Fingerprints: []models.Fingerprint{
					{
						Type:        models.FingerprintTypePhash,
						Fingerprint: int64(1234567890123456789),
					},
				},
			},
		},
		{
			name: "phash as legacy int64 #6894",
			input: &jsonschema.BaseFile{
				Fingerprints: []jsonschema.Fingerprint{
					{
						Type:        models.FingerprintTypePhash,
						Fingerprint: gojson.RawMessage(`6391326009969271747`),
					},
				},
			},
			expected: &models.BaseFile{
				Fingerprints: []models.Fingerprint{
					{
						Type:        models.FingerprintTypePhash,
						Fingerprint: int64(6391326009969271747),
					},
				},
			},
		},
		{
			name: "floating point fingerprint value",
			input: &jsonschema.BaseFile{
				Fingerprints: []jsonschema.Fingerprint{
					{
						Type:        models.FingerprintTypeMD5,
						Fingerprint: gojson.RawMessage(`12345.6789`),
					},
				},
			},
			err: true,
		},
		{
			name: "invalid hex string phash",
			input: &jsonschema.BaseFile{
				Fingerprints: []jsonschema.Fingerprint{
					{
						Type:        models.FingerprintTypePhash,
						Fingerprint: gojson.RawMessage(`"not a hex string"`),
					},
				},
			},
			err: true,
		},
	}

	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)

			importer := &Importer{
				Input: tc.input,
			}
			err := importer.PreImport(ctx)
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.err {
				assert.Error(err)
				return
			}

			// unset timestamps and basename to avoid test failures due to time differences
			if importer.file != nil {
				importer.file.Base().CreatedAt = time.Time{}
				importer.file.Base().UpdatedAt = time.Time{}
				importer.file.Base().ModTime = time.Time{}
				importer.file.Base().Basename = ""
			}

			assert.Equal(tc.expected, importer.file)
		})
	}

}

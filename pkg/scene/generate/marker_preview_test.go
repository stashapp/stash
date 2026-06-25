package generate

import "testing"

func TestMarkerPreviewDuration(t *testing.T) {
	ptr := func(v float64) *float64 { return &v }

	const defaultDuration = 20

	tests := []struct {
		name            string
		seconds         float64
		endSeconds      *float64
		maxDuration     int
		defaultDuration int
		wantDuration    float64
		wantGenerate    bool
	}{
		{
			name:            "nil end uses default",
			seconds:         10,
			endSeconds:      nil,
			maxDuration:     0,
			defaultDuration: defaultDuration,
			wantDuration:    20,
			wantGenerate:    true,
		},
		{
			name:            "explicit interval honored with ceiling disabled",
			seconds:         10,
			endSeconds:      ptr(40),
			maxDuration:     0,
			defaultDuration: defaultDuration,
			wantDuration:    30,
			wantGenerate:    true,
		},
		{
			name:            "interval within ceiling honored",
			seconds:         10,
			endSeconds:      ptr(40),
			maxDuration:     60,
			defaultDuration: defaultDuration,
			wantDuration:    30,
			wantGenerate:    true,
		},
		{
			name:            "interval exceeding ceiling is capped",
			seconds:         10,
			endSeconds:      ptr(40),
			maxDuration:     20,
			defaultDuration: defaultDuration,
			wantDuration:    20,
			wantGenerate:    true,
		},
		{
			name:            "negative ceiling disables the cap",
			seconds:         10,
			endSeconds:      ptr(40),
			maxDuration:     -5,
			defaultDuration: defaultDuration,
			wantDuration:    30,
			wantGenerate:    true,
		},
		{
			name:            "fractional interval is preserved",
			seconds:         12.5,
			endSeconds:      ptr(47.25),
			maxDuration:     0,
			defaultDuration: defaultDuration,
			wantDuration:    34.75,
			wantGenerate:    true,
		},
		{
			name:            "interval equal to ceiling is honored",
			seconds:         10,
			endSeconds:      ptr(30),
			maxDuration:     20,
			defaultDuration: defaultDuration,
			wantDuration:    20,
			wantGenerate:    true,
		},
		{
			name:            "zero interval falls back to default",
			seconds:         10,
			endSeconds:      ptr(10),
			maxDuration:     0,
			defaultDuration: defaultDuration,
			wantDuration:    20,
			wantGenerate:    true,
		},
		{
			name:            "negative interval falls back to default",
			seconds:         10,
			endSeconds:      ptr(5),
			maxDuration:     0,
			defaultDuration: defaultDuration,
			wantDuration:    20,
			wantGenerate:    true,
		},
		{
			name:            "non-positive default with nil end skips generation",
			seconds:         10,
			endSeconds:      nil,
			maxDuration:     0,
			defaultDuration: 0,
			wantDuration:    0,
			wantGenerate:    false,
		},
		{
			name:            "non-positive default with zero interval skips generation",
			seconds:         10,
			endSeconds:      ptr(10),
			maxDuration:     0,
			defaultDuration: 0,
			wantDuration:    0,
			wantGenerate:    false,
		},
		{
			name:            "negative default with nil end skips generation",
			seconds:         10,
			endSeconds:      nil,
			maxDuration:     0,
			defaultDuration: -10,
			wantDuration:    -10,
			wantGenerate:    false,
		},
		{
			name:            "valid interval generates despite misconfigured default",
			seconds:         10,
			endSeconds:      ptr(40),
			maxDuration:     0,
			defaultDuration: 0,
			wantDuration:    30,
			wantGenerate:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDuration, gotGenerate := markerPreviewDuration(tt.seconds, tt.endSeconds, tt.maxDuration, tt.defaultDuration)
			if gotDuration != tt.wantDuration {
				t.Errorf("duration = %v, want %v", gotDuration, tt.wantDuration)
			}
			if gotGenerate != tt.wantGenerate {
				t.Errorf("generate = %v, want %v", gotGenerate, tt.wantGenerate)
			}
		})
	}
}

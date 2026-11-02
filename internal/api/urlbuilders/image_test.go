package urlbuilders

import "testing"

func TestImageURLBuilder_cacheBuster(t *testing.T) {
	tests := []struct {
		name      string
		checksum  string
		updatedAt string
		want      string
	}{
		{
			name:      "prefers checksum when present",
			checksum:  "abc123",
			updatedAt: "1700000000",
			want:      "abc123",
		},
		{
			name:      "falls back to updatedAt when checksum is empty",
			checksum:  "",
			updatedAt: "1700000000",
			want:      "1700000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := ImageURLBuilder{Checksum: tt.checksum, UpdatedAt: tt.updatedAt}
			if got := b.cacheBuster(); got != tt.want {
				t.Errorf("cacheBuster() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestImageURLBuilder_GetImageURL(t *testing.T) {
	tests := []struct {
		name      string
		checksum  string
		updatedAt string
		want      string
	}{
		{
			name:      "uses checksum so the URL is stable across metadata edits",
			checksum:  "abc123",
			updatedAt: "1700000000",
			want:      "http://localhost:9999/image/42/image?t=abc123",
		},
		{
			name:      "falls back to updatedAt when there is no primary file",
			checksum:  "",
			updatedAt: "1700000000",
			want:      "http://localhost:9999/image/42/image?t=1700000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := ImageURLBuilder{
				BaseURL:   "http://localhost:9999",
				ImageID:   "42",
				Checksum:  tt.checksum,
				UpdatedAt: tt.updatedAt,
			}
			if got := b.GetImageURL(); got != tt.want {
				t.Errorf("GetImageURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestImageURLBuilder_GetThumbnailURL(t *testing.T) {
	tests := []struct {
		name      string
		checksum  string
		updatedAt string
		want      string
	}{
		{
			name:      "uses checksum so the URL is stable across metadata edits",
			checksum:  "abc123",
			updatedAt: "1700000000",
			want:      "http://localhost:9999/image/42/thumbnail?t=abc123",
		},
		{
			name:      "falls back to updatedAt when there is no primary file",
			checksum:  "",
			updatedAt: "1700000000",
			want:      "http://localhost:9999/image/42/thumbnail?t=1700000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := ImageURLBuilder{
				BaseURL:   "http://localhost:9999",
				ImageID:   "42",
				Checksum:  tt.checksum,
				UpdatedAt: tt.updatedAt,
			}
			if got := b.GetThumbnailURL(); got != tt.want {
				t.Errorf("GetThumbnailURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

package config

type ScanMetadataOptions struct {
	// Generate scene covers during scan
	ScanGenerateCovers bool `json:"scanGenerateCovers"`
	// Generate previews during scan
	ScanGeneratePreviews bool `json:"scanGeneratePreviews"`
	// Generate image previews during scan
	ScanGenerateImagePreviews bool `json:"scanGenerateImagePreviews"`
	// Generate sprites during scan
	ScanGenerateSprites bool `json:"scanGenerateSprites"`
	// Generate phashes during scan
	ScanGeneratePhashes bool `json:"scanGeneratePhashes"`
	// Generate image thumbnails during scan
	ScanGenerateThumbnails bool `json:"scanGenerateThumbnails"`
	// Generate image thumbnails during scan
	ScanGenerateClipPreviews bool `json:"scanGenerateClipPreviews"`
}

type AutoTagMetadataOptions struct {
	// IDs of performers to tag files with, or "*" for all
	Performers []string `json:"performers"`
	// IDs of studios to tag files with, or "*" for all
	Studios []string `json:"studios"`
	// IDs of tags to tag files with, or "*" for all
	Tags []string `json:"tags"`
}

package config

type ScanMetadataOptions struct {
	// Set name, date, details from metadata (if present)
	UseFileMetadata bool `json:"useFileMetadata"`
	// Strip file extension from title
	StripFileExtension bool `json:"stripFileExtension"`
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
}

type AutoTagMetadataOptions struct {
	// IDs of performers to tag files with, or "*" for all
	Performers []string `json:"performers"`
	// IDs of studios to tag files with, or "*" for all
	Studios []string `json:"studios"`
	// IDs of tags to tag files with, or "*" for all
	Tags []string `json:"tags"`
}

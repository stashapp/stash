package file

// VideoFile is an extension of BaseFile to represent video files.
type VideoFile struct {
	*BaseFile
	Format     string  `json:"format"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	Duration   float64 `json:"duration"`
	VideoCodec string  `json:"video_codec"`
	AudioCodec string  `json:"audio_codec"`
	FrameRate  float64 `json:"frame_rate"`
	BitRate    int64   `json:"bitrate"`

	Interactive      bool `json:"interactive"`
	InteractiveSpeed *int `json:"interactive_speed"`
}

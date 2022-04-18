package ffmpeg

type VideoCodec string

func (c VideoCodec) Args() []string {
	if c == "" {
		return nil
	}

	return []string{"-c:v", string(c)}
}

var (
	VideoCodecLibX264 VideoCodec = "libx264"
	VideoCodecLibWebP VideoCodec = "libwebp"
	VideoCodecBMP     VideoCodec = "bmp"
	VideoCodecMJpeg   VideoCodec = "mjpeg"
	VideoCodecVP9     VideoCodec = "libvpx-vp9"
	VideoCodecVPX     VideoCodec = "libvpx"
	VideoCodecLibX265 VideoCodec = "libx265"
	VideoCodecCopy    VideoCodec = "copy"
)

type AudioCodec string

func (c AudioCodec) Args() []string {
	if c == "" {
		return nil
	}

	return []string{"-c:a", string(c)}
}

var (
	AudioCodecAAC     AudioCodec = "aac"
	AudioCodecLibOpus AudioCodec = "libopus"
	AudioCodecCopy    AudioCodec = "copy"
)

package ffmpeg

type VideoCodec struct {
	Name     string // The full name of the codec including profile/quality
	CodeName string // The core codec name without profile/quality suffix
}

func makeVideoCodec(name string, codename string) VideoCodec {
	return VideoCodec{name, codename}
}

func (c VideoCodec) Args() []string {
	if c.CodeName == "" {
		return nil
	}

	return []string{"-c:v", string(c.CodeName)}
}

var (
	// Software codec's
	VideoCodecLibX264 VideoCodec = makeVideoCodec("x264", "libx264")
	VideoCodecLibWebP VideoCodec = makeVideoCodec("WebP", "libwebp")
	VideoCodecBMP     VideoCodec = makeVideoCodec("BMP", "bmp")
	VideoCodecMJpeg   VideoCodec = makeVideoCodec("Jpeg", "mjpeg")
	VideoCodecVP9     VideoCodec = makeVideoCodec("VPX-VP9", "libvpx-vp9")
	VideoCodecVPX     VideoCodec = makeVideoCodec("VPX-VP8", "libvpx")
	VideoCodecLibX265 VideoCodec = makeVideoCodec("x265", "libx265")
	VideoCodecCopy    VideoCodec = makeVideoCodec("Copy", "copy")
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

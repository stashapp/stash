package ffmpeg

type VideoCodec string

func (c VideoCodec) Args() []string {
	if c == "" {
		return nil
	}

	return []string{"-c:v", string(c)}
}

var (
	// Software codec's
	VideoCodecLibX264 VideoCodec = "libx264"
	VideoCodecLibWebP VideoCodec = "libwebp"
	VideoCodecBMP     VideoCodec = "bmp"
	VideoCodecMJpeg   VideoCodec = "mjpeg"
	VideoCodecVP9     VideoCodec = "libvpx-vp9"
	VideoCodecVPX     VideoCodec = "libvpx"
	VideoCodecLibX265 VideoCodec = "libx265"
	VideoCodecCopy    VideoCodec = "copy"

	// Hardware codec's
	VideoCodecN264 VideoCodec = "h264_nvenc"
	VideoCodecI264 VideoCodec = "h264_qsv"
	VideoCodecA264 VideoCodec = "h264_amf"
	VideoCodecM264 VideoCodec = "h264_videotoolbox"
	VideoCodecV264 VideoCodec = "h264_vaapi"
	VideoCodecR264 VideoCodec = "h264_v4l2m2m"
	VideoCodecO264 VideoCodec = "h264_omx"
	VideoCodecIVP9 VideoCodec = "vp9_qsv"
	VideoCodecVVP9 VideoCodec = "vp9_vaapi"
	VideoCodecVVPX VideoCodec = "vp8_vaapi"
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

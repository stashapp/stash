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
	VideoCodecLibX264 = makeVideoCodec("x264", "libx264")
	VideoCodecLibWebP = makeVideoCodec("WebP", "libwebp")
	VideoCodecBMP     = makeVideoCodec("BMP", "bmp")
	VideoCodecMJpeg   = makeVideoCodec("Jpeg", "mjpeg")
	VideoCodecVP9     = makeVideoCodec("VPX-VP9", "libvpx-vp9")
	VideoCodecVPX     = makeVideoCodec("VPX-VP8", "libvpx")
	VideoCodecLibX265 = makeVideoCodec("x265", "libx265")
	VideoCodecCopy    = makeVideoCodec("Copy", "copy")
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

func (codec VideoCodec) ExtraArgs() (args Args) {
	args = args.VideoCodec(codec)

	switch codec {
	// CPU Codecs
	case VideoCodecLibX264:
		args = append(args,
			"-pix_fmt", "yuv420p",
			"-preset", "veryfast",
			"-crf", "25",
			"-sc_threshold", "0",
		)
	case VideoCodecVP9:
		args = append(args,
			"-pix_fmt", "yuv420p",
			"-deadline", "realtime",
			"-cpu-used", "5",
			"-row-mt", "1",
			"-crf", "30",
			"-b:v", "0",
		)
	// HW Codecs
	case VideoCodecN264:
		args = append(args,
			"-rc", "vbr",
			"-cq", "15",
		)
	case VideoCodecN264H:
		args = append(args,
			"-profile", "p7",
			"-tune", "hq",
			"-profile", "high",
			"-rc", "vbr",
			"-rc-lookahead", "60",
			"-surfaces", "64",
			"-spatial-aq", "1",
			"-aq-strength", "15",
			"-cq", "15",
			"-coder", "cabac",
			"-b_ref_mode", "middle",
		)
	case VideoCodecI264, VideoCodecIVP9:
		args = append(args,
			"-global_quality", "20",
			"-preset", "faster",
		)
	case VideoCodecI264C:
		args = append(args,
			"-q", "20",
			"-preset", "faster",
		)
	case VideoCodecV264, VideoCodecVVP9:
		args = append(args,
			"-qp", "20",
		)
	case VideoCodecA264:
		args = append(args,
			"-quality", "speed",
		)
	case VideoCodecM264:
		args = append(args,
			"-realtime", "1",
		)
	case VideoCodecO264:
		args = append(args,
			"-preset", "superfast",
			"-crf", "25",
		)
	}

	return args
}

func (codec VideoCodec) ExtraArgsHQ(preset string) (args Args) {
	switch codec {
	// CPU Codecs
	case VideoCodecLibX264:
		args = append(args,
			"-pix_fmt", "yuv420p",
			"-profile:v", "high",
			"-preset", preset,
			"-crf", "21",
			"-sc_threshold", "0",
		)
	case VideoCodecVP9:
		args = append(args,
			"-pix_fmt", "yuv420p",
			"-deadline", "good",
			"-cpu-used", "0",
			"-row-mt", "1",
			"-crf", "20",
			"-b:v", "0",
		)
	// HW Codecs - adjust for higher quality
	case VideoCodecN264:
		args = append(args,
			"-rc", "vbr",
			"-cq", "10",
		)
	case VideoCodecN264H:
		args = append(args,
			"-profile", "p7",
			"-tune", "hq",
			"-profile", "high",
			"-rc", "vbr",
			"-rc-lookahead", "60",
			"-surfaces", "64",
			"-spatial-aq", "1",
			"-aq-strength", "15",
			"-cq", "10",
			"-coder", "cabac",
			"-b_ref_mode", "middle",
		)
	case VideoCodecI264, VideoCodecIVP9:
		args = append(args,
			"-global_quality", "15",
			"-preset", "slow",
		)
	case VideoCodecI264C:
		args = append(args,
			"-q", "15",
			"-preset", "slow",
		)
	case VideoCodecV264, VideoCodecVVP9:
		args = append(args,
			"-qp", "15",
		)
	case VideoCodecA264:
		args = append(args,
			"-quality", "quality",
		)
	case VideoCodecM264:
		args = append(args,
			"-realtime", "0",
		)
	case VideoCodecO264:
		args = append(args,
			"-preset", "slow",
			"-crf", "21",
		)
	}

	return args
}

package ffmpeg

type Container string
type ProbeAudioCodec string

const (
	Mp4      Container = "mp4"
	M4v      Container = "m4v"
	Mov      Container = "mov"
	Wmv      Container = "wmv"
	Webm     Container = "webm"
	Matroska Container = "matroska"
	Avi      Container = "avi"
	Flv      Container = "flv"
	Mpegts   Container = "mpegts"

	Aac                ProbeAudioCodec = "aac"
	Mp3                ProbeAudioCodec = "mp3"
	Opus               ProbeAudioCodec = "opus"
	Vorbis             ProbeAudioCodec = "vorbis"
	MissingUnsupported ProbeAudioCodec = ""

	Mp4Ffmpeg      string = "mov,mp4,m4a,3gp,3g2,mj2" // browsers support all of them
	M4vFfmpeg      string = "mov,mp4,m4a,3gp,3g2,mj2" // so we don't care that ffmpeg
	MovFfmpeg      string = "mov,mp4,m4a,3gp,3g2,mj2" // can't differentiate between them
	WmvFfmpeg      string = "asf"
	WebmFfmpeg     string = "matroska,webm"
	MatroskaFfmpeg string = "matroska,webm"
	AviFfmpeg      string = "avi"
	FlvFfmpeg      string = "flv"
	MpegtsFfmpeg   string = "mpegts"
	H264           string = "h264"
	H265           string = "h265" // found in rare cases from a faulty encoder
	Hevc           string = "hevc"
	Vp8            string = "vp8"
	Vp9            string = "vp9"
	Mkv            string = "mkv" // only used from the browser to indicate mkv support
	Hls            string = "hls" // only used from the browser to indicate hls support
)

var ffprobeToContainer = map[string]Container{
	Mp4Ffmpeg:      Mp4,
	WmvFfmpeg:      Wmv,
	AviFfmpeg:      Avi,
	FlvFfmpeg:      Flv,
	MpegtsFfmpeg:   Mpegts,
	MatroskaFfmpeg: Matroska,
}

func MatchContainer(format string, filePath string) (Container, error) { // match ffprobe string to our Container
	container := ffprobeToContainer[format]
	if container == Matroska {
		return magicContainer(filePath) // use magic number instead of ffprobe for matroska,webm
	}
	if container == "" { // if format is not in our Container list leave it as ffprobes reported format_name
		container = Container(format)
	}
	return container, nil
}

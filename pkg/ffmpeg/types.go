package ffmpeg

import (
	"time"

	"github.com/stashapp/stash/pkg/models"
)

type FFProbeJSON struct {
	Format struct {
		BitRate        string `json:"bit_rate"`
		Duration       string `json:"duration"`
		Filename       string `json:"filename"`
		FormatLongName string `json:"format_long_name"`
		FormatName     string `json:"format_name"`
		NbPrograms     int    `json:"nb_programs"`
		NbStreams      int    `json:"nb_streams"`
		ProbeScore     int    `json:"probe_score"`
		Size           string `json:"size"`
		StartTime      string `json:"start_time"`
		Tags           struct {
			CompatibleBrands string          `json:"compatible_brands"`
			CreationTime     models.JSONTime `json:"creation_time"`
			Encoder          string          `json:"encoder"`
			MajorBrand       string          `json:"major_brand"`
			MinorVersion     string          `json:"minor_version"`
			Title            string          `json:"title"`
			Comment          string          `json:"comment"`
		} `json:"tags"`
	} `json:"format"`
	Streams []FFProbeStream `json:"streams"`
	Error   struct {
		Code   int    `json:"code"`
		String string `json:"string"`
	} `json:"error"`
}

type FFProbeStream struct {
	AvgFrameRate       string `json:"avg_frame_rate"`
	BitRate            string `json:"bit_rate"`
	BitsPerRawSample   string `json:"bits_per_raw_sample,omitempty"`
	ChromaLocation     string `json:"chroma_location,omitempty"`
	CodecLongName      string `json:"codec_long_name"`
	CodecName          string `json:"codec_name"`
	CodecTag           string `json:"codec_tag"`
	CodecTagString     string `json:"codec_tag_string"`
	CodecTimeBase      string `json:"codec_time_base"`
	CodecType          string `json:"codec_type"`
	CodedHeight        int    `json:"coded_height,omitempty"`
	CodedWidth         int    `json:"coded_width,omitempty"`
	DisplayAspectRatio string `json:"display_aspect_ratio,omitempty"`
	Disposition        struct {
		AttachedPic     int `json:"attached_pic"`
		CleanEffects    int `json:"clean_effects"`
		Comment         int `json:"comment"`
		Default         int `json:"default"`
		Dub             int `json:"dub"`
		Forced          int `json:"forced"`
		HearingImpaired int `json:"hearing_impaired"`
		Karaoke         int `json:"karaoke"`
		Lyrics          int `json:"lyrics"`
		Original        int `json:"original"`
		TimedThumbnails int `json:"timed_thumbnails"`
		VisualImpaired  int `json:"visual_impaired"`
	} `json:"disposition"`
	Duration          string `json:"duration"`
	DurationTs        int    `json:"duration_ts"`
	HasBFrames        int    `json:"has_b_frames,omitempty"`
	Height            int    `json:"height,omitempty"`
	Index             int    `json:"index"`
	IsAvc             string `json:"is_avc,omitempty"`
	Level             int    `json:"level,omitempty"`
	NalLengthSize     string `json:"nal_length_size,omitempty"`
	NbFrames          string `json:"nb_frames"`
	PixFmt            string `json:"pix_fmt,omitempty"`
	Profile           string `json:"profile"`
	RFrameRate        string `json:"r_frame_rate"`
	Refs              int    `json:"refs,omitempty"`
	SampleAspectRatio string `json:"sample_aspect_ratio,omitempty"`
	StartPts          int    `json:"start_pts"`
	StartTime         string `json:"start_time"`
	Tags              struct {
		CreationTime models.JSONTime `json:"creation_time"`
		HandlerName  string          `json:"handler_name"`
		Language     string          `json:"language"`
		Rotate       string          `json:"rotate"`
	} `json:"tags"`
	TimeBase      string `json:"time_base"`
	Width         int    `json:"width,omitempty"`
	BitsPerSample int    `json:"bits_per_sample,omitempty"`
	ChannelLayout string `json:"channel_layout,omitempty"`
	Channels      int    `json:"channels,omitempty"`
	MaxBitRate    string `json:"max_bit_rate,omitempty"`
	SampleFmt     string `json:"sample_fmt,omitempty"`
	SampleRate    string `json:"sample_rate,omitempty"`
}
type Container string
type AudioCodec string

const (
	Mp4                Container  = "mp4"
	M4v                Container  = "m4v"
	Mov                Container  = "mov"
	Wmv                Container  = "wmv"
	Webm               Container  = "webm"
	Matroska           Container  = "matroska"
	Avi                Container  = "avi"
	Flv                Container  = "flv"
	Mpegts             Container  = "mpegts"
	Aac                AudioCodec = "aac"
	Mp3                AudioCodec = "mp3"
	Opus               AudioCodec = "opus"
	Vorbis             AudioCodec = "vorbis"
	MissingUnsupported AudioCodec = ""
	Mp4Ffmpeg          string     = "mov,mp4,m4a,3gp,3g2,mj2" // browsers support all of them
	M4vFfmpeg          string     = "mov,mp4,m4a,3gp,3g2,mj2" // so we don't care that ffmpeg
	MovFfmpeg          string     = "mov,mp4,m4a,3gp,3g2,mj2" // can't differentiate between them
	WmvFfmpeg          string     = "asf"
	WebmFfmpeg         string     = "matroska,webm"
	MatroskaFfmpeg     string     = "matroska,webm"
	AviFfmpeg          string     = "avi"
	FlvFfmpeg          string     = "flv"
	MpegtsFfmpeg       string     = "mpegts"
	H264               string     = "h264"
	H265               string     = "h265" // found in rare cases from a faulty encoder
	Hevc               string     = "hevc"
	Vp8                string     = "vp8"
	Vp9                string     = "vp9"
	MimeWebm           string     = "video/webm"
	MimeMkv            string     = "video/x-matroska"
)

var ValidCodecs = []string{H264, H265, Vp8, Vp9}

var validForH264Mkv = []Container{Mp4, Matroska}
var validForH264 = []Container{Mp4}
var validForH265Mkv = []Container{Mp4, Matroska}
var validForH265 = []Container{Mp4}
var validForVp8 = []Container{Webm}
var validForVp9Mkv = []Container{Webm, Matroska}
var validForVp9 = []Container{Webm}
var validForHevcMkv = []Container{Mp4, Matroska}
var validForHevc = []Container{Mp4}

var validAudioForMkv = []AudioCodec{Aac, Mp3, Vorbis, Opus}
var validAudioForWebm = []AudioCodec{Vorbis, Opus}
var validAudioForMp4 = []AudioCodec{Aac, Mp3}

//maps user readable container strings to ffprobe's format_name
//on some formats ffprobe can't differentiate
var ContainerToFfprobe = map[Container]string{
	Mp4:      Mp4Ffmpeg,
	M4v:      M4vFfmpeg,
	Mov:      MovFfmpeg,
	Wmv:      WmvFfmpeg,
	Webm:     WebmFfmpeg,
	Matroska: MatroskaFfmpeg,
	Avi:      AviFfmpeg,
	Flv:      FlvFfmpeg,
	Mpegts:   MpegtsFfmpeg,
}

var FfprobeToContainer = map[string]Container{
	Mp4Ffmpeg:      Mp4,
	WmvFfmpeg:      Wmv,
	AviFfmpeg:      Avi,
	FlvFfmpeg:      Flv,
	MpegtsFfmpeg:   Mpegts,
	MatroskaFfmpeg: Matroska,
}

type VideoFile struct {
	JSON        FFProbeJSON
	AudioStream *FFProbeStream
	VideoStream *FFProbeStream

	Path         string
	Title        string
	Comment      string
	Container    string
	Duration     float64
	StartTime    float64
	Bitrate      int64
	Size         int64
	CreationTime time.Time

	VideoCodec   string
	VideoBitrate int64
	Width        int
	Height       int
	FrameRate    float64
	Rotation     int64

	AudioCodec string
}

var X264Presets = map[models.StreamingProfile]string{
	models.StreamingProfileUltrafast: "ultrafast",
	models.StreamingProfileMedium:    "medium",
	models.StreamingProfileSlow:      "slow",
}

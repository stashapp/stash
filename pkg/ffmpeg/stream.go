package ffmpeg

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type Stream struct {
	Stdout   io.ReadCloser
	Process  *os.Process
	options  TranscodeStreamOptions
	mimeType string
}

func (s *Stream) Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", s.mimeType)
	w.WriteHeader(http.StatusOK)

	logger.Infof("[stream] transcoding video file to %s", s.mimeType)

	// handle if client closes the connection
	notify := r.Context().Done()
	go func() {
		<-notify
		s.Process.Kill()
	}()

	_, err := io.Copy(w, s.Stdout)
	if err != nil {
		logger.Errorf("[stream] error serving transcoded video file: %s", err.Error())
	}
}

type Codec struct {
	Codec     string
	format    string
	MimeType  string
	extraArgs []string
}

var CodecHLS = Codec{
	Codec:    "libx264",
	format:   "mpegts",
	MimeType: MimeMpegts,
	extraArgs: []string{
		"-acodec", "aac",
		"-pix_fmt", "yuv420p",
		"-preset", "veryfast",
		"-crf", "30",
	},
}

var CodecH264 = Codec{
	Codec:    "libx264",
	format:   "mp4",
	MimeType: MimeMp4,
	extraArgs: []string{
		"-movflags", "frag_keyframe",
		"-pix_fmt", "yuv420p",
		"-preset", "veryfast",
		"-crf", "30",
	},
}

var CodecVP9 = Codec{
	Codec:    "libvpx-vp9",
	format:   "webm",
	MimeType: MimeWebm,
	extraArgs: []string{
		"-deadline", "realtime",
		"-cpu-used", "5",
		"-row-mt", "1",
		"-crf", "30",
	},
}

var CodecVP8 = Codec{
	Codec:    "libvpx",
	format:   "webm",
	MimeType: MimeWebm,
	extraArgs: []string{
		"-deadline", "realtime",
		"-cpu-used", "5",
		"-crf", "30",
	},
}

var CodecHEVC = Codec{
	Codec:    "libx265",
	format:   "mp4",
	MimeType: MimeMp4,
	extraArgs: []string{
		"-movflags", "frag_keyframe",
		"-preset", "veryfast",
		"-crf", "30",
	},
}

// it is very common in MKVs to have just the audio codec unsupported
// copy the video stream, transcode the audio and serve as Matroska
var CodecMKVAudio = Codec{
	Codec:    "copy",
	format:   "matroska",
	MimeType: MimeMkv,
	extraArgs: []string{
		"-c:a", "libopus",
		"-b:a", "96k",
		"-vbr", "on",
	},
}

type TranscodeStreamOptions struct {
	ProbeResult      VideoFile
	Codec            Codec
	Hls              bool
	StartTime        string
	MaxTranscodeSize models.StreamingResolutionEnum
	// transcode the video, remove the audio
	// in some videos where the audio codec is not supported by ffmpeg
	// ffmpeg fails if you try to transcode the audio
	VideoOnly bool
}

func GetTranscodeStreamOptions(probeResult VideoFile, videoCodec string, audioCodec AudioCodec, container Container, supportedVideoCodecs []string) TranscodeStreamOptions {
	options := TranscodeStreamOptions{
		ProbeResult: probeResult,
	}

	options.setTranscodeCodec(supportedVideoCodecs)

	if audioCodec == MissingUnsupported {
		//ffmpeg fails if it trys to transcode a non supported audio codec
		options.VideoOnly = true
	} else {
		// try to be smart if the video to be transcoded is in a Matroska container
		// mp4 has always supported audio so it doesn't need to be checked
		// while mpeg_ts has seeking issues if we don't reencode the video

		// If MKV is supported and video codec is also supported then only transcode audio
		if IsValidCodec(Mkv, supportedVideoCodecs) && Container(container) == Matroska && IsValidCodec(videoCodec, supportedVideoCodecs) {
			// copy video stream instead of transcoding it
			options.Codec = CodecMKVAudio
		}
	}

	return options
}

func (o *TranscodeStreamOptions) setTranscodeCodec(supportedVideoCodecs []string) {
	if len(supportedVideoCodecs) == 0 {
		supportedVideoCodecs = DefaultSupportedCodecs
	}

	logger.Debugf("Choosing transcode codec from: %s", strings.Join(supportedVideoCodecs, ","))

	// TODO - make preferred order configurable
	if IsValidCodec(Vp9, supportedVideoCodecs) {
		logger.Debug("Using VP9")
		o.Codec = CodecVP9
	} else if IsValidCodec(Vp8, supportedVideoCodecs) {
		logger.Debug("Using VP8")
		o.Codec = CodecVP8
	} else if IsValidCodec(Hevc, supportedVideoCodecs) {
		logger.Debug("Using HEVC")
		o.Codec = CodecHEVC
	} else if IsValidCodec(Hls, supportedVideoCodecs) {
		logger.Debug("Using HLS (with H264)")
		o.Codec = CodecHLS
		o.Hls = true
	} else {
		logger.Debug("Using H264")
		o.Codec = CodecH264
	}
}

func (o TranscodeStreamOptions) getStreamArgs() []string {
	scale := calculateTranscodeScale(o.ProbeResult, o.MaxTranscodeSize)

	args := []string{
		"-hide_banner",
		"-v", "error",
	}

	if o.StartTime != "" {
		args = append(args, "-ss", o.StartTime)
	}

	if o.Hls {
		// we only serve a fixed segment length
		args = append(args, "-t", strconv.Itoa(int(hlsSegmentLength)))
	}

	args = append(args,
		"-i", o.ProbeResult.Path,
	)

	if o.VideoOnly {
		args = append(args, "-an")
	}

	args = append(args,
		"-c:v", o.Codec.Codec,
		"-vf", "scale="+scale,
	)

	if len(o.Codec.extraArgs) > 0 {
		args = append(args, o.Codec.extraArgs...)
	}

	args = append(args,
		// this is needed for 5-channel ac3 files
		"-ac", "2",
		"-b:v", "0",
		"-f", o.Codec.format,
		"pipe:",
	)

	return args
}

func (e *Encoder) GetTranscodeStream(options TranscodeStreamOptions) (*Stream, error) {
	return e.stream(options.ProbeResult, options)
}

func (e *Encoder) stream(probeResult VideoFile, options TranscodeStreamOptions) (*Stream, error) {
	args := options.getStreamArgs()
	cmd := exec.Command(e.Path, args...)
	logger.Debugf("Streaming via: %s", strings.Join(cmd.Args, " "))

	stdout, err := cmd.StdoutPipe()
	if nil != err {
		logger.Error("FFMPEG stdout not available: " + err.Error())
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if nil != err {
		logger.Error("FFMPEG stderr not available: " + err.Error())
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	registerRunningEncoder(probeResult.Path, cmd.Process)
	go waitAndDeregister(probeResult.Path, cmd)

	// stderr must be consumed or the process deadlocks
	go func() {
		stderrData, _ := ioutil.ReadAll(stderr)
		stderrString := string(stderrData)
		if len(stderrString) > 0 {
			logger.Debugf("[stream] ffmpeg stderr: %s", stderrString)
		}
	}()

	ret := &Stream{
		Stdout:   stdout,
		Process:  cmd.Process,
		options:  options,
		mimeType: options.Codec.MimeType,
	}
	return ret, nil
}

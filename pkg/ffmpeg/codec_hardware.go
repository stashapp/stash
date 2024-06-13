package ffmpeg

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

var (
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

const minHeight int = 256

// Tests all (given) hardware codec's
func (f *FFMpeg) InitHWSupport(ctx context.Context) {
	var hwCodecSupport []VideoCodec

	for _, codec := range []VideoCodec{
		VideoCodecN264,
		VideoCodecI264,
		VideoCodecV264,
		VideoCodecR264,
		VideoCodecIVP9,
		VideoCodecVVP9,
	} {
		var args Args
		args = append(args, "-hide_banner")
		args = args.LogLevel(LogLevelWarning)
		args = f.hwDeviceInit(args, codec, false)
		args = args.Format("lavfi")
		args = args.Input(fmt.Sprintf("color=c=red:s=%dx%d", 1280, 720))
		args = args.Duration(0.1)

		// Test scaling
		videoFilter := f.hwMaxResFilter(codec, 1280, 720, minHeight, false)
		args = append(args, CodecInit(codec)...)
		args = args.VideoFilter(videoFilter)

		args = args.Format("null")
		args = args.Output("-")

		cmd := f.Command(ctx, args)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			errOutput := stderr.String()

			if len(errOutput) == 0 {
				errOutput = err.Error()
			}

			logger.Debugf("[InitHWSupport] Codec %s not supported. Error output:\n%s", codec, errOutput)
		} else {
			hwCodecSupport = append(hwCodecSupport, codec)
		}
	}

	outstr := fmt.Sprintf("[InitHWSupport] Supported HW codecs [%d]:\n", len(hwCodecSupport))
	for _, codec := range hwCodecSupport {
		outstr += fmt.Sprintf("\t%s\n", codec)
	}
	logger.Info(outstr)

	f.hwCodecSupport = hwCodecSupport
}

func (f *FFMpeg) hwCanFullHWTranscode(ctx context.Context, codec VideoCodec, vf *models.VideoFile, reqHeight int) bool {
	if codec == VideoCodecCopy {
		return false
	}

	var args Args
	args = append(args, "-hide_banner")
	args = args.LogLevel(LogLevelWarning)
	args = args.XError()
	args = f.hwDeviceInit(args, codec, true)
	args = args.Input(vf.Path)
	args = args.Duration(0.1)

	videoFilter := f.hwMaxResFilter(codec, vf.Width, vf.Height, reqHeight, true)
	args = append(args, CodecInit(codec)...)
	args = args.VideoFilter(videoFilter)

	args = args.Format("null")
	args = args.Output("-")

	cmd := f.Command(ctx, args)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errOutput := stderr.String()

		if len(errOutput) == 0 {
			errOutput = err.Error()
		}

		logger.Debugf("[InitHWSupport] Full hardware transcode for file %s not supported. Error output:\n%s", vf.Basename, errOutput)
		return false
	}

	return true
}

// Prepend input for hardware encoding only
func (f *FFMpeg) hwDeviceInit(args Args, toCodec VideoCodec, fullhw bool) Args {
	switch toCodec {
	case VideoCodecN264:
		args = append(args, "-hwaccel_device")
		args = append(args, "0")
		if fullhw {
			args = append(args, "-hwaccel")
			args = append(args, "cuda")
			args = append(args, "-hwaccel_output_format")
			args = append(args, "cuda")
			args = append(args, "-extra_hw_frames")
			args = append(args, "5")
		}
	case VideoCodecV264,
		VideoCodecVVP9:
		args = append(args, "-vaapi_device")
		args = append(args, "/dev/dri/renderD128")
		if fullhw {
			args = append(args, "-hwaccel")
			args = append(args, "vaapi")
			args = append(args, "-hwaccel_output_format")
			args = append(args, "vaapi")
		}
	case VideoCodecI264,
		VideoCodecIVP9:
		if fullhw {
			args = append(args, "-hwaccel")
			args = append(args, "qsv")
			args = append(args, "-hwaccel_output_format")
			args = append(args, "qsv")
		} else {
			args = append(args, "-init_hw_device")
			args = append(args, "qsv=hw")
			args = append(args, "-filter_hw_device")
			args = append(args, "hw")
		}
	}

	return args
}

// Initialise a video filter for HW encoding
func (f *FFMpeg) hwFilterInit(toCodec VideoCodec, fullhw bool) VideoFilter {
	var videoFilter VideoFilter
	switch toCodec {
	case VideoCodecV264,
		VideoCodecVVP9:
		if !fullhw {
			videoFilter = videoFilter.Append("format=nv12")
			videoFilter = videoFilter.Append("hwupload")
		}
	case VideoCodecN264:
		if !fullhw {
			videoFilter = videoFilter.Append("format=nv12")
			videoFilter = videoFilter.Append("hwupload_cuda")
		}
	case VideoCodecI264,
		VideoCodecIVP9:
		if !fullhw {
			videoFilter = videoFilter.Append("hwupload=extra_hw_frames=64")
			videoFilter = videoFilter.Append("format=qsv")
		}
	}

	return videoFilter
}

var scaler_re = regexp.MustCompile(`scale=(?P<value>[-\d]+:[-\d]+)`)

func templateReplaceScale(input string, template string, match []int, minusonehack bool) string {
	result := []byte{}

	res := string(scaler_re.ExpandString(result, template, input, match))

	// BUG: [scale_qsv]: Size values less than -1 are not acceptable.
	// Fix: Replace all instances of -2 with -1 in a scale operation
	if minusonehack {
		res = strings.ReplaceAll(res, "-2", "-1")
	}

	matchStart := match[0]
	matchEnd := match[1]

	return input[0:matchStart] + res + input[matchEnd:]
}

// Replace video filter scaling with hardware scaling for full hardware transcoding
func (f *FFMpeg) hwCodecFilter(args VideoFilter, codec VideoCodec, fullhw bool) VideoFilter {
	sargs := string(args)

	match := scaler_re.FindStringSubmatchIndex(sargs)
	if match == nil {
		return args
	}

	switch codec {
	case VideoCodecN264:
		template := "scale_cuda=$value"
		// In 10bit inputs you might get an error like "10 bit encode not supported"
		if fullhw && f.version.major >= 5 {
			template += ":format=nv12"
		}
		args = VideoFilter(templateReplaceScale(sargs, template, match, false))
	case VideoCodecV264,
		VideoCodecVVP9:
		template := "scale_vaapi=$value"
		args = VideoFilter(templateReplaceScale(sargs, template, match, false))
	case VideoCodecI264,
		VideoCodecIVP9:
		template := "scale_qsv=$value"
		args = VideoFilter(templateReplaceScale(sargs, template, match, true))
	}

	return args
}

// Returns the max resolution for a given codec, or a default
func (f *FFMpeg) hwCodecMaxRes(codec VideoCodec, dW int, dH int) (int, int) {
	switch codec {
	case VideoCodecN264,
		VideoCodecI264:
		return 4096, 4096
	}

	return dW, dH
}

// Return a maxres filter
func (f *FFMpeg) hwMaxResFilter(toCodec VideoCodec, width int, height int, reqHeight int, fullhw bool) VideoFilter {
	if width == 0 || height == 0 {
		return ""
	}
	videoFilter := f.hwFilterInit(toCodec, fullhw)
	maxWidth, maxHeight := f.hwCodecMaxRes(toCodec, width, height)
	videoFilter = videoFilter.ScaleMaxLM(width, height, reqHeight, maxWidth, maxHeight)
	return f.hwCodecFilter(videoFilter, toCodec, fullhw)
}

// Return if a hardware accelerated for HLS is available
func (f *FFMpeg) hwCodecHLSCompatible() *VideoCodec {
	for _, element := range f.hwCodecSupport {
		switch element {
		case VideoCodecN264,
			VideoCodecI264,
			VideoCodecV264,
			VideoCodecR264:
			return &element
		}
	}
	return nil
}

// Return if a hardware accelerated codec for MP4 is available
func (f *FFMpeg) hwCodecMP4Compatible() *VideoCodec {
	for _, element := range f.hwCodecSupport {
		switch element {
		case VideoCodecN264,
			VideoCodecI264:
			return &element
		}
	}
	return nil
}

// Return if a hardware accelerated codec for WebM is available
func (f *FFMpeg) hwCodecWEBMCompatible() *VideoCodec {
	for _, element := range f.hwCodecSupport {
		switch element {
		case VideoCodecIVP9,
			VideoCodecVVP9:
			return &element
		}
	}
	return nil
}

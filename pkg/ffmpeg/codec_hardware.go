package ffmpeg

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"regexp"
	"strconv"
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

const minHeight int = 480

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
		VideoCodecM264,
	} {
		var args Args
		args = append(args, "-hide_banner")
		args = args.LogLevel(LogLevelWarning)
		args = f.hwDeviceInit(args, codec, false)
		args = args.Format("lavfi")
		vFile := &models.VideoFile{Width: 1280, Height: 720}
		args = args.Input(fmt.Sprintf("color=c=red:s=%dx%d", vFile.Width, vFile.Height))
		args = args.Duration(0.1)

		// Test scaling
		videoFilter := f.hwMaxResFilter(codec, vFile, minHeight, false)
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
	args = args.Duration(1)

	videoFilter := f.hwMaxResFilter(codec, vf, reqHeight, true)
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
			args = append(args, "-threads")
			args = append(args, "1")
			args = append(args, "-hwaccel")
			args = append(args, "cuda")
			args = append(args, "-hwaccel_output_format")
			args = append(args, "cuda")
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
	case VideoCodecM264:
		if fullhw {
			args = append(args, "-hwaccel")
			args = append(args, "videotoolbox")
			args = append(args, "-hwaccel_output_format")
			args = append(args, "videotoolbox_vld")
		} else {
			args = append(args, "-init_hw_device")
			args = append(args, "videotoolbox=vt")
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
			videoFilter = videoFilter.Append("format=yuv420p")
			videoFilter = videoFilter.Append("hwupload_cuda")
		}
	case VideoCodecI264,
		VideoCodecIVP9:
		if !fullhw {
			videoFilter = videoFilter.Append("hwupload=extra_hw_frames=64")
			videoFilter = videoFilter.Append("format=qsv")
		}
	case VideoCodecM264:
		if !fullhw {
			videoFilter = videoFilter.Append("format=nv12")
			videoFilter = videoFilter.Append("hwupload")
		}
	}

	return videoFilter
}

var scaler_re = regexp.MustCompile(`scale=(?P<value>([-\d]+):([-\d]+))`)

func templateReplaceScale(input string, template string, match []int, vf *models.VideoFile, minusonehack bool) string {
	result := []byte{}

	if minusonehack {
		// Parse width and height
		w, err := strconv.Atoi(input[match[4]:match[5]])
		if err != nil {
			logger.Error("failed to parse width")
			return input
		}
		h, err := strconv.Atoi(input[match[6]:match[7]])
		if err != nil {
			logger.Error("failed to parse height")
			return input
		}

		// Calculate ratio
		ratio := float64(vf.Width) / float64(vf.Height)
		if w < 0 {
			w = int(math.Round(float64(h) * ratio))
		} else if h < 0 {
			h = int(math.Round(float64(w) / ratio))
		}

		// Fix not divisible by 2 errors
		if w%2 != 0 {
			w++
		}
		if h%2 != 0 {
			h++
		}

		template = strings.ReplaceAll(template, "$value", fmt.Sprintf("%d:%d", w, h))
	}

	res := string(scaler_re.ExpandString(result, template, input, match))

	matchStart := match[0]
	matchEnd := match[1]

	return input[0:matchStart] + res + input[matchEnd:]
}

// Replace video filter scaling with hardware scaling for full hardware transcoding (also fixes the format)
func (f *FFMpeg) hwCodecFilter(args VideoFilter, codec VideoCodec, vf *models.VideoFile, fullhw bool) VideoFilter {
	sargs := string(args)

	match := scaler_re.FindStringSubmatchIndex(sargs)
	if match == nil {
		return f.hwApplyFullHWFilter(args, codec, fullhw)
	}

	return f.hwApplyScaleTemplate(sargs, codec, match, vf, fullhw)
}

// Apply format switching if applicable
func (f *FFMpeg) hwApplyFullHWFilter(args VideoFilter, codec VideoCodec, fullhw bool) VideoFilter {
	switch codec {
	case VideoCodecN264:
		if fullhw && f.version.Gteq(FFMpegVersion{major: 5}) { // Added in FFMpeg 5
			args = args.Append("scale_cuda=format=yuv420p")
		}
	case VideoCodecV264, VideoCodecVVP9:
		if fullhw && f.version.Gteq(FFMpegVersion{major: 3, minor: 1}) { // Added in FFMpeg 3.1
			args = args.Append("scale_vaapi=format=nv12")
		}
	case VideoCodecI264, VideoCodecIVP9:
		if fullhw && f.version.Gteq(FFMpegVersion{major: 3, minor: 3}) { // Added in FFMpeg 3.3
			args = args.Append("scale_qsv=format=nv12")
		}
	}

	return args
}

// Switch scaler
func (f *FFMpeg) hwApplyScaleTemplate(sargs string, codec VideoCodec, match []int, vf *models.VideoFile, fullhw bool) VideoFilter {
	var template string

	switch codec {
	case VideoCodecN264:
		template = "scale_cuda=$value"
		if fullhw && f.version.Gteq(FFMpegVersion{major: 5}) { // Added in FFMpeg 5
			template += ":format=yuv420p"
		}
	case VideoCodecV264, VideoCodecVVP9:
		template = "scale_vaapi=$value"
		if fullhw && f.version.Gteq(FFMpegVersion{major: 3, minor: 1}) { // Added in FFMpeg 3.1
			template += ":format=nv12"
		}
	case VideoCodecI264, VideoCodecIVP9:
		template = "scale_qsv=$value"
		if fullhw && f.version.Gteq(FFMpegVersion{major: 3, minor: 3}) { // Added in FFMpeg 3.3
			template += ":format=nv12"
		}
	case VideoCodecM264:
		template = "scale_vt=$value"
	default:
		return VideoFilter(sargs)
	}

	// BUG: [scale_qsv]: Size values less than -1 are not acceptable.
	isIntel := codec == VideoCodecI264 || codec == VideoCodecIVP9
	// BUG: scale_vt doesn't call ff_scale_adjust_dimensions, thus cant accept negative size values
	isApple := codec == VideoCodecM264
	return VideoFilter(templateReplaceScale(sargs, template, match, vf, isIntel || isApple))
}

// Returns the max resolution for a given codec, or a default
func (f *FFMpeg) hwCodecMaxRes(codec VideoCodec) (int, int) {
	switch codec {
	case VideoCodecN264,
		VideoCodecI264:
		return 4096, 4096
	}

	return 0, 0
}

// Return a maxres filter
func (f *FFMpeg) hwMaxResFilter(toCodec VideoCodec, vf *models.VideoFile, reqHeight int, fullhw bool) VideoFilter {
	if vf.Width == 0 || vf.Height == 0 {
		return ""
	}
	videoFilter := f.hwFilterInit(toCodec, fullhw)
	maxWidth, maxHeight := f.hwCodecMaxRes(toCodec)
	videoFilter = videoFilter.ScaleMaxLM(vf.Width, vf.Height, reqHeight, maxWidth, maxHeight)
	return f.hwCodecFilter(videoFilter, toCodec, vf, fullhw)
}

// Return if a hardware accelerated for HLS is available
func (f *FFMpeg) hwCodecHLSCompatible() *VideoCodec {
	for _, element := range f.hwCodecSupport {
		switch element {
		case VideoCodecN264,
			VideoCodecI264,
			VideoCodecV264,
			VideoCodecR264,
			VideoCodecM264: // Note that the Apple encoder sucks at startup, thus HLS quality is crap
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
			VideoCodecI264,
			VideoCodecM264:
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

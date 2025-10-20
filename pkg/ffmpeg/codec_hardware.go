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
	VideoCodecN264  = makeVideoCodec("H264 NVENC", "h264_nvenc")
	VideoCodecN264H = makeVideoCodec("H264 NVENC HQ profile", "h264_nvenc")
	VideoCodecNHEVC = makeVideoCodec("HEVC NVENC", "hevc_nvenc")
	VideoCodecNAV1  = makeVideoCodec("AV1 NVENC", "av1_nvenc")
	VideoCodecI264  = makeVideoCodec("H264 Intel Quick Sync Video (QSV)", "h264_qsv")
	VideoCodecI264C = makeVideoCodec("H264 Intel Quick Sync Video (QSV) Compatibility profile", "h264_qsv")
	VideoCodecIHEVC = makeVideoCodec("HEVC Intel Quick Sync Video (QSV)", "hevc_qsv")
	VideoCodecIAV1  = makeVideoCodec("AV1 Intel Quick Sync Video (QSV)", "av1_qsv")
	VideoCodecIVP9  = makeVideoCodec("VP9 Intel Quick Sync Video (QSV)", "vp9_qsv")
	VideoCodecA264  = makeVideoCodec("H264 Advanced Media Framework (AMF)", "h264_amf")
	VideoCodecAHEVC = makeVideoCodec("HEVC Advanced Media Framework (AMF)", "hevc_amf")
	VideoCodecAAV1  = makeVideoCodec("AV1 Advanced Media Framework (AMF)", "av1_amf")
	VideoCodecM264  = makeVideoCodec("H264 VideoToolbox", "h264_videotoolbox")
	VideoCodecMHEVC = makeVideoCodec("HEVC VideoToolbox", "hevc_videotoolbox")
	VideoCodecMAV1  = makeVideoCodec("AV1 VideoToolbox", "av1_videotoolbox")
	VideoCodecV264  = makeVideoCodec("H264 VAAPI", "h264_vaapi")
	VideoCodecVHEVC = makeVideoCodec("HEVC VAAPI", "hevc_vaapi")
	VideoCodecVAV1  = makeVideoCodec("AV1 VAAPI", "av1_vaapi")
	VideoCodecVVP9  = makeVideoCodec("VP9 VAAPI", "vp9_vaapi")
	VideoCodecVVPX  = makeVideoCodec("VP8 VAAPI", "vp8_vaapi")
	VideoCodecR264  = makeVideoCodec("H264 V4L2M2M", "h264_v4l2m2m")
	VideoCodecO264  = makeVideoCodec("H264 OMX", "h264_omx")
)

const minHeight int = 480

// Tests all (given) hardware codec's
func (f *FFMpeg) InitHWSupport(ctx context.Context) {
	var hwCodecSupport []VideoCodec

	// Note that the first compatible codec is returned, so order is important
	// Priority: Modern codecs first, then legacy codecs
	for _, codec := range []VideoCodec{
		// NVIDIA modern codecs
		VideoCodecNAV1,
		VideoCodecNHEVC,
		VideoCodecN264H,
		VideoCodecN264,

		// Intel modern codecs
		VideoCodecIAV1,
		VideoCodecIHEVC,
		VideoCodecI264,
		VideoCodecI264C,

		// AMD modern codecs (VAAPI)
		VideoCodecVAV1,
		VideoCodecVHEVC,
		VideoCodecV264,
		VideoCodecVVP9,

		// AMD legacy codecs (AMF)
		VideoCodecAAV1,
		VideoCodecAHEVC,
		VideoCodecA264,

		// Apple modern codecs
		VideoCodecMAV1,
		VideoCodecMHEVC,
		VideoCodecM264,

		// Other legacy codecs
		VideoCodecR264,
		VideoCodecO264,
		VideoCodecVVPX,
		VideoCodecIVP9,
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
		outstr += fmt.Sprintf("\t%s - %s\n", codec.Name, codec.CodeName)
	}
	logger.Info(outstr)

	f.hwCodecSupport = hwCodecSupport
}

func (f *FFMpeg) hwCanFullHWTranscode(ctx context.Context, codec VideoCodec, vf *models.VideoFile, reqHeight int) bool {
	if codec == VideoCodecCopy {
		logger.Infof("[transcode] Codec is VideoCodecCopy, full hardware transcode not applicable")
		return false
	}

	logger.Infof("[transcode] Testing full hardware transcode for file %s with codec: %s (%s)", vf.Basename, codec.Name, codec.CodeName)

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

		logger.Infof("[transcode] Full hardware transcode test failed for file %s with codec %s: %s", vf.Basename, codec.Name, errOutput)
		logger.Debugf("[transcode] Full hardware transcode for file %s not supported. Error output:\n%s", vf.Basename, errOutput)
		return false
	}

	logger.Infof("[transcode] Full hardware transcode test passed for file %s with codec %s", vf.Basename, codec.Name)
	return true
}

// Prepend input for hardware encoding only
func (f *FFMpeg) hwDeviceInit(args Args, toCodec VideoCodec, fullhw bool) Args {
	switch toCodec {
	// NVIDIA codecs
	case VideoCodecN264,
		VideoCodecN264H,
		VideoCodecNHEVC,
		VideoCodecNAV1:
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

	// VAAPI codecs (AMD/Intel)
	case VideoCodecV264,
		VideoCodecVHEVC,
		VideoCodecVAV1,
		VideoCodecVVP9,
		VideoCodecVVPX:
		args = append(args, "-vaapi_device")
		args = append(args, "/dev/dri/renderD128")
		if fullhw {
			args = append(args, "-hwaccel")
			args = append(args, "vaapi")
			args = append(args, "-hwaccel_output_format")
			args = append(args, "vaapi")
		}

	// Intel QSV codecs
	case VideoCodecI264,
		VideoCodecI264C,
		VideoCodecIHEVC,
		VideoCodecIAV1,
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

	// Apple VideoToolbox codecs
	case VideoCodecM264,
		VideoCodecMHEVC,
		VideoCodecMAV1:
		if fullhw {
			args = append(args, "-hwaccel")
			args = append(args, "videotoolbox")
			args = append(args, "-hwaccel_output_format")
			args = append(args, "videotoolbox_vld")
		} else {
			args = append(args, "-init_hw_device")
			args = append(args, "videotoolbox=vt")
		}

	// AMD AMF codecs
	case VideoCodecA264,
		VideoCodecAHEVC,
		VideoCodecAAV1:
		// AMF uses software decoding with hardware encoding
		// No special hardware acceleration needed for decoding
		if fullhw {
			// For AMF, full hardware transcode is not typically supported
			// Use software decoding with hardware encoding
		}

	// Legacy codecs - no special hardware initialization needed
	case VideoCodecR264,
		VideoCodecO264:
		// V4L2M2M and OMX don't require special hardware initialization
	}

	return args
}

// Initialise a video filter for HW encoding
func (f *FFMpeg) hwFilterInit(toCodec VideoCodec, fullhw bool) VideoFilter {
	var videoFilter VideoFilter
	switch toCodec {
	// VAAPI codecs (AMD/Intel)
	case VideoCodecV264,
		VideoCodecVHEVC,
		VideoCodecVAV1,
		VideoCodecVVP9,
		VideoCodecVVPX:
		if !fullhw {
			videoFilter = videoFilter.Append("format=nv12")
			videoFilter = videoFilter.Append("hwupload")
		}

	// NVIDIA codecs
	case VideoCodecN264, VideoCodecN264H,
		VideoCodecNHEVC, VideoCodecNAV1:
		if !fullhw {
			videoFilter = videoFilter.Append("format=nv12")
			videoFilter = videoFilter.Append("hwupload_cuda")
		}

	// Intel QSV codecs
	case VideoCodecI264,
		VideoCodecI264C,
		VideoCodecIHEVC,
		VideoCodecIAV1,
		VideoCodecIVP9:
		if !fullhw {
			videoFilter = videoFilter.Append("hwupload=extra_hw_frames=64")
			videoFilter = videoFilter.Append("format=qsv")
		}

	// Apple VideoToolbox codecs
	case VideoCodecM264,
		VideoCodecMHEVC,
		VideoCodecMAV1:
		if !fullhw {
			videoFilter = videoFilter.Append("format=nv12")
			videoFilter = videoFilter.Append("hwupload")
		}

	// AMD AMF codecs
	case VideoCodecA264,
		VideoCodecAHEVC,
		VideoCodecAAV1:
		// AMF typically uses software decoding, so no special filter needed
		// The format conversion is handled by the encoder

	// Legacy codecs
	case VideoCodecR264,
		VideoCodecO264:
		// V4L2M2M and OMX don't require special filter initialization
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
	// NVIDIA codecs
	case VideoCodecN264, VideoCodecN264H,
		VideoCodecNHEVC, VideoCodecNAV1:
		if fullhw && f.version.Gteq(Version{major: 5}) { // Added in FFMpeg 5
			args = args.Append("scale_cuda=format=yuv420p")
		}

	// VAAPI codecs (AMD/Intel)
	case VideoCodecV264, VideoCodecVVP9,
		VideoCodecVHEVC, VideoCodecVAV1:
		if fullhw && f.version.Gteq(Version{major: 3, minor: 1}) { // Added in FFMpeg 3.1
			args = args.Append("scale_vaapi=format=nv12")
		}

	// Intel QSV codecs
	case VideoCodecI264, VideoCodecI264C, VideoCodecIVP9,
		VideoCodecIHEVC, VideoCodecIAV1:
		if fullhw && f.version.Gteq(Version{major: 3, minor: 3}) { // Added in FFMpeg 3.3
			args = args.Append("scale_qsv=format=nv12")
		}

	// Apple VideoToolbox codecs
	case VideoCodecM264, VideoCodecMHEVC, VideoCodecMAV1:
		if fullhw && f.version.Gteq(Version{major: 4, minor: 3}) { // Added in FFMpeg 4.3
			args = args.Append("scale_vt=format=nv12")
		}

	// AMD AMF codecs - typically don't support full hardware scaling
	// Legacy codecs - no hardware scaling support
	}

	return args
}

// Switch scaler
func (f *FFMpeg) hwApplyScaleTemplate(sargs string, codec VideoCodec, match []int, vf *models.VideoFile, fullhw bool) VideoFilter {
	var template string

	switch codec {
	// NVIDIA codecs
	case VideoCodecN264, VideoCodecN264H,
		VideoCodecNHEVC, VideoCodecNAV1:
		template = "scale_cuda=$value"
		if fullhw && f.version.Gteq(Version{major: 5}) { // Added in FFMpeg 5
			template += ":format=yuv420p"
		}

	// VAAPI codecs (AMD/Intel)
	case VideoCodecV264, VideoCodecVVP9,
		VideoCodecVHEVC, VideoCodecVAV1:
		template = "scale_vaapi=$value"
		if fullhw && f.version.Gteq(Version{major: 3, minor: 1}) { // Added in FFMpeg 3.1
			template += ":format=nv12"
		}

	// Intel QSV codecs
	case VideoCodecI264, VideoCodecI264C, VideoCodecIVP9,
		VideoCodecIHEVC, VideoCodecIAV1:
		template = "scale_qsv=$value"
		if fullhw && f.version.Gteq(Version{major: 3, minor: 3}) { // Added in FFMpeg 3.3
			template += ":format=nv12"
		}

	// Apple VideoToolbox codecs
	case VideoCodecM264, VideoCodecMHEVC, VideoCodecMAV1:
		template = "scale_vt=$value"

	// AMD AMF and legacy codecs - use software scaling
	default:
		return VideoFilter(sargs)
	}

	// BUG: [scale_qsv]: Size values less than -1 are not acceptable.
	isIntel := codec == VideoCodecI264 || codec == VideoCodecI264C || codec == VideoCodecIVP9 ||
		codec == VideoCodecIHEVC || codec == VideoCodecIAV1
	// BUG: scale_vt doesn't call ff_scale_adjust_dimensions, thus cant accept negative size values
	isApple := codec == VideoCodecM264 || codec == VideoCodecMHEVC || codec == VideoCodecMAV1
	return VideoFilter(templateReplaceScale(sargs, template, match, vf, isIntel || isApple))
}

// Returns the max resolution for a given codec, or a default
func (f *FFMpeg) hwCodecMaxRes(codec VideoCodec) (int, int) {
	switch codec {
	// Modern codecs with 8K support
	case VideoCodecNHEVC, VideoCodecNAV1,
		VideoCodecIHEVC, VideoCodecIAV1,
		VideoCodecVHEVC, VideoCodecVAV1,
		VideoCodecAHEVC, VideoCodecAAV1,
		VideoCodecMHEVC, VideoCodecMAV1:
		return 8192, 8192 // 8K support

	// Legacy codecs with 4K support
	case VideoCodecN264, VideoCodecN264H,
		VideoCodecI264, VideoCodecI264C,
		VideoCodecV264, VideoCodecVVP9,
		VideoCodecA264, VideoCodecM264:
		return 4096, 4096 // 4K support

	// Other codecs - use default resolution
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
		// H.264 codecs (best compatibility for HLS)
		case VideoCodecN264,
			VideoCodecN264H,
			VideoCodecI264,
			VideoCodecI264C,
			VideoCodecV264,
			VideoCodecR264,
			VideoCodecA264,
			VideoCodecM264: // Note that the Apple encoder sucks at startup, thus HLS quality is crap
			return &element

		// HEVC codecs (modern HLS support)
		case VideoCodecNHEVC,
			VideoCodecIHEVC,
			VideoCodecVHEVC,
			VideoCodecAHEVC,
			VideoCodecMHEVC:
			return &element
		}
	}
	return nil
}

// Return if a hardware accelerated codec for MP4 is available
func (f *FFMpeg) hwCodecMP4Compatible() *VideoCodec {
	for _, element := range f.hwCodecSupport {
		switch element {
		// H.264 codecs (best MP4 compatibility)
		case VideoCodecN264,
			VideoCodecN264H,
			VideoCodecI264,
			VideoCodecI264C,
			VideoCodecA264,
			VideoCodecM264:
			return &element

		// HEVC codecs (modern MP4 support)
		case VideoCodecNHEVC,
			VideoCodecIHEVC,
			VideoCodecAHEVC,
			VideoCodecMHEVC:
			return &element

		// AV1 codecs (latest MP4 support)
		case VideoCodecNAV1,
			VideoCodecIAV1,
			VideoCodecAAV1,
			VideoCodecMAV1:
			return &element
		}
	}
	return nil
}

// Return if a hardware accelerated codec for WebM is available
func (f *FFMpeg) hwCodecWEBMCompatible() *VideoCodec {
	for _, element := range f.hwCodecSupport {
		switch element {
		// VP9 codecs (native WebM support)
		case VideoCodecIVP9,
			VideoCodecVVP9:
			return &element

		// AV1 codecs (modern WebM support)
		case VideoCodecNAV1,
			VideoCodecIAV1,
			VideoCodecAAV1,
			VideoCodecMAV1,
			VideoCodecVAV1:
			return &element
		}
	}
	return nil
}

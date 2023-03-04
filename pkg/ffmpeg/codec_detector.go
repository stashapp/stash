package ffmpeg

import (
	"context"
	"regexp"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
)

var HWCodecSupport []VideoCodec

// Tests all (given) hardware codec's
func FindHWCodecs(ctx context.Context, encoder FFMpeg) {
	HWCodecSupport = HWCodecSupport[:0]

	for _, codec := range []VideoCodec{
		VideoCodecN264,
		VideoCodecI264,
		VideoCodecV264,
		VideoCodecIVP9,
		VideoCodecVVP9,
	} {
		var args Args
		args = append(args, "-hide_banner")
		args = args.LogLevel(LogLevelQuiet)
		args = HWDeviceInit(args, codec)
		args = args.Format("lavfi")
		args = args.Input("color=c=red")
		args = args.Duration(0.1)

		videoFilter := HWFilterInit(codec)
		// Test scaling
		videoFilter = videoFilter.ScaleDimensions(-2, 160)
		videoFilter = HWCodecFilter(videoFilter, codec)
		args = append(args, CodecInit(codec)...)
		args = args.VideoFilter(videoFilter)

		args = args.Format("null")
		args = args.Output("-")

		cmd := encoder.Command(ctx, args)

		if err := cmd.Run(); err == nil {
			HWCodecSupport = append(HWCodecSupport, codec)
		}
	}

	logger.Info("Supported HW codecs: ")
	for _, codec := range HWCodecSupport {
		logger.Info("\t", codec)
	}
}

// Return if given codec is hardware accelerated
func HWCodecDetect(codec VideoCodec) bool {
	switch codec {
	case VideoCodecN264,
		VideoCodecA264,
		VideoCodecM264,
		VideoCodecV264,
		VideoCodecI264,
		VideoCodecR264,
		VideoCodecO264,
		VideoCodecIVP9,
		VideoCodecVVP9:
		return true
	default:
		return false
	}
}

// Test full-hardware transcoding on an input video
/*func HWCodecVideoSupported(ctx context.Context, encoder FFMpeg, o TranscodeStreamOptions) bool {
	if !HWCodecDetect(o.Codec.codec) {
		return false
	}

	var args Args
	args = append(args, "-hide_banner")
	args = append(args, o.ExtraInputArgs...)
	args = args.LogLevel(LogLevelQuiet)
	args = HWDeviceInit_Full(args, o.Codec.codec)
	args = args.Input(o.Input)
	args = args.Duration(0.1)

	args = args.VideoCodec(o.Codec.codec)
	if len(o.Codec.extraArgs) > 0 {
		args = append(args, o.Codec.extraArgs...)
	}

	// Test scaling
	videoFilter := HWFilterInit(o.Codec.codec)
	videoFilter = videoFilter.ScaleDimensions(-2, 160)
	videoFilter = HWCodecFilter(videoFilter, o.Codec.codec)
	args = args.VideoFilter(videoFilter)

	args = args.Format("null")
	args = args.Output("-")

	cmd := encoder.Command(ctx, args)

	err := cmd.Run()
	return err == nil
}*/

// Prepend input for hardware encoding only
func HWDeviceInit(args Args, codec VideoCodec) Args {
	switch codec {
	case VideoCodecN264:
		args = append(args, "-hwaccel_device")
		args = append(args, "0")
	case VideoCodecV264,
		VideoCodecVVP9:
		args = append(args, "-vaapi_device")
		args = append(args, "/dev/dri/renderD128")
	case VideoCodecI264,
		VideoCodecIVP9:
		args = append(args, "-init_hw_device")
		args = append(args, "qsv=hw")
		args = append(args, "-filter_hw_device")
		args = append(args, "hw")
	}

	return args
}

// Initialise a video filter for HW encoding
func HWFilterInit(codec VideoCodec) VideoFilter {
	var videoFilter VideoFilter
	switch codec {
	case VideoCodecV264,
		VideoCodecVVP9:
		videoFilter = videoFilter.Append("format=nv12")
		videoFilter = videoFilter.Append("hwupload")
	case VideoCodecN264:
		videoFilter = videoFilter.Append("format=nv12")
		videoFilter = videoFilter.Append("hwupload_cuda")
	case VideoCodecI264,
		VideoCodecIVP9:
		videoFilter = videoFilter.Append("hwupload=extra_hw_frames=64")
		videoFilter = videoFilter.Append("format=qsv")
	}

	return videoFilter
}

/*
Prepend input for full hardware transcoding

Currently unused
One strategy is to use HWCodecVideoSupported and test if its supported, and then apply this instead of HWDeviceInit and HWFilterInit.
*/
func HWDeviceInit_Full(args Args, codec VideoCodec) Args {
	switch codec {
	case VideoCodecN264:
		args = append(args, "-hwaccel")
		args = append(args, "cuda")
		args = append(args, "-hwaccel_output_format")
		args = append(args, "cuda")
		args = append(args, "-hwaccel_device")
		args = append(args, "0")
	case VideoCodecV264,
		VideoCodecVVP9:
		args = append(args, "-hwaccel")
		args = append(args, "vaapi")
		args = append(args, "-hwaccel_output_format")
		args = append(args, "vaapi")
		args = append(args, "-vaapi_device")
		args = append(args, "/dev/dri/renderD128")
	case VideoCodecI264,
		VideoCodecIVP9:
		args = append(args, "-hwaccel")
		args = append(args, "qsv")
	}

	return args
}

// Replace video filter scaling with hardware scaling for full hardware transcoding
func HWCodecFilter(args VideoFilter, codec VideoCodec) VideoFilter {
	sargs := string(args)

	if strings.Contains(sargs, "scale=") {
		switch codec {
		case VideoCodecN264:
			args = VideoFilter(strings.Replace(sargs, "scale=", "scale_cuda=", 1))
		case VideoCodecV264,
			VideoCodecVVP9:
			args = VideoFilter(strings.Replace(sargs, "scale=", "scale_vaapi=", 1))
		case VideoCodecI264,
			VideoCodecIVP9:
			// BUG: [scale_qsv]: Size values less than -1 are not acceptable.
			// Fix: Replace all instances of -2 with -1 in a scale operation
			re := regexp.MustCompile(`(scale=)([\d:]*)(-2)(.*)`)
			sargs = re.ReplaceAllString(sargs, "scale=$2-1$4")
			args = VideoFilter(strings.Replace(sargs, "scale=", "scale_qsv=", 1))
		}
	}

	return args
}

func HWCodecMaxRes(codec VideoCodec, dW int, dH int) (int, int) {
	if codec == VideoCodecN264 {
		return 4096, 4096
	}

	return dW, dH
}

// Return if a hardware accelerated for HLS is available
func HWCodecHLSCompatible() *VideoCodec {
	for _, element := range HWCodecSupport {
		switch element {
		case VideoCodecN264,
			VideoCodecI264,
			VideoCodecV264:
			return &element
		}
	}
	return nil
}

// Return if a hardware accelerated codec for MP4 is available
func HWCodecMP4Compatible() *VideoCodec {
	for _, element := range HWCodecSupport {
		switch element {
		case VideoCodecN264,
			VideoCodecI264:
			return &element
		}
	}
	return nil
}

// Return if a hardware accelerated codec for WebM is available
func HWCodecWEBMCompatible() *VideoCodec {
	for _, element := range HWCodecSupport {
		switch element {
		case VideoCodecIVP9,
			VideoCodecVVP9:
			return &element
		}
	}
	return nil
}

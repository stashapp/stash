package ffmpeg

import (
	"context"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
)

var HWCodecSupport []StreamFormat

//Tests all (given) hardware codec's
func FindHWCodecs(ctx context.Context, encoder FFMpeg) {
	for _, codec := range []StreamFormat{
		StreamFormatN264,
		StreamFormatI264,
		/*
			Untested:
				StreamFormatA264,
				StreamFormatM264,
				StreamFormatV264,
				StreamFormatO264,
				StreamFormatIVP9,
				StreamFormatVVP9,
		*/
	} {
		var args Args
		args = append(args, "-hide_banner")
		args = args.LogLevel(LogLevelQuiet)
		args = HWCodecDevice_Encode(args, codec.codec)
		args = args.Format("lavfi")
		args = args.Input("color=c=red")
		args = args.Duration(0.1)

		args = HWCodecPrepend_Encode(args, codec.codec)
		args = args.VideoCodec(codec.codec)
		if len(codec.extraArgs) > 0 {
			args = append(args, codec.extraArgs...)
		}

		//Test scaling
		var videoFilter VideoFilter
		videoFilter = videoFilter.ScaleDimensions(-2, 160)
		videoFilter = HWCodecFilter(videoFilter, codec.codec)
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
		logger.Info("\t", codec.codec)
	}
}

//Return if given codec is hardware accelerated
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

//Test full-hardware transcoding on an input video
func HWCodecVideoSupported(ctx context.Context, encoder FFMpeg, o TranscodeStreamOptions) bool {
	if !HWCodecDetect(o.Codec.codec) {
		return false
	}

	var args Args
	args = append(args, "-hide_banner")
	args = append(args, o.ExtraInputArgs...)
	args = args.LogLevel(LogLevelQuiet)
	args = HWCodecDevice_Full(args, o.Codec.codec)
	args = args.Input(o.Input)
	args = args.Duration(0.1)

	//Test scaling
	var videoFilter VideoFilter
	videoFilter = videoFilter.ScaleDimensions(-2, 160)
	videoFilter = HWCodecFilter(videoFilter, o.Codec.codec)
	args = args.VideoFilter(videoFilter)

	args = args.VideoCodec(o.Codec.codec)
	if len(o.Codec.extraArgs) > 0 {
		args = append(args, o.Codec.extraArgs...)
	}

	args = args.Format("null")
	args = args.Output("-")

	cmd := encoder.Command(ctx, args)

	err := cmd.Run()
	return err == nil
}

//Prepend input for hardware encoding only
func HWCodecDevice_Encode(args Args, codec VideoCodec) Args {
	switch codec {
	case VideoCodecN264:
		args = append(args, "-hwaccel_device")
		args = append(args, "0")
	case VideoCodecV264,
		VideoCodecVVP9:
		args = append(args, "-hwaccel_device")
		args = append(args, "/dev/dri/renderD128")
	}

	return args
}

//Prepend codec for hardware encoding only
func HWCodecPrepend_Encode(args Args, codec VideoCodec) Args {
	switch codec {
	case VideoCodecV264,
		VideoCodecVVP9:
		args = append(args, "-vf")
		args = append(args, "format=nv12,hwupload")
	}

	return args
}

/*
Prepend input for full hardware transcoding

Currently unused
One strategy is to use HWCodecVideoSupported and test if its supported, and then apply this instead of the _Encode functions.
*/
func HWCodecDevice_Full(args Args, codec VideoCodec) Args {
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
		args = append(args, "-hwaccel_device")
		args = append(args, "/dev/dri/renderD128")
	case VideoCodecI264,
		VideoCodecIVP9:
		args = append(args, "-hwaccel")
		args = append(args, "qsv")
		args = append(args, "-hwaccel_device")
		args = append(args, "/dev/dri/renderD128")
	}

	return args
}

//Replace video filter scaling with hardware scaling for full hardware transcoding
func HWCodecFilter(args VideoFilter, codec VideoCodec) VideoFilter {
	sargs := string(args)

	if strings.Contains(sargs, "scale=") {
		switch codec {
		case VideoCodecN264:
			args = VideoFilter(strings.Replace(sargs, "scale=", "hwupload_cuda,scale_cuda=", 1)).Append("hwdownload")
		case VideoCodecV264,
			VideoCodecVVP9:
			args = VideoFilter(strings.Replace(sargs, "scale=", "hwupload,scale_vaapi=", 1)).Append("hwdownload")
			//BUG: scale_qsv is seemingly broken on windows?
			/*case VideoCodecI264,
			VideoCodecIVP9:
			args = VideoFilter(strings.Replace(sargs, "scale=", "hwupload,scale_qsv=", 1)).Append("hwdownload")*/
		}
	}

	return args
}

//Return if a hardware accelerated H264 codec is available
func HWCodecH264Compatible() *StreamFormat {
	for _, element := range HWCodecSupport {
		switch element.codec {
		case VideoCodecN264,
			VideoCodecA264,
			VideoCodecM264,
			VideoCodecV264,
			VideoCodecI264,
			VideoCodecR264,
			VideoCodecO264:
			return &element
		}
	}
	return nil
}

//Return if a hardware accelerated VP9 codec is available
func HWCodecVP9Compatible() *StreamFormat {
	for _, element := range HWCodecSupport {
		switch element.codec {
		case VideoCodecIVP9,
			VideoCodecVVP9:
			return &element
		}
	}
	return nil
}

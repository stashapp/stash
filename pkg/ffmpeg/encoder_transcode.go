package ffmpeg

type TranscodeOptions struct {
	OutputPath string
}

func (e *Encoder) Transcode(probeResult VideoFile, options TranscodeOptions) {
	args := []string{
		"-i", probeResult.Path,
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-profile:v", "high",
		"-level", "4.2",
		"-preset", "superfast",
		"-crf", "23",
		"-vf", "scale=iw:-2",
		"-c:a", "aac",
		"-strict", "-2",
		options.OutputPath,
	}
	_, _ = e.run(probeResult, args)
}

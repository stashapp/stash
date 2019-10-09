package ffmpeg

import "fmt"

type ScreenshotOptions struct {
	OutputPath string
	Quality    int
	Time       float64
	Width      int
	Verbosity  string
}

func (e *Encoder) Screenshot(probeResult VideoFile, options ScreenshotOptions) {
	if options.Verbosity == "" {
		options.Verbosity = "quiet"
	}
	if options.Quality == 0 {
		options.Quality = 1
	}
	args := []string{
		"-v " + options.Verbosity,
		fmt.Sprintf("-ss %v", options.Time),
		"-y",
		"-i \"" + probeResult.Path + "\"",
		"-vframes 1",
		fmt.Sprintf("-q:v %v", options.Quality),
		fmt.Sprintf("-vf scale=%v:-1", options.Width),
		"-f image2",
		options.OutputPath,
	}
	_, _ = e.run(probeResult, args)
}

// TODO: document in README.md
package main

import (
	"context"
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/hash/videophash"
	"github.com/stashapp/stash/pkg/models"
)

func customUsage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "%s [OPTIONS] VIDEOFILE...\n\nOptions:\n", os.Args[0])
	flag.PrintDefaults()
}

func printPhash(ff *ffmpeg.FFMpeg, ffp ffmpeg.FFProbe, inputfile string, quiet *bool) error {
	ffvideoFile, err := ffp.NewVideoFile(inputfile)
	if err != nil {
		return err
	}

	// All we need for videophash.Generate() is
	// videoFile.Path (from BaseFile)
	// videoFile.Duration
	// The rest of the struct isn't needed.
	vf := &models.VideoFile{
		BaseFile: &models.BaseFile{Path: inputfile},
		Duration: ffvideoFile.FileDuration,
	}

	phash, err := videophash.Generate(ff, vf)
	if err != nil {
		return err
	}

	if *quiet {
		fmt.Printf("%x\n", *phash)
	} else {
		fmt.Printf("%x %v\n", *phash, vf.Path)
	}
	return nil
}

func main() {
	flag.Usage = customUsage
	quiet := flag.BoolP("quiet", "q", false, "print only the phash")
	help := flag.BoolP("help", "h", false, "print this help output")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(2)
	}

	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Missing VIDEOFILE argument.\n")
		flag.Usage()
		os.Exit(2)
	}

	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Files will be processed sequentially! Consier using GNU Parallel.")
		fmt.Fprintf(os.Stderr, "Example: parallel %v ::: *.mp4\n", os.Args[0])
	}

	ffmpegPath, ffprobePath := ffmpeg.GetPaths(nil)
	encoder := ffmpeg.NewEncoder(ffmpegPath)
	encoder.InitHWSupport(context.TODO())
	ffprobe := ffmpeg.FFProbe(ffprobePath)

	for _, item := range args {
		if err := printPhash(encoder, ffprobe, item, quiet); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

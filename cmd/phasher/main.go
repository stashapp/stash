// TODO: document in README.md
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	flag "github.com/spf13/pflag"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/hash/imagephash"
	"github.com/stashapp/stash/pkg/hash/videophash"
	"github.com/stashapp/stash/pkg/models"
)

func customUsage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "%s [OPTIONS] FILE...\n\nOptions:\n", os.Args[0])
	flag.PrintDefaults()
}

func printPhash(ff *ffmpeg.FFMpeg, ffp *ffmpeg.FFProbe, inputfile string, quiet *bool) error {
	// Determine if this is a video or image file based on extension
	ext := filepath.Ext(inputfile)
	ext = ext[1:] // remove the leading dot

	// Common image extensions
	imageExts := map[string]bool{
		"jpg": true, "jpeg": true, "png": true, "gif": true, "webp": true, "bmp": true, "avif": true,
	}

	if imageExts[ext] {
		return printImagePhash(ff, inputfile, quiet)
	}

	return printVideoPhash(ff, ffp, inputfile, quiet)
}

func printVideoPhash(ff *ffmpeg.FFMpeg, ffp *ffmpeg.FFProbe, inputfile string, quiet *bool) error {
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

func printImagePhash(ff *ffmpeg.FFMpeg, inputfile string, quiet *bool) error {
	imgFile := &models.ImageFile{
		BaseFile: &models.BaseFile{Path: inputfile},
	}

	phash, err := imagephash.Generate(ff, imgFile)
	if err != nil {
		return err
	}

	if *quiet {
		fmt.Printf("%x\n", *phash)
	} else {
		fmt.Printf("%x %v\n", *phash, imgFile.Path)
	}
	return nil
}

func getPaths() (string, string) {
	ffmpegPath, _ := exec.LookPath("ffmpeg")
	ffprobePath, _ := exec.LookPath("ffprobe")

	return ffmpegPath, ffprobePath
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
		fmt.Fprintf(os.Stderr, "Missing FILE argument.\n")
		flag.Usage()
		os.Exit(2)
	}

	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Files will be processed sequentially! If required, use e.g. GNU Parallel to run concurrently.")
		fmt.Fprintf(os.Stderr, "Example: parallel %v ::: *.mp4\n", os.Args[0])
	}

	ffmpegPath, ffprobePath := getPaths()
	encoder := ffmpeg.NewEncoder(ffmpegPath)
	// don't need to InitHWSupport, phashing doesn't use hw acceleration
	ffprobe := ffmpeg.NewFFProbe(ffprobePath)

	for _, item := range args {
		if err := printPhash(encoder, ffprobe, item, quiet); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

}

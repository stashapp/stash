// TODO: document in README.md
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/hash/videophash"
)

func main() {
	// TODO: Tidier argument handling (spf13/cobra?) and print usage info.
	if len(os.Args) < 2 {
		panic("missing argument")
	}
	inputfile := os.Args[1]

	ffmpegPath, ffprobePath := ffmpeg.GetPaths(nil)
	FFMPEG := ffmpeg.NewEncoder(ffmpegPath)
	FFMPEG.InitHWSupport(context.TODO())

	FFPROBE := ffmpeg.FFProbe(ffprobePath)
	ffvideoFile, err := FFPROBE.NewVideoFile(inputfile)
	if err != nil {
		fmt.Println(err)
	}

	// All we need for videophash.Generate() is
	// videoFile.Path (from BaseFile)
	// videoFile.Duration
	// The rest of the struct isn't needed.
	vf := &file.VideoFile{
		BaseFile: &file.BaseFile{Path: inputfile},
		Duration: ffvideoFile.FileDuration,
	}

	phash, err := videophash.Generate(FFMPEG, vf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%x %v\n", *phash, vf.Path)
}

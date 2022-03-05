package encoder

import (
	"bytes"
	"context"
	"fmt"
	"image"

	"github.com/stashapp/stash/pkg/ffmpeg2"
	"github.com/stashapp/stash/pkg/video"
)

func doGenerate(encoder ffmpeg2.FFMpeg, fn string, args ffmpeg2.Args) error {
	ctx, cancel := readLockManager.ReadLock(context.Background(), fn)
	defer cancel()
	return video.Generate(encoder, ctx, args)
}

func doGenerateOutput(encoder ffmpeg2.FFMpeg, fn string, args ffmpeg2.Args) ([]byte, error) {
	ctx, cancel := readLockManager.ReadLock(context.Background(), fn)
	defer cancel()
	return video.GenerateOutput(encoder, ctx, args)
}

func doGenerateImage(encoder ffmpeg2.FFMpeg, input string, args ffmpeg2.Args) (image.Image, error) {
	out, err := doGenerateOutput(encoder, input, args)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(out))
	if err != nil {
		return nil, fmt.Errorf("decoding image from ffmpeg: %w", err)
	}

	return img, nil
}

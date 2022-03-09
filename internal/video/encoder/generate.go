package encoder

import (
	"bytes"
	"context"
	"fmt"
	"image"

	"github.com/stashapp/stash/pkg/ffmpeg"
)

func doGenerate(encoder ffmpeg.FFMpeg, fn string, args ffmpeg.Args) error {
	ctx, cancel := readLockManager.ReadLock(context.Background(), fn)
	defer cancel()
	return ffmpeg.Generate(ctx, encoder, args)
}

func doGenerateOutput(encoder ffmpeg.FFMpeg, fn string, args ffmpeg.Args) ([]byte, error) {
	ctx, cancel := readLockManager.ReadLock(context.Background(), fn)
	defer cancel()
	return ffmpeg.GenerateOutput(ctx, encoder, args)
}

func doGenerateImage(encoder ffmpeg.FFMpeg, input string, args ffmpeg.Args) (image.Image, error) {
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

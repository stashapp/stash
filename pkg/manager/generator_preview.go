package manager

import (
	"bufio"
	"fmt"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
	"os"
	"path/filepath"
)

const previewImageWidth = 640

type PreviewGenerator struct {
	Info *GeneratorInfo

	VideoFilename   string
	ImageFilename   string
	OutputDirectory string
}

func NewPreviewGenerator(videoFile ffmpeg.VideoFile, videoFilename string, imageFilename string, outputDirectory string) (*PreviewGenerator, error) {
	exists, err := utils.FileExists(videoFile.Path)
	if !exists {
		return nil, err
	}
	generator, err := newGeneratorInfo(videoFile)
	if err != nil {
		return nil, err
	}
	generator.ChunkCount = 12 // 12 segments to the preview
	if err := generator.configure(); err != nil {
		return nil, err
	}

	return &PreviewGenerator{
		Info:            generator,
		VideoFilename:   videoFilename,
		ImageFilename:   imageFilename,
		OutputDirectory: outputDirectory,
	}, nil
}

func (g *PreviewGenerator) Generate() error {
	logger.Infof("[generator] generating scene preview for %s", g.Info.VideoFile.Path)
	encoder := ffmpeg.NewEncoder(instance.FFMPEGPath)

	if err := g.generateConcatFile(); err != nil {
		return err
	}
	if err := g.generateVideo(&encoder); err != nil {
		return err
	}
	if err := g.generateDefaultImage(&encoder, false); err != nil {
		return err
	}
	return nil
}

func (g *PreviewGenerator) GenerateDefaultImage() error {
	logger.Infof("[generator] generating scene preview for %s", g.Info.VideoFile.Path)
	encoder := ffmpeg.NewEncoder(instance.FFMPEGPath)
	return g.generateDefaultImage(&encoder, true)
}

func (g *PreviewGenerator) GenerateImageAt(at float64) error {
	logger.Infof("[generator] generating scene preview for %s", g.Info.VideoFile.Path)
	encoder := ffmpeg.NewEncoder(instance.FFMPEGPath)
	return g.generateImageAt(&encoder, at)
}

func (g *PreviewGenerator) generateConcatFile() error {
	f, err := os.Create(g.getConcatFilePath())
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for i := 0; i < g.Info.ChunkCount; i++ {
		num := fmt.Sprintf("%.3d", i)
		filename := "preview" + num + ".mp4"
		_, _ = w.WriteString(fmt.Sprintf("file '%s'\n", filename))
	}
	return w.Flush()
}

func (g *PreviewGenerator) generateVideo(encoder *ffmpeg.Encoder) error {
	outputPath := filepath.Join(g.OutputDirectory, g.VideoFilename)
	outputExists, _ := utils.FileExists(outputPath)
	if outputExists {
		return nil
	}

	stepSize := int(g.Info.VideoFile.Duration / float64(g.Info.ChunkCount))
	for i := 0; i < g.Info.ChunkCount; i++ {
		time := i * stepSize
		num := fmt.Sprintf("%.3d", i)
		filename := "preview" + num + ".mp4"
		chunkOutputPath := instance.Paths.Generated.GetTmpPath(filename)

		options := ffmpeg.ScenePreviewChunkOptions{
			Time:       time,
			Width:      640,
			OutputPath: chunkOutputPath,
		}
		encoder.ScenePreviewVideoChunk(g.Info.VideoFile, options)
	}

	videoOutputPath := filepath.Join(g.OutputDirectory, g.VideoFilename)
	encoder.ScenePreviewVideoChunkCombine(g.Info.VideoFile, g.getConcatFilePath(), videoOutputPath)
	logger.Debug("created video preview: ", videoOutputPath)
	return nil
}

func (g *PreviewGenerator) generateDefaultImage(encoder *ffmpeg.Encoder, overwrite bool) error {
	outputPath := filepath.Join(g.OutputDirectory, g.ImageFilename)
	outputExists, _ := utils.FileExists(outputPath)
	if !overwrite && outputExists {
		return nil
	}

	videoPreviewPath := filepath.Join(g.OutputDirectory, g.VideoFilename)
	tmpOutputPath := instance.Paths.Generated.GetTmpPath(g.ImageFilename)
	if err := encoder.ScenePreviewVideoToImage(g.Info.VideoFile, previewImageWidth, videoPreviewPath, tmpOutputPath); err != nil {
		return err
	}
	if err := os.Rename(tmpOutputPath, outputPath); err != nil {
		return err
	}
	logger.Debug("created video preview image: ", outputPath)

	return nil
}

func (g *PreviewGenerator) generateImageAt(encoder *ffmpeg.Encoder, at float64) error {
	outputPath := filepath.Join(g.OutputDirectory, g.ImageFilename)

	// assume always overwrite

	tmpOutputPath := instance.Paths.Generated.GetTmpPath(g.ImageFilename)
	options := ffmpeg.ScreenshotOptions{
		OutputPath: tmpOutputPath,
		Time:       at,
		Width:      previewImageWidth,
	}

	if err := encoder.Screenshot(g.Info.VideoFile, options); err != nil {
		return err
	}

	if err := os.Rename(tmpOutputPath, outputPath); err != nil {
		return err
	}

	logger.Debug("created video preview image: ", outputPath)

	return nil
}

func (g *PreviewGenerator) getConcatFilePath() string {
	return instance.Paths.Generated.GetTmpPath("files.txt")
}

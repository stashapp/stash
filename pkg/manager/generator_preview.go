package manager

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
)

type PreviewGenerator struct {
	Info *GeneratorInfo

	VideoChecksum   string
	VideoFilename   string
	ImageFilename   string
	OutputDirectory string

	GenerateVideo bool
	GenerateImage bool

	PreviewPreset string

	Overwrite bool
}

func NewPreviewGenerator(videoFile ffmpeg.VideoFile, videoChecksum string, videoFilename string, imageFilename string, outputDirectory string, generateVideo bool, generateImage bool, previewPreset string) (*PreviewGenerator, error) {
	exists, err := utils.FileExists(videoFile.Path)
	if !exists {
		return nil, err
	}
	generator, err := newGeneratorInfo(videoFile)
	if err != nil {
		return nil, err
	}
	generator.ChunkCount = 12 // 12 segments to the preview

	return &PreviewGenerator{
		Info:            generator,
		VideoChecksum:   videoChecksum,
		VideoFilename:   videoFilename,
		ImageFilename:   imageFilename,
		OutputDirectory: outputDirectory,
		GenerateVideo:   generateVideo,
		GenerateImage:   generateImage,
		PreviewPreset:   previewPreset,
	}, nil
}

func (g *PreviewGenerator) Generate() error {
	logger.Infof("[generator] generating scene preview for %s", g.Info.VideoFile.Path)

	if err := g.Info.configure(); err != nil {
		return err
	}

	encoder := ffmpeg.NewEncoder(instance.FFMPEGPath)

	tmpFiles, err := g.generateConcatFile()
	if err != nil {
		return err
	}

	if g.GenerateVideo {
		if err := g.generateVideo(&encoder, false); err != nil {
			logger.Warnf("[generator] failed generating scene preview, trying fallback")
			if err := g.generateVideo(&encoder, true); err != nil {
				//remove intermediate files
				removeFiles(tmpFiles)
				return err
			}
		}
		//remove intermediate files
		removeFiles(tmpFiles)
	}
	if g.GenerateImage {
		if err := g.generateImage(&encoder); err != nil {
			return err
		}
	}
	return nil
}

func (g *PreviewGenerator) generateConcatFile() (*[]string, error) {
	f, err := os.Create(g.getConcatFilePath())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var fileList []string
	fileList = append(fileList, g.getConcatFilePath()) // add concat file to fileList

	w := bufio.NewWriter(f)
	for i := 0; i < g.Info.ChunkCount; i++ {
		num := fmt.Sprintf("%.3d", i)
		filename := "preview_" + g.VideoChecksum + "_" + num + ".mp4"
		_, _ = w.WriteString(fmt.Sprintf("file '%s'\n", filename))
		filename = instance.Paths.Generated.GetTmpPath(filename) // make filename absolute
		fileList = append(fileList, filename)                    // and add it to fileList
	}
	err = w.Flush()
	return &fileList, err
}

func (g *PreviewGenerator) generateVideo(encoder *ffmpeg.Encoder, fallback bool) error {
	outputPath := filepath.Join(g.OutputDirectory, g.VideoFilename)
	outputExists, _ := utils.FileExists(outputPath)
	if !g.Overwrite && outputExists {
		return nil
	}

	stepSize, offset := g.Info.getStepSizeAndOffset()

	durationSegment := g.Info.ChunkDuration
	if durationSegment < 0.75 { // a very short duration can create files without a video stream
		durationSegment = 0.75 // use 0.75 in that case
		logger.Warnf("[generator] Segment duration (%f) too small.Using 0.75 instead.", g.Info.ChunkDuration)
	}

	for i := 0; i < g.Info.ChunkCount; i++ {
		time := offset + (float64(i) * stepSize)
		num := fmt.Sprintf("%.3d", i)
		filename := "preview_" + g.VideoChecksum + "_" + num + ".mp4"
		chunkOutputPath := instance.Paths.Generated.GetTmpPath(filename)

		options := ffmpeg.ScenePreviewChunkOptions{
			StartTime:  time,
			Duration:   durationSegment,
			Width:      640,
			OutputPath: chunkOutputPath,
		}
		if err := encoder.ScenePreviewVideoChunk(g.Info.VideoFile, options, g.PreviewPreset, fallback); err != nil {
			return err
		}
	}

	videoOutputPath := filepath.Join(g.OutputDirectory, g.VideoFilename)
	if err := encoder.ScenePreviewVideoChunkCombine(g.Info.VideoFile, g.getConcatFilePath(), videoOutputPath); err != nil {
		return err
	}
	logger.Debug("created video preview: ", videoOutputPath)
	return nil
}

func (g *PreviewGenerator) generateImage(encoder *ffmpeg.Encoder) error {
	outputPath := filepath.Join(g.OutputDirectory, g.ImageFilename)
	outputExists, _ := utils.FileExists(outputPath)
	if !g.Overwrite && outputExists {
		return nil
	}

	videoPreviewPath := filepath.Join(g.OutputDirectory, g.VideoFilename)
	tmpOutputPath := instance.Paths.Generated.GetTmpPath(g.ImageFilename)
	if err := encoder.ScenePreviewVideoToImage(g.Info.VideoFile, 640, videoPreviewPath, tmpOutputPath); err != nil {
		return err
	}
	if err := utils.SafeMove(tmpOutputPath, outputPath); err != nil {
		return err
	}
	logger.Debug("created video preview image: ", outputPath)

	return nil
}

func (g *PreviewGenerator) getConcatFilePath() string {
	return instance.Paths.Generated.GetTmpPath(fmt.Sprintf("files_%s.txt", g.VideoChecksum))
}

func removeFiles(list *[]string) {
	if list != nil {
		for _, f := range *list {
			if err := os.Remove(f); err != nil {
				logger.Warnf("[generator] Delete error: %s", err)
			}
		}
	}

}

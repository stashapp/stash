package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

// DestroyScene deletes a scene and its associated relationships from the
// database. Returns a function to perform any post-commit actions.
func DestroyScene(scene *models.Scene, repo models.Repository) (func(), error) {
	qb := repo.Scene()
	mqb := repo.SceneMarker()

	markers, err := mqb.FindBySceneID(scene.ID)
	if err != nil {
		return nil, err
	}

	var funcs []func()
	for _, m := range markers {
		f, err := DestroySceneMarker(scene, m, mqb)
		if err != nil {
			return nil, err
		}
		funcs = append(funcs, f)
	}

	if err := qb.Destroy(scene.ID); err != nil {
		return nil, err
	}

	return func() {
		for _, f := range funcs {
			f()
		}
	}, nil
}

// DestroySceneMarker deletes the scene marker from the database and returns a
// function that removes the generated files, to be executed after the
// transaction is successfully committed.
func DestroySceneMarker(scene *models.Scene, sceneMarker *models.SceneMarker, qb models.SceneMarkerWriter) (func(), error) {
	if err := qb.Destroy(sceneMarker.ID); err != nil {
		return nil, err
	}

	// delete the preview for the marker
	return func() {
		seconds := int(sceneMarker.Seconds)
		DeleteSceneMarkerFiles(scene, seconds, config.GetInstance().GetVideoFileNamingAlgorithm())
	}, nil
}

// DeleteGeneratedSceneFiles deletes generated files for the provided scene.
func DeleteGeneratedSceneFiles(scene *models.Scene, fileNamingAlgo models.HashAlgorithm) {
	sceneHash := scene.GetHash(fileNamingAlgo)

	if sceneHash == "" {
		return
	}

	markersFolder := filepath.Join(GetInstance().Paths.Generated.Markers, sceneHash)

	exists, _ := utils.FileExists(markersFolder)
	if exists {
		err := os.RemoveAll(markersFolder)
		if err != nil {
			logger.Warnf("Could not delete folder %s: %s", markersFolder, err.Error())
		}
	}

	thumbPath := GetInstance().Paths.Scene.GetThumbnailScreenshotPath(sceneHash)
	exists, _ = utils.FileExists(thumbPath)
	if exists {
		err := os.Remove(thumbPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", thumbPath, err.Error())
		}
	}

	normalPath := GetInstance().Paths.Scene.GetScreenshotPath(sceneHash)
	exists, _ = utils.FileExists(normalPath)
	if exists {
		err := os.Remove(normalPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", normalPath, err.Error())
		}
	}

	streamPreviewPath := GetInstance().Paths.Scene.GetStreamPreviewPath(sceneHash)
	exists, _ = utils.FileExists(streamPreviewPath)
	if exists {
		err := os.Remove(streamPreviewPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", streamPreviewPath, err.Error())
		}
	}

	streamPreviewImagePath := GetInstance().Paths.Scene.GetStreamPreviewImagePath(sceneHash)
	exists, _ = utils.FileExists(streamPreviewImagePath)
	if exists {
		err := os.Remove(streamPreviewImagePath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", streamPreviewImagePath, err.Error())
		}
	}

	transcodePath := GetInstance().Paths.Scene.GetTranscodePath(sceneHash)
	exists, _ = utils.FileExists(transcodePath)
	if exists {
		// kill any running streams
		KillRunningStreams(transcodePath)

		err := os.Remove(transcodePath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", transcodePath, err.Error())
		}
	}

	spritePath := GetInstance().Paths.Scene.GetSpriteImageFilePath(sceneHash)
	exists, _ = utils.FileExists(spritePath)
	if exists {
		err := os.Remove(spritePath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", spritePath, err.Error())
		}
	}

	vttPath := GetInstance().Paths.Scene.GetSpriteVttFilePath(sceneHash)
	exists, _ = utils.FileExists(vttPath)
	if exists {
		err := os.Remove(vttPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", vttPath, err.Error())
		}
	}
}

// DeleteSceneMarkerFiles deletes generated files for a scene marker with the
// provided scene and timestamp.
func DeleteSceneMarkerFiles(scene *models.Scene, seconds int, fileNamingAlgo models.HashAlgorithm) {
	videoPath := GetInstance().Paths.SceneMarkers.GetStreamPath(scene.GetHash(fileNamingAlgo), seconds)
	imagePath := GetInstance().Paths.SceneMarkers.GetStreamPreviewImagePath(scene.GetHash(fileNamingAlgo), seconds)
	screenshotPath := GetInstance().Paths.SceneMarkers.GetStreamScreenshotPath(scene.GetHash(fileNamingAlgo), seconds)

	exists, _ := utils.FileExists(videoPath)
	if exists {
		err := os.Remove(videoPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", videoPath, err.Error())
		}
	}

	exists, _ = utils.FileExists(imagePath)
	if exists {
		err := os.Remove(imagePath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", imagePath, err.Error())
		}
	}

	exists, _ = utils.FileExists(screenshotPath)
	if exists {
		err := os.Remove(screenshotPath)
		if err != nil {
			logger.Warnf("Could not delete file %s: %s", screenshotPath, err.Error())
		}
	}
}

// DeleteSceneFile deletes the scene video file from the filesystem.
func DeleteSceneFile(scene *models.Scene) {
	// kill any running encoders
	KillRunningStreams(scene.Path)

	err := os.Remove(scene.Path)
	if err != nil {
		logger.Warnf("Could not delete file %s: %s", scene.Path, err.Error())
	}
}

func GetSceneFileContainer(scene *models.Scene) (ffmpeg.Container, error) {
	var container ffmpeg.Container
	if scene.Format.Valid {
		container = ffmpeg.Container(scene.Format.String)
	} else { // container isn't in the DB
		// shouldn't happen, fallback to ffprobe
		ffprobe := GetInstance().FFProbe
		tmpVideoFile, err := ffprobe.NewVideoFile(scene.Path, false)
		if err != nil {
			return ffmpeg.Container(""), fmt.Errorf("error reading video file: %v", err)
		}

		container = ffmpeg.MatchContainer(tmpVideoFile.Container, scene.Path)
	}

	return container, nil
}

func includeSceneStreamPath(scene *models.Scene, streamingResolution models.StreamingResolutionEnum, maxStreamingTranscodeSize models.StreamingResolutionEnum) bool {
	// convert StreamingResolutionEnum to ResolutionEnum so we can get the min
	// resolution
	convertedRes := models.ResolutionEnum(streamingResolution)

	minResolution := int64(convertedRes.GetMinResolution())
	sceneResolution := scene.GetMinResolution()

	// don't include if scene resolution is smaller than the streamingResolution
	if sceneResolution != 0 && sceneResolution < minResolution {
		return false
	}

	// if we always allow everything, then return true
	if maxStreamingTranscodeSize == models.StreamingResolutionEnumOriginal {
		return true
	}

	// convert StreamingResolutionEnum to ResolutionEnum
	maxStreamingResolution := models.ResolutionEnum(maxStreamingTranscodeSize)
	return int64(maxStreamingResolution.GetMinResolution()) >= minResolution
}

func makeStreamEndpoint(streamURL string, streamingResolution models.StreamingResolutionEnum, mimeType, label string) *models.SceneStreamEndpoint {
	return &models.SceneStreamEndpoint{
		URL:      fmt.Sprintf("%s?resolution=%s", streamURL, streamingResolution.String()),
		MimeType: &mimeType,
		Label:    &label,
	}
}

func GetSceneStreamPaths(scene *models.Scene, directStreamURL string, maxStreamingTranscodeSize models.StreamingResolutionEnum) ([]*models.SceneStreamEndpoint, error) {
	if scene == nil {
		return nil, fmt.Errorf("nil scene")
	}

	var ret []*models.SceneStreamEndpoint
	mimeWebm := ffmpeg.MimeWebm
	mimeHLS := ffmpeg.MimeHLS
	mimeMp4 := ffmpeg.MimeMp4

	labelWebm := "webm"
	labelHLS := "HLS"

	// direct stream should only apply when the audio codec is supported
	audioCodec := ffmpeg.MissingUnsupported
	if scene.AudioCodec.Valid {
		audioCodec = ffmpeg.AudioCodec(scene.AudioCodec.String)
	}

	// don't care if we can't get the container
	container, _ := GetSceneFileContainer(scene)

	if HasTranscode(scene, config.GetInstance().GetVideoFileNamingAlgorithm()) || ffmpeg.IsValidAudioForContainer(audioCodec, container) {
		label := "Direct stream"
		ret = append(ret, &models.SceneStreamEndpoint{
			URL:      directStreamURL,
			MimeType: &mimeMp4,
			Label:    &label,
		})
	}

	// only add mkv stream endpoint if the scene container is an mkv already
	if container == ffmpeg.Matroska {
		label := "mkv"
		ret = append(ret, &models.SceneStreamEndpoint{
			URL: directStreamURL + ".mkv",
			// set mkv to mp4 to trick the client, since many clients won't try mkv
			MimeType: &mimeMp4,
			Label:    &label,
		})
	}

	hls := models.SceneStreamEndpoint{
		URL:      directStreamURL + ".m3u8",
		MimeType: &mimeHLS,
		Label:    &labelHLS,
	}
	ret = append(ret, &hls)

	// WEBM quality transcoding options
	// Note: These have the wrong mime type intentionally to allow jwplayer to selection between mp4/webm
	webmLabelFourK := "WEBM 4K (2160p)"         // "FOUR_K"
	webmLabelFullHD := "WEBM Full HD (1080p)"   // "FULL_HD"
	webmLabelStandardHD := "WEBM HD (720p)"     // "STANDARD_HD"
	webmLabelStandard := "WEBM Standard (480p)" // "STANDARD"
	webmLabelLow := "WEBM Low (240p)"           // "LOW"

	// Setup up lower quality transcoding options (MP4)
	mp4LabelFourK := "MP4 4K (2160p)"         // "FOUR_K"
	mp4LabelFullHD := "MP4 Full HD (1080p)"   // "FULL_HD"
	mp4LabelStandardHD := "MP4 HD (720p)"     // "STANDARD_HD"
	mp4LabelStandard := "MP4 Standard (480p)" // "STANDARD"
	mp4LabelLow := "MP4 Low (240p)"           // "LOW"

	var webmStreams []*models.SceneStreamEndpoint
	var mp4Streams []*models.SceneStreamEndpoint

	webmURL := directStreamURL + ".webm"
	mp4URL := directStreamURL + ".mp4"

	if includeSceneStreamPath(scene, models.StreamingResolutionEnumFourK, maxStreamingTranscodeSize) {
		webmStreams = append(webmStreams, makeStreamEndpoint(webmURL, models.StreamingResolutionEnumFourK, mimeMp4, webmLabelFourK))
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp4URL, models.StreamingResolutionEnumFourK, mimeMp4, mp4LabelFourK))
	}

	if includeSceneStreamPath(scene, models.StreamingResolutionEnumFullHd, maxStreamingTranscodeSize) {
		webmStreams = append(webmStreams, makeStreamEndpoint(webmURL, models.StreamingResolutionEnumFullHd, mimeMp4, webmLabelFullHD))
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp4URL, models.StreamingResolutionEnumFullHd, mimeMp4, mp4LabelFullHD))
	}

	if includeSceneStreamPath(scene, models.StreamingResolutionEnumStandardHd, maxStreamingTranscodeSize) {
		webmStreams = append(webmStreams, makeStreamEndpoint(webmURL, models.StreamingResolutionEnumStandardHd, mimeMp4, webmLabelStandardHD))
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp4URL, models.StreamingResolutionEnumStandardHd, mimeMp4, mp4LabelStandardHD))
	}

	if includeSceneStreamPath(scene, models.StreamingResolutionEnumStandard, maxStreamingTranscodeSize) {
		webmStreams = append(webmStreams, makeStreamEndpoint(webmURL, models.StreamingResolutionEnumStandard, mimeMp4, webmLabelStandard))
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp4URL, models.StreamingResolutionEnumStandard, mimeMp4, mp4LabelStandard))
	}

	if includeSceneStreamPath(scene, models.StreamingResolutionEnumLow, maxStreamingTranscodeSize) {
		webmStreams = append(webmStreams, makeStreamEndpoint(webmURL, models.StreamingResolutionEnumLow, mimeMp4, webmLabelLow))
		mp4Streams = append(mp4Streams, makeStreamEndpoint(mp4URL, models.StreamingResolutionEnumLow, mimeMp4, mp4LabelLow))
	}

	ret = append(ret, webmStreams...)
	ret = append(ret, mp4Streams...)

	defaultStreams := []*models.SceneStreamEndpoint{
		{
			URL:      directStreamURL + ".webm",
			MimeType: &mimeWebm,
			Label:    &labelWebm,
		},
	}

	ret = append(ret, defaultStreams...)

	return ret, nil
}

// HasTranscode returns true if a transcoded video exists for the provided
// scene. It will check using the OSHash of the scene first, then fall back
// to the checksum.
func HasTranscode(scene *models.Scene, fileNamingAlgo models.HashAlgorithm) bool {
	if scene == nil {
		return false
	}

	sceneHash := scene.GetHash(fileNamingAlgo)
	if sceneHash == "" {
		return false
	}

	transcodePath := instance.Paths.Scene.GetTranscodePath(sceneHash)
	ret, _ := utils.FileExists(transcodePath)
	return ret
}

package heresphere

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

/*
 * Returns the primary media source
 */
func getPrimaryMediaSource(rs Routes, r *http.Request, scene *models.Scene) HeresphereVideoMediaSource {
	mediaFile := scene.Files.Primary()
	if mediaFile == nil {
		return HeresphereVideoMediaSource{} // Return empty source if no primary file
	}

	sourceUrl := urlbuilders.NewSceneURLBuilder(manager.GetBaseURL(r), scene).GetStreamURL("").String()
	sourceUrlWithApiKey := addApiKey(sourceUrl)

	return HeresphereVideoMediaSource{
		Resolution: mediaFile.Height,
		Height:     mediaFile.Height,
		Width:      mediaFile.Width,
		Size:       mediaFile.Size,
		Url:        sourceUrlWithApiKey,
	}
}

/*
 * This auxiliary function gathers a script if applicable
 */
func getVideoScripts(rs Routes, r *http.Request, scene *models.Scene) []HeresphereVideoScript {
	processedScripts := []HeresphereVideoScript{}

	primaryFile := scene.Files.Primary()
	if primaryFile != nil && primaryFile.Interactive {
		processedScript := HeresphereVideoScript{
			Name:   "Default script",
			Url:    addApiKey(urlbuilders.NewSceneURLBuilder(manager.GetBaseURL(r), scene).GetFunscriptURL()),
			Rating: 5,
		}
		processedScripts = append(processedScripts, processedScript)
	}

	return processedScripts
}

/*
 * This auxiliary function gathers subtitles if applicable
 */
func getVideoSubtitles(rs Routes, r *http.Request, scene *models.Scene) []HeresphereVideoSubtitle {
	processedSubtitles := []HeresphereVideoSubtitle{}

	txn.WithReadTxn(r.Context(), rs.TxnManager, func(ctx context.Context) error {
		var err error

		primaryFile := scene.Files.Primary()
		if primaryFile != nil {
			if captions_id, err := rs.Repository.File.GetCaptions(ctx, primaryFile.ID); err == nil {
				for _, caption := range captions_id {
					processedCaption := HeresphereVideoSubtitle{
						Name:     caption.Filename,
						Language: caption.LanguageCode,
						Url: addApiKey(fmt.Sprintf("%s?lang=%s&type=%s",
							urlbuilders.NewSceneURLBuilder(manager.GetBaseURL(r), scene).GetCaptionURL(),
							caption.LanguageCode,
							caption.CaptionType,
						)),
					}
					processedSubtitles = append(processedSubtitles, processedCaption)
				}
			}
		}

		return err
	})

	return processedSubtitles
}

/*
 * Function to get transcoded media sources
 */
func getTranscodedMediaSources(sceneURL string, transcodeSize int, mediaFile *file.VideoFile) map[string][]HeresphereVideoMediaSource {
	transcodedSources := make(map[string][]HeresphereVideoMediaSource)
	transNames := []string{"HLS", "DASH"}
	resRatio := float32(mediaFile.Width) / float32(mediaFile.Height)

	for i, trans := range []string{".m3u8", ".mpd"} {
		for _, res := range models.AllStreamingResolutionEnum {
			if transcodeSize == 0 || transcodeSize >= res.GetMaxResolution() {
				if height := res.GetMaxResolution(); height <= mediaFile.Height {
					transcodedUrl, err := url.Parse(sceneURL + trans)
					if err != nil {
						panic(err)
					}
					transcodedUrl.Query().Add("resolution", res.String())

					processedEntry := HeresphereVideoMediaSource{
						Resolution: height,
						Height:     height,
						Width:      int(resRatio * float32(height)),
						Size:       0,
						Url:        transcodedUrl.String(),
					}

					typeName := transNames[i]
					transcodedSources[typeName] = append(transcodedSources[typeName], processedEntry)
				}
			}
		}
	}

	return transcodedSources
}

/*
 * Main function to gather media information and transcoding options
 */
func getVideoMedia(rs Routes, r *http.Request, scene *models.Scene) []HeresphereVideoMedia {
	processedMedia := []HeresphereVideoMedia{}

	if err := txn.WithTxn(r.Context(), rs.Repository.TxnManager, func(ctx context.Context) error {
		return scene.LoadPrimaryFile(ctx, rs.Repository.File)
	}); err != nil {
		return processedMedia
	}

	primarySource := getPrimaryMediaSource(rs, r, scene)
	if primarySource.Url != "" {
		processedMedia = append(processedMedia, HeresphereVideoMedia{
			Name:    "direct stream",
			Sources: []HeresphereVideoMediaSource{primarySource},
		})
	}

	mediaFile := scene.Files.Primary()
	if mediaFile != nil {
		sceneURL := urlbuilders.NewSceneURLBuilder(manager.GetBaseURL(r), scene).GetStreamURL(config.GetInstance().GetAPIKey()).String()
		transcodeSize := config.GetInstance().GetMaxStreamingTranscodeSize().GetMaxResolution()
		transcodedSources := getTranscodedMediaSources(sceneURL, transcodeSize, mediaFile)

		// Reconstruct tables for transcoded sources
		for codec, sources := range transcodedSources {
			processedMedia = append(processedMedia, HeresphereVideoMedia{
				Name:    codec,
				Sources: sources,
			})
		}
	}

	return processedMedia
}

package heresphere

import (
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

/*
 * Reads relevant VR tags and sets projection settings
 */
func findProjectionTagsFromTags(processedScene *HeresphereVideoEntry, tags []HeresphereVideoTag) {
	for _, tag := range tags {
		tagPre := strings.TrimPrefix(tag.Name, "Tag:")

		// Has degrees tag
		if strings.HasSuffix(tagPre, "°") {
			deg := strings.TrimSuffix(tagPre, "°")
			if s, err := strconv.ParseFloat(deg, 64); err == nil {
				processedScene.Fov = s
			}
		}
		// Has VR tag
		vrTag, err := getVrTag()
		if err == nil && tagPre == vrTag {
			if processedScene.Projection == HeresphereProjectionPerspective {
				processedScene.Projection = HeresphereProjectionEquirectangular
			}
			if processedScene.Stereo == HeresphereStereoMono {
				processedScene.Stereo = HeresphereStereoSbs
			}
		}
		// Has Fisheye tag
		if tagPre == "Fisheye" {
			processedScene.Projection = HeresphereProjectionFisheye
			if processedScene.Stereo == HeresphereStereoMono {
				processedScene.Stereo = HeresphereStereoSbs
			}
		}
	}
}

/*
 * Reads relevant VR strings from a filename and sets projection settings
 */
func findProjectionTagsFromFilename(processedScene *HeresphereVideoEntry, filename string) {
	path := strings.ToUpper(filename)

	// Stereo settings
	if strings.Contains(path, "_LR") || strings.Contains(path, "_3DH") {
		processedScene.Stereo = HeresphereStereoSbs
	}
	if strings.Contains(path, "_RL") {
		processedScene.Stereo = HeresphereStereoSbs
		processedScene.IsEyeSwapped = true
	}
	if strings.Contains(path, "_TB") || strings.Contains(path, "_3DV") {
		processedScene.Stereo = HeresphereStereoTB
	}
	if strings.Contains(path, "_BT") {
		processedScene.Stereo = HeresphereStereoTB
		processedScene.IsEyeSwapped = true
	}

	// Projection settings
	if strings.Contains(path, "_EAC360") || strings.Contains(path, "_360EAC") {
		processedScene.Projection = HeresphereProjectionEquirectangularCubemap
		processedScene.Fov = 360.0
	}
	if strings.Contains(path, "_360") {
		processedScene.Projection = HeresphereProjectionEquirectangular360
		processedScene.Fov = 360.0
	}
	if strings.Contains(path, "_F180") || strings.Contains(path, "_180F") || strings.Contains(path, "_VR180") {
		processedScene.Projection = HeresphereProjectionFisheye
		processedScene.Fov = 180.0
	} else if strings.Contains(path, "_180") {
		processedScene.Projection = HeresphereProjectionEquirectangular
		processedScene.Fov = 180.0
	}
	if strings.Contains(path, "_MKX200") {
		processedScene.Projection = HeresphereProjectionFisheye
		processedScene.Fov = 200.0
		processedScene.Lens = HeresphereLensMKX200
	}
	if strings.Contains(path, "_MKX220") {
		processedScene.Projection = HeresphereProjectionFisheye
		processedScene.Fov = 220.0
		processedScene.Lens = HeresphereLensMKX220
	}
	if strings.Contains(path, "_FISHEYE") {
		processedScene.Projection = HeresphereProjectionFisheye
	}
	if strings.Contains(path, "_RF52") || strings.Contains(path, "_FISHEYE190") {
		processedScene.Projection = HeresphereProjectionFisheye
		processedScene.Fov = 190.0
	}
	if strings.Contains(path, "_VRCA220") {
		processedScene.Projection = HeresphereProjectionFisheye
		processedScene.Fov = 220.0
		processedScene.Lens = HeresphereLensVRCA220
	}
}

/*
 * This auxiliary function finds vr projection modes from tags and the filename.
 */
func FindProjectionTags(scene *models.Scene, processedScene *HeresphereVideoEntry) {
	findProjectionTagsFromTags(processedScene, processedScene.Tags)

	file := scene.Files.Primary()
	if file != nil {
		findProjectionTagsFromFilename(processedScene, file.Basename)
	}
}

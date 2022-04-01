package scene

import (
	"path/filepath"
	"strings"
)

func GetCaptionDEPath(path string) string {
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)
	return fn + ".de.vtt"
}

func GetCaptionENPath(path string) string {
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)
	return fn + ".en.vtt"
}

func GetCaptionESPath(path string) string {
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)
	return fn + ".es.vtt"
}

func GetCaptionFRPath(path string) string {
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)
	return fn + ".fr.vtt"
}

func GetCaptionITPath(path string) string {
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)
	return fn + ".it.vtt"
}

func GetCaptionNLPath(path string) string {
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)
	return fn + ".nl.vtt"
}

func GetCaptionPTPath(path string) string {
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)
	return fn + ".pt.vtt"
}

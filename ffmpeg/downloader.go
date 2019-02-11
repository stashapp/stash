package ffmpeg

import (
	"archive/zip"
	"fmt"
	"github.com/stashapp/stash/utils"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetPaths(configDirectory string) (string, string) {
	var ffmpegPath, ffprobePath string

	// Check if ffmpeg exists in the PATH
	if pathBinaryHasCorrectFlags() {
		ffmpegPath, _ = exec.LookPath("ffmpeg")
		ffprobePath, _ = exec.LookPath("ffprobe")
	}

	// Check if ffmpeg exists in the config directory
	ffmpegConfigPath := filepath.Join(configDirectory, getFFMPEGFilename())
	ffprobeConfigPath := filepath.Join(configDirectory, getFFProbeFilename())
	ffmpegConfigExists, _ := utils.FileExists(ffmpegConfigPath)
	ffprobeConfigExists, _ := utils.FileExists(ffprobeConfigPath)
	if ffmpegPath == "" && ffmpegConfigExists {
		ffmpegPath = ffmpegConfigPath
	}
	if ffprobePath == "" && ffprobeConfigExists {
		ffprobePath = ffprobeConfigPath
	}

	return ffmpegPath, ffprobePath
}

func Download(configDirectory string) error {
	url := getFFMPEGURL()
	if url == "" {
		return fmt.Errorf("no ffmpeg url for this platform")
	}

	// Configure where we want to download the archive
	urlExt := path.Ext(url)
	archivePath := filepath.Join(configDirectory, "ffmpeg"+urlExt)
	_ = os.Remove(archivePath) // remove archive if it already exists
	out, err := os.Create(archivePath)
	if err != nil  {
		return err
	}
	defer out.Close()

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the response to the archive file location
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return err
	}

	if urlExt == ".zip" {
		if err := unzip(archivePath, configDirectory); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("FFMPeg was downloaded to %s.  ")
	}

	return nil
}

func getFFMPEGURL() string {
	switch runtime.GOOS {
	case "darwin":
		return "https://ffmpeg.zeranoe.com/builds/macos64/static/ffmpeg-4.1-macos64-static.zip"
	case "linux":
		// TODO: untar this
		//return "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"
		return ""
	case "windows":
		return "https://ffmpeg.zeranoe.com/builds/win64/static/ffmpeg-4.1-win64-static.zip"
	default:
		return ""
	}
}

func getFFMPEGFilename() string {
	if runtime.GOOS == "windows" {
		return "ffmpeg.exe"
	} else {
		return "ffmpeg"
	}
}

func getFFProbeFilename() string {
	if runtime.GOOS == "windows" {
		return "ffprobe.exe"
	} else {
		return "ffprobe"
	}
}

// Checks if FFMPEG in the path has the correct flags
func pathBinaryHasCorrectFlags() bool {
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return false
	}
	bytes, _ := exec.Command(ffmpegPath).CombinedOutput()
	output := string(bytes)
	hasOpus := strings.Contains(output, "--enable-libopus")
	hasVpx := strings.Contains(output, "--enable-libvpx")
	hasX264 := strings.Contains(output, "--enable-libx264")
	hasX265 := strings.Contains(output, "--enable-libx265")
	hasWebp := strings.Contains(output, "--enable-libwebp")
	return hasOpus && hasVpx && hasX264 && hasX265 && hasWebp
}

func unzip(src, configDirectory string) error {
	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		filename := f.FileInfo().Name()
		if filename != "ffprobe" && filename != "ffmpeg" && filename != "ffprobe.exe" && filename != "ffmpeg.exe" {
			continue
		}

		rc, err := f.Open()

		unzippedPath := filepath.Join(configDirectory, filename)
		unzippedOutput, err := os.Create(unzippedPath)
		if err != nil  {
			return err
		}

		_, err = io.Copy(unzippedOutput, rc)
		if err != nil {
			return err
		}

		if err := unzippedOutput.Close(); err != nil {
			return err
		}
	}

	return nil
}
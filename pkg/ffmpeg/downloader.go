package ffmpeg

import (
	"archive/zip"
	"fmt"
	"github.com/stashapp/stash/pkg/utils"
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
	for _, url := range getFFMPEGURL() {
		err := DownloadSingle(configDirectory, url)
		if err != nil {
			return err
		}
	}
	return nil
}

func DownloadSingle(configDirectory, url string) error {
	if url == "" {
		return fmt.Errorf("no ffmpeg url for this platform")
	}

	// Configure where we want to download the archive
	urlExt := path.Ext(url)
	urlBase := path.Base(url)
	archivePath := filepath.Join(configDirectory, urlBase)
	_ = os.Remove(archivePath) // remove archive if it already exists
	out, err := os.Create(archivePath)
	if err != nil {
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
	if err != nil {
		return err
	}

	if urlExt == ".zip" {
		if err := unzip(archivePath, configDirectory); err != nil {
			return err
		}

		// On OSX or Linux set downloaded files permissions
		if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
			if err := os.Chmod(filepath.Join(configDirectory, "ffmpeg"), 0755); err != nil {
				return err
			}

			if err := os.Chmod(filepath.Join(configDirectory, "ffprobe"), 0755); err != nil {
				return err
			}

			// TODO: In future possible clear xattr to allow running on osx without user intervention
			// TODO: this however may not be required.
			// xattr -c /path/to/binary -- xattr.Remove(path, "com.apple.quarantine")
		}

	} else {
		return fmt.Errorf("ffmpeg was downloaded to %s", archivePath)
	}

	return nil
}

func getFFMPEGURL() []string {
	urls := []string{""}
	switch runtime.GOOS {
	case "darwin":
		urls = []string{"https://evermeet.cx/ffmpeg/ffmpeg-4.3.1.zip", "https://evermeet.cx/ffmpeg/ffprobe-4.3.1.zip"}
	case "linux":
		// TODO: get appropriate arch (arm,arm64,amd64) and xz untar from https://johnvansickle.com/ffmpeg/
		//       or get the ffmpeg,ffprobe zip repackaged ones from  https://ffbinaries.com/downloads
		urls = []string{""}
	case "windows":
		urls = []string{"https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip"}
	default:
		urls = []string{""}
	}
	return urls
}

func getFFMPEGFilename() string {
	if runtime.GOOS == "windows" {
		return "ffmpeg.exe"
	}
	return "ffmpeg"
}

func getFFProbeFilename() string {
	if runtime.GOOS == "windows" {
		return "ffprobe.exe"
	}
	return "ffprobe"
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
		if err != nil {
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

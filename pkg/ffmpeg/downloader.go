package ffmpeg

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	stashExec "github.com/stashapp/stash/pkg/exec"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

func GetPaths(paths []string) (string, string) {
	var ffmpegPath, ffprobePath string

	// Check if ffmpeg exists in the PATH
	if pathBinaryHasCorrectFlags() {
		ffmpegPath, _ = exec.LookPath("ffmpeg")
		ffprobePath, _ = exec.LookPath("ffprobe")
	}

	// Check if ffmpeg exists in the config directory
	if ffmpegPath == "" {
		ffmpegPath = fsutil.FindInPaths(paths, getFFMpegFilename())
	}
	if ffprobePath == "" {
		ffprobePath = fsutil.FindInPaths(paths, getFFProbeFilename())
	}

	return ffmpegPath, ffprobePath
}

func Download(ctx context.Context, configDirectory string) error {
	for _, url := range getFFmpegURL() {
		err := downloadSingle(ctx, configDirectory, url)
		if err != nil {
			return err
		}
	}

	// validate that the urls contained what we needed
	executables := []string{getFFMpegFilename(), getFFProbeFilename()}
	for _, executable := range executables {
		_, err := os.Stat(filepath.Join(configDirectory, executable))
		if err != nil {
			return err
		}
	}
	return nil
}

type progressReader struct {
	io.Reader
	lastProgress int64
	bytesRead    int64
	total        int64
}

func (r *progressReader) Read(p []byte) (int, error) {
	read, err := r.Reader.Read(p)
	if err == nil {
		r.bytesRead += int64(read)
		if r.total > 0 {
			progress := int64(float64(r.bytesRead) / float64(r.total) * 100)
			if progress/5 > r.lastProgress {
				logger.Infof("%d%% downloaded...", progress)
				r.lastProgress = progress / 5
			}
		}
	}

	return read, err
}

func downloadSingle(ctx context.Context, configDirectory, url string) error {
	if url == "" {
		return fmt.Errorf("no ffmpeg url for this platform")
	}

	// Configure where we want to download the archive
	urlBase := path.Base(url)
	archivePath := filepath.Join(configDirectory, urlBase)
	_ = os.Remove(archivePath) // remove archive if it already exists
	out, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer out.Close()

	logger.Infof("Downloading %s...", url)

	// Make the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	reader := &progressReader{
		Reader: resp.Body,
		total:  resp.ContentLength,
	}

	// Write the response to the archive file location
	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}
	logger.Info("Downloading complete")

	mime := resp.Header.Get("Content-Type")
	if mime != "application/zip" { // try detecting MIME type since some servers don't return the correct one
		data := make([]byte, 500) // http.DetectContentType only reads up to 500 bytes
		_, _ = out.ReadAt(data, 0)
		mime = http.DetectContentType(data)
	}

	if mime == "application/zip" {
		logger.Infof("Unzipping %s...", archivePath)
		if err := unzip(archivePath, configDirectory); err != nil {
			return err
		}

		// On OSX or Linux set downloaded files permissions
		if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
			_, err = os.Stat(filepath.Join(configDirectory, "ffmpeg"))
			if !os.IsNotExist(err) {
				if err = os.Chmod(filepath.Join(configDirectory, "ffmpeg"), 0755); err != nil {
					return err
				}
			}

			_, err = os.Stat(filepath.Join(configDirectory, "ffprobe"))
			if !os.IsNotExist(err) {
				if err := os.Chmod(filepath.Join(configDirectory, "ffprobe"), 0755); err != nil {
					return err
				}
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

func getFFmpegURL() []string {
	var urls []string
	switch runtime.GOOS {
	case "darwin":
		urls = []string{"https://evermeet.cx/ffmpeg/getrelease/zip", "https://evermeet.cx/ffmpeg/getrelease/ffprobe/zip"}
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			urls = []string{"https://github.com/ffbinaries/ffbinaries-prebuilt/releases/download/v4.2.1/ffmpeg-4.2.1-linux-64.zip", "https://github.com/ffbinaries/ffbinaries-prebuilt/releases/download/v4.2.1/ffprobe-4.2.1-linux-64.zip"}
		case "arm":
			urls = []string{"https://github.com/ffbinaries/ffbinaries-prebuilt/releases/download/v4.2.1/ffmpeg-4.2.1-linux-armhf-32.zip", "https://github.com/ffbinaries/ffbinaries-prebuilt/releases/download/v4.2.1/ffprobe-4.2.1-linux-armhf-32.zip"}
		case "arm64":
			urls = []string{"https://github.com/ffbinaries/ffbinaries-prebuilt/releases/download/v4.2.1/ffmpeg-4.2.1-linux-arm-64.zip", "https://github.com/ffbinaries/ffbinaries-prebuilt/releases/download/v4.2.1/ffprobe-4.2.1-linux-arm-64.zip"}
		}
	case "windows":
		urls = []string{"https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip"}
	default:
		urls = []string{""}
	}
	return urls
}

func getFFMpegFilename() string {
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

// Checks if ffmpeg in the path has the correct flags
func pathBinaryHasCorrectFlags() bool {
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return false
	}
	cmd := stashExec.Command(ffmpegPath)
	bytes, _ := cmd.CombinedOutput()
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
		if err != nil {
			return err
		}

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

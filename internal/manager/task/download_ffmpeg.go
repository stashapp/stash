package task

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
)

type DownloadFFmpegJob struct {
	ConfigDirectory string
	OnComplete      func(ctx context.Context)
	urls            []string
	downloaded      int
}

func (s *DownloadFFmpegJob) Execute(ctx context.Context, progress *job.Progress) error {
	if err := s.download(ctx, progress); err != nil {
		if job.IsCancelled(ctx) {
			return nil
		}
		return err
	}

	if s.OnComplete != nil {
		s.OnComplete(ctx)
	}

	return nil
}

func (s *DownloadFFmpegJob) setTaskProgress(taskProgress float64, progress *job.Progress) {
	progress.SetPercent((float64(s.downloaded) + taskProgress) / float64(len(s.urls)))
}

func (s *DownloadFFmpegJob) download(ctx context.Context, progress *job.Progress) error {
	s.urls = ffmpeg.GetFFmpegURL()

	// set steps based on the number of URLs

	for _, url := range s.urls {
		err := s.downloadSingle(ctx, url, progress)
		if err != nil {
			return err
		}
		s.downloaded++
	}

	// validate that the urls contained what we needed
	executables := []string{fsutil.GetExeName("ffmpeg"), fsutil.GetExeName("ffprobe")}
	for _, executable := range executables {
		_, err := os.Stat(filepath.Join(s.ConfigDirectory, executable))
		if err != nil {
			return err
		}
	}
	return nil
}

type downloadProgressReader struct {
	io.Reader
	setProgress func(taskProgress float64)
	bytesRead   int64
	total       int64
}

func (r *downloadProgressReader) Read(p []byte) (int, error) {
	read, err := r.Reader.Read(p)
	if err == nil {
		r.bytesRead += int64(read)
		if r.total > 0 {
			progress := float64(r.bytesRead) / float64(r.total)
			r.setProgress(progress)
		}
	}

	return read, err
}

func (s *DownloadFFmpegJob) downloadSingle(ctx context.Context, url string, progress *job.Progress) error {
	if url == "" {
		return fmt.Errorf("no ffmpeg url for this platform")
	}

	configDirectory := s.ConfigDirectory

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

	progress.ExecuteTask(fmt.Sprintf("Downloading %s", url), func() {
		err = s.downloadFile(ctx, url, out, progress)
	})

	if err != nil {
		return fmt.Errorf("failed to download ffmpeg from %s: %w", url, err)
	}

	logger.Info("Downloading complete")

	logger.Infof("Unzipping %s...", archivePath)
	progress.ExecuteTask(fmt.Sprintf("Unzipping %s", archivePath), func() {
		err = s.unzip(archivePath)
	})

	if err != nil {
		return fmt.Errorf("failed to unzip ffmpeg archive: %w", err)
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

	return nil
}

func (s *DownloadFFmpegJob) downloadFile(ctx context.Context, url string, out *os.File, progress *job.Progress) error {
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

	reader := &downloadProgressReader{
		Reader: resp.Body,
		total:  resp.ContentLength,
		setProgress: func(taskProgress float64) {
			s.setTaskProgress(taskProgress, progress)
		},
	}

	// Write the response to the archive file location
	if _, err := io.Copy(out, reader); err != nil {
		return err
	}

	mime := resp.Header.Get("Content-Type")
	if mime != "application/zip" { // try detecting MIME type since some servers don't return the correct one
		data := make([]byte, 500) // http.DetectContentType only reads up to 500 bytes
		_, _ = out.ReadAt(data, 0)
		mime = http.DetectContentType(data)
	}

	if mime != "application/zip" {
		return fmt.Errorf("downloaded file is not a zip archive")
	}

	return nil
}

func (s *DownloadFFmpegJob) unzip(src string) error {
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

		unzippedPath := filepath.Join(s.ConfigDirectory, filename)
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

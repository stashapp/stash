package ffmpeg

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

const (
	MimeWebmVideo string = "video/webm"
	MimeWebmAudio string = "audio/webm"
	MimeMkvVideo  string = "video/x-matroska"
	MimeMkvAudio  string = "audio/x-matroska"
	MimeMp4Video  string = "video/mp4"
	MimeMp4Audio  string = "audio/mp4"
)

type StreamManager struct {
	cacheDir string
	encoder  *FFMpeg
	ffprobe  *FFProbe

	config      StreamManagerConfig
	lockManager *fsutil.ReadLockManager

	context    context.Context
	cancelFunc context.CancelFunc

	runningStreams map[string]*runningStream
	streamsMutex   sync.Mutex
}

type StreamManagerConfig interface {
	GetMaxStreamingTranscodeSize() models.StreamingResolutionEnum
	GetLiveTranscodeInputArgs() []string
	GetLiveTranscodeOutputArgs() []string
	GetTranscodeHardwareAcceleration() bool
}

func NewStreamManager(cacheDir string, encoder *FFMpeg, ffprobe *FFProbe, config StreamManagerConfig, lockManager *fsutil.ReadLockManager) *StreamManager {
	if cacheDir == "" {
		logger.Warn("cache directory is not set. Live HLS/DASH transcoding will be disabled")
	}

	ctx, cancel := context.WithCancel(context.Background())

	ret := &StreamManager{
		cacheDir:       cacheDir,
		encoder:        encoder,
		ffprobe:        ffprobe,
		config:         config,
		lockManager:    lockManager,
		context:        ctx,
		cancelFunc:     cancel,
		runningStreams: make(map[string]*runningStream),
	}

	go func() {
		for {
			select {
			case <-time.After(monitorInterval):
				ret.monitorStreams()
			case <-ctx.Done():
				return
			}
		}
	}()

	return ret
}

// Shutdown shuts down the stream manager, killing any running transcoding processes and removing all cached files.
func (sm *StreamManager) Shutdown() {
	sm.cancelFunc()
	sm.stopAndRemoveAll()
}

type StreamRequestContext struct {
	context.Context
	ResponseWriter http.ResponseWriter
}

func NewStreamRequestContext(w http.ResponseWriter, r *http.Request) *StreamRequestContext {
	return &StreamRequestContext{
		Context:        r.Context(),
		ResponseWriter: w,
	}
}

func (c *StreamRequestContext) Cancel() {
	hj, ok := (c.ResponseWriter).(http.Hijacker)
	if !ok {
		return
	}

	// hijack and close the connection
	conn, bw, _ := hj.Hijack()
	if conn != nil {
		if bw != nil {
			// notify end of stream
			_, err := bw.WriteString("0\r\n")
			if err != nil {
				logger.Warnf("unable to write end of stream: %v", err)
			}
			_, err = bw.WriteString("\r\n")
			if err != nil {
				logger.Warnf("unable to write end of stream: %v", err)
			}

			// flush the buffer, but don't wait indefinitely
			timeout := make(chan struct{}, 1)
			go func() {
				_ = bw.Flush()
				close(timeout)
			}()

			const waitTime = time.Second

			select {
			case <-timeout:
			case <-time.After(waitTime):
				logger.Warnf("unable to flush buffer - closing connection")
			}
		}

		conn.Close()
	}
}

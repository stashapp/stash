package models

import (
	"context"
	"time"
)

// PlaybackStreamType identifies which stream endpoint format to use.
type PlaybackStreamType string

const (
	PlaybackStreamTypeDirect PlaybackStreamType = "DIRECT"
	PlaybackStreamTypeMP4    PlaybackStreamType = "MP4"
	PlaybackStreamTypeWEBM   PlaybackStreamType = "WEBM"
	PlaybackStreamTypeMKV    PlaybackStreamType = "MKV"
	PlaybackStreamTypeHLS    PlaybackStreamType = "HLS"
	PlaybackStreamTypeDASH   PlaybackStreamType = "DASH"
)

func (t PlaybackStreamType) IsValid() bool {
	switch t {
	case PlaybackStreamTypeDirect, PlaybackStreamTypeMP4, PlaybackStreamTypeWEBM,
		PlaybackStreamTypeMKV, PlaybackStreamTypeHLS, PlaybackStreamTypeDASH:
		return true
	}
	return false
}

func (t PlaybackStreamType) String() string {
	return string(t)
}

// PlaybackDefault maps a user-agent substring pattern to a preferred stream
// type and quality. Lower Priority value wins when multiple patterns match.
type PlaybackDefault struct {
	ID                 int
	UserAgentPattern   string
	Priority           int
	StreamType         PlaybackStreamType
	Quality            *StreamingResolutionEnum
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type PlaybackDefaultCreateInput struct {
	UserAgentPattern string
	Priority         int
	StreamType       string
	Quality          *StreamingResolutionEnum
}

type PlaybackDefaultUpdateInput struct {
	ID               string
	UserAgentPattern *string
	Priority         *int
	StreamType       *string
	Quality          *StreamingResolutionEnum
}

type PlaybackDefaultReader interface {
	GetAll(ctx context.Context) ([]*PlaybackDefault, error)
	// FindByUserAgent returns the highest-priority rule whose UserAgentPattern
	// is a case-insensitive substring of ua, or nil if none match.
	FindByUserAgent(ctx context.Context, ua string) (*PlaybackDefault, error)
	Find(ctx context.Context, id int) (*PlaybackDefault, error)
}

type PlaybackDefaultWriter interface {
	Create(ctx context.Context, pd *PlaybackDefault) error
	Update(ctx context.Context, pd *PlaybackDefault) error
	Destroy(ctx context.Context, id int) error
	// Upsert inserts or replaces the stream_type and quality for the given
	// user_agent_pattern, leaving priority unchanged on conflict.
	Upsert(ctx context.Context, pd *PlaybackDefault) error
}

type PlaybackDefaultReaderWriter interface {
	PlaybackDefaultReader
	PlaybackDefaultWriter
}

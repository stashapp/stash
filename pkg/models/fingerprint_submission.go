package models

import (
	"context"
	"time"
)

type FingerprintVote string

const (
	FingerprintVoteValid   FingerprintVote = "VALID"
	FingerprintVoteInvalid FingerprintVote = "INVALID"
)

func (e FingerprintVote) IsValid() bool {
	switch e {
	case FingerprintVoteValid, FingerprintVoteInvalid:
		return true
	}
	return false
}

func (e FingerprintVote) String() string {
	return string(e)
}

type FingerprintSubmission struct {
	Endpoint  string          `json:"endpoint"`
	StashID   string          `json:"stash_id"`
	SceneID   int             `json:"scene_id"`
	Vote      FingerprintVote `json:"vote"`
	CreatedAt time.Time       `json:"created_at"`
}

type FingerprintSubmissionReader interface {
	FindByEndpoint(ctx context.Context, endpoint string) ([]*FingerprintSubmission, error)
	Find(ctx context.Context, endpoint string, stashID string) (*FingerprintSubmission, error)
}

type FingerprintSubmissionWriter interface {
	Create(ctx context.Context, newObject *FingerprintSubmission) error
	Delete(ctx context.Context, endpoint string, stashID string) error
	DeleteByEndpoint(ctx context.Context, endpoint string) error
}

type FingerprintSubmissionReaderWriter interface {
	FingerprintSubmissionReader
	FingerprintSubmissionWriter
}

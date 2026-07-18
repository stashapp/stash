package blob

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/stashapp/stash/pkg/logger"
)

// S3Options configures the S3-compatible blob store.
// Works with any S3-compatible service (AWS S3, RustFS, Garage, SeaweedFS, ...).
type S3Options struct {
	// Endpoint is the S3 host, e.g. "s3.amazonaws.com" or "rustfs:9000"
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	// Prefix is prepended to every object key, e.g. "blobs/"
	Prefix string
	Region string
	UseSSL bool
}

// Configured returns true if the options carry enough information to build a client.
func (o S3Options) Configured() bool {
	return o.Endpoint != "" && o.Bucket != ""
}

// S3Store stores blobs as objects keyed by prefix + checksum.
// The client is created lazily so that construction never fails.
type S3Store struct {
	options S3Options

	initOnce sync.Once
	client   *s3.Client
	initErr  error
}

func NewS3Store(options S3Options) *S3Store {
	return &S3Store{options: options}
}

func (s *S3Store) getClient() (*s3.Client, error) {
	s.initOnce.Do(func() {
		if !s.options.Configured() {
			s.initErr = fmt.Errorf("s3 blob store is not configured (endpoint and bucket are required)")
			return
		}

		scheme := "http"
		if s.options.UseSSL {
			scheme = "https"
		}

		region := s.options.Region
		if region == "" {
			region = "us-east-1"
		}

		s.client = s3.New(s3.Options{
			BaseEndpoint: aws.String(scheme + "://" + s.options.Endpoint),
			Region:       region,
			Credentials:  credentials.NewStaticCredentialsProvider(s.options.AccessKeyID, s.options.SecretAccessKey, ""),
			// path-style addressing for self-hosted S3-compatible services
			UsePathStyle: true,
		})
	})

	return s.client, s.initErr
}

func (s *S3Store) checksumToKey(checksum string) string {
	return s.options.Prefix + checksum
}

// isNotFound returns true if err indicates a missing object.
func isNotFound(err error) bool {
	var noSuchKey *types.NoSuchKey
	if errors.As(err, &noSuchKey) {
		return true
	}

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "NoSuchKey", "NotFound":
			return true
		}
	}

	return false
}

// Read returns the blob for checksum. A missing object is reported as an
// error wrapping fs.ErrNotExist, matching the filesystem store semantics.
func (s *S3Store) Read(ctx context.Context, checksum string) ([]byte, error) {
	client, err := s.getClient()
	if err != nil {
		return nil, err
	}

	key := s.checksumToKey(checksum)
	out, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.options.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("object %q: %w", key, fs.ErrNotExist)
		}
		return nil, fmt.Errorf("getting object %q: %w", key, err)
	}
	defer out.Body.Close()

	data, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, fmt.Errorf("reading object %q: %w", key, err)
	}

	return data, nil
}

func (s *S3Store) Write(ctx context.Context, checksum string, data []byte) error {
	client, err := s.getClient()
	if err != nil {
		return err
	}

	key := s.checksumToKey(checksum)
	logger.Debugf("Writing blob object %s", key)

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.options.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/octet-stream"),
	})
	if err != nil {
		return fmt.Errorf("putting object %q: %w", key, err)
	}

	return nil
}

// Delete removes the blob object. Deleting a missing object is not an error.
func (s *S3Store) Delete(ctx context.Context, checksum string) error {
	client, err := s.getClient()
	if err != nil {
		return err
	}

	key := s.checksumToKey(checksum)
	_, err = client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.options.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if isNotFound(err) {
			return nil
		}
		return fmt.Errorf("removing object %q: %w", key, err)
	}

	return nil
}

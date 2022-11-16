package utils

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"syscall"
	"time"
)

// Timeout to get the image. Includes transfer time. May want to make this
// configurable at some point.
const imageGetTimeout = time.Second * 60

const base64RE = `^data:.+\/(.+);base64,(.*)$`

// ProcessImageInput transforms an image string either from a base64 encoded
// string, or from a URL, and returns the image as a byte slice
func ProcessImageInput(ctx context.Context, imageInput string) ([]byte, error) {
	regex := regexp.MustCompile(base64RE)
	if regex.MatchString(imageInput) {
		d, err := ProcessBase64Image(imageInput)
		return d, err
	}

	// assume input is a URL. Read it.
	return ReadImageFromURL(ctx, imageInput)
}

// ReadImageFromURL returns image data from a URL
func ReadImageFromURL(ctx context.Context, url string) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{ // ignore insecure certificates
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},

		Timeout: imageGetTimeout,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// assume is a URL for now

	// set the host of the URL as the referer
	if req.URL.Scheme != "" {
		req.Header.Set("Referer", req.URL.Scheme+"://"+req.Host+"/")
	}
	req.Header.Set("User-Agent", getUserAgent())

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http error %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// ProcessBase64Image transforms a base64 encoded string from a form post and
// returns the image itself as a byte slice.
func ProcessBase64Image(imageString string) ([]byte, error) {
	if imageString == "" {
		return nil, fmt.Errorf("empty image string")
	}

	regex := regexp.MustCompile(base64RE)
	matches := regex.FindStringSubmatch(imageString)
	var encodedString string
	if len(matches) > 2 {
		encodedString = regex.FindStringSubmatch(imageString)[2]
	} else {
		encodedString = imageString
	}
	imageData, err := GetDataFromBase64String(encodedString)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

// GetDataFromBase64String returns the given base64 encoded string as a byte slice
func GetDataFromBase64String(encodedString string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedString)
}

// GetBase64StringFromData returns the given byte slice as a base64 encoded string
func GetBase64StringFromData(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func ServeImage(image []byte, w http.ResponseWriter, r *http.Request) error {
	etag := fmt.Sprintf("%x", md5.Sum(image))

	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			w.WriteHeader(http.StatusNotModified)
			return nil
		}
	}

	contentType := http.DetectContentType(image)
	if contentType == "text/xml; charset=utf-8" || contentType == "text/plain; charset=utf-8" {
		contentType = "image/svg+xml"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Add("Etag", etag)
	w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
	_, err := w.Write(image)
	// Broken pipe errors are common when serving images and the remote
	// connection closes the connection. Filter them out of the error
	// messages, as they are benign.
	if errors.Is(err, syscall.EPIPE) {
		return nil
	}
	return err
}

// Package brotli provides encoding for statigz.Server.
package brotli

import (
	"bytes"
	"io"

	"github.com/andybalholm/brotli"
	"github.com/vearutop/statigz"
)

// AddEncoding is an option that prepends brotli to encodings of statigz.Server.
//
// It is located in a separate package to allow better control of imports graph.
func AddEncoding(server *statigz.Server) {
	enc := statigz.Encoding{
		FileExt:         ".br",
		ContentEncoding: "br",
		Decoder: func(r io.Reader) (io.Reader, error) {
			return brotli.NewReader(r), nil
		},
		Encoder: func(r io.Reader) ([]byte, error) {
			res := bytes.NewBuffer(nil)
			w := brotli.NewWriterLevel(res, 8)

			if _, err := io.Copy(w, r); err != nil {
				return nil, err
			}

			if err := w.Close(); err != nil {
				return nil, err
			}

			return res.Bytes(), nil
		},
	}

	server.Encodings = append([]statigz.Encoding{enc}, server.Encodings...)
}

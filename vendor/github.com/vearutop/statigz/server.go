// Package statigz serves pre-compressed embedded files with http.
package statigz

import (
	"bytes"
	"compress/gzip"
	"hash/fnv"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Server is a http.Handler that directly serves compressed files from file system to capable agents.
//
// Please use FileServer to create an instance of Server.
//
// If agent does not accept encoding and uncompressed file is not available in file system,
// it would decompress the file before serving.
//
// Compressed files should have an additional extension to indicate their encoding,
// for example "style.css.gz" or "bundle.js.br".
//
// Caching is implemented with ETag and If-None-Match headers. Range requests are supported
// with help of http.ServeContent.
//
// Behavior is similar to http://nginx.org/en/docs/http/ngx_http_gzip_static_module.html and
// https://github.com/lpar/gzipped, except compressed data can be decompressed for an incapable agent.
type Server struct {
	// OnError controls error handling during Serve.
	OnError func(rw http.ResponseWriter, r *http.Request, err error)

	// Encodings contains supported encodings, default GzipEncoding.
	Encodings []Encoding

	// EncodeOnInit encodes files that does not have encoded version on Server init.
	// This allows embedding uncompressed files and still leverage one time compression
	// for multiple requests.
	// Enabling this option can degrade startup performance and memory usage in case
	// of large embeddings, use with caution.
	EncodeOnInit bool

	info map[string]fileInfo
	fs   fs.ReadDirFS
}

const (
	// minSizeToEncode is minimal file size to apply encoding in runtime, 0.5KiB.
	minSizeToEncode = 512

	// minCompressionRatio is a minimal compression ratio to serve encoded data, 97%.
	minCompressionRatio = 0.97
)

// SkipCompressionExt lists file extensions of data that is already compressed.
var SkipCompressionExt = []string{".gz", ".br", ".gif", ".jpg", ".png", ".webp"}

// FileServer creates an instance of Server from file system.
//
// This function indexes provided file system to optimize further serving,
// so it is not recommended running it in the loop (for example for each request).
//
// Typically, file system would be an embed.FS.
//
//   //go:embed *.png *.br
//	 var FS embed.FS
//
// Brotli support is optionally available with brotli.AddEncoding.
func FileServer(fs fs.ReadDirFS, options ...func(server *Server)) *Server {
	s := Server{
		fs:   fs,
		info: make(map[string]fileInfo),
		OnError: func(rw http.ResponseWriter, r *http.Request, err error) {
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		},
		Encodings: []Encoding{GzipEncoding()},
	}

	for _, o := range options {
		o(&s)
	}

	// Reading from "." is not expected to fail.
	if err := s.hashDir("."); err != nil {
		panic(err)
	}

	if s.EncodeOnInit {
		err := s.encodeFiles()
		if err != nil {
			panic(err)
		}
	}

	return &s
}

func (s *Server) encodeFiles() error {
	for _, enc := range s.Encodings {
		if enc.Encoder == nil {
			continue
		}

		for fn, i := range s.info {
			isEncoded := false

			for _, ext := range SkipCompressionExt {
				if strings.HasSuffix(fn, ext) {
					isEncoded = true

					break
				}
			}

			if isEncoded {
				continue
			}

			if _, found := s.info[fn+enc.FileExt]; found {
				continue
			}

			// Skip encoding of small data.
			if i.size < minSizeToEncode {
				continue
			}

			f, err := s.fs.Open(fn)
			if err != nil {
				return err
			}

			b, err := enc.Encoder(f)
			if err != nil {
				return err
			}

			// Skip encoding for non-compressible data.
			if float64(len(b))/float64(i.size) > minCompressionRatio {
				continue
			}

			s.info[fn+enc.FileExt] = fileInfo{
				hash:    i.hash + enc.FileExt,
				size:    len(b),
				content: b[0:len(b):len(b)],
			}
		}
	}

	return nil
}

func (s *Server) hashDir(p string) error {
	files, err := s.fs.ReadDir(p)
	if err != nil {
		return err
	}

	for _, f := range files {
		fn := path.Join(p, f.Name())

		if f.IsDir() {
			s.info[path.Clean(fn)] = fileInfo{
				isDir: true,
			}

			if err = s.hashDir(fn); err != nil {
				return err
			}

			continue
		}

		h := fnv.New64()

		f, err := s.fs.Open(fn)
		if err != nil {
			return err
		}

		n, err := io.Copy(h, f)
		if err != nil {
			return err
		}

		s.info[path.Clean(fn)] = fileInfo{
			hash: strconv.FormatUint(h.Sum64(), 36),
			size: int(n),
		}
	}

	return nil
}

func (s *Server) reader(fn string, info fileInfo) (io.Reader, error) {
	if info.content != nil {
		return bytes.NewReader(info.content), nil
	}

	return s.fs.Open(fn)
}

func (s *Server) serve(rw http.ResponseWriter, req *http.Request, fn, suf, enc string, info fileInfo,
	decompress func(r io.Reader) (io.Reader, error),
) {
	if m := req.Header.Get("If-None-Match"); m == info.hash {
		rw.WriteHeader(http.StatusNotModified)

		return
	}

	ctype := mime.TypeByExtension(filepath.Ext(fn))
	if ctype == "" {
		ctype = "application/octet-stream" // Prevent unreliable Content-Type detection on compressed data.
	}

	// This is used to enforce application/javascript MIME on Windows (https://github.com/golang/go/issues/32350)
	if strings.HasSuffix(req.URL.Path, ".js") {
		ctype = "application/javascript"
	}

	rw.Header().Set("Content-Type", ctype)
	rw.Header().Set("Etag", info.hash)

	if enc != "" {
		rw.Header().Set("Content-Encoding", enc)
	}

	if info.size > 0 {
		rw.Header().Set("Content-Length", strconv.Itoa(info.size))
	}

	if req.Method == http.MethodHead {
		return
	}

	r, err := s.reader(fn+suf, info)
	if err != nil {
		s.OnError(rw, req, err)

		return
	}

	if decompress != nil {
		r, err = decompress(r)
		if err != nil {
			rw.Header().Del("Etag")
			s.OnError(rw, req, err)

			return
		}
	}

	if rs, ok := r.(io.ReadSeeker); ok {
		http.ServeContent(rw, req, fn, time.Time{}, rs)

		return
	}

	_, err = io.Copy(rw, r)
	if err != nil {
		s.OnError(rw, req, err)

		return
	}
}

func (s *Server) minEnc(accessEncoding string, fn string) (fileInfo, Encoding) {
	var (
		minEnc  Encoding
		minInfo = fileInfo{size: -1}
	)

	for _, enc := range s.Encodings {
		if !strings.Contains(accessEncoding, enc.ContentEncoding) {
			continue
		}

		info, found := s.info[fn+enc.FileExt]
		if !found {
			continue
		}

		if minInfo.size == -1 || info.size < minInfo.size {
			minEnc = enc
			minInfo = info
		}
	}

	return minInfo, minEnc
}

// ServeHTTP serves static files.
//
// For compatibility with std http.FileServer:
// if request path ends with /index.html, it is redirected to base directory;
// if request path points to a directory without trailing "/", it is redirected to a path with trailing "/".
func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		rw.Header().Set("Allow", http.MethodGet+", "+http.MethodHead)
		http.Error(rw, "Method Not Allowed\n\nmethod should be GET or HEAD", http.StatusMethodNotAllowed)

		return
	}

	if strings.HasSuffix(req.URL.Path, "/index.html") {
		localRedirect(rw, req, "./")

		return
	}

	fn := strings.TrimPrefix(req.URL.Path, "/")
	ae := req.Header.Get("Accept-Encoding")

	if s.info[fn].isDir {
		localRedirect(rw, req, path.Base(req.URL.Path)+"/")

		return
	}

	if fn == "" || strings.HasSuffix(fn, "/") {
		fn += "index.html"
	}

	// Always add Accept-Encoding to Vary to prevent intermediate caches corruption.
	rw.Header().Add("Vary", "Accept-Encoding")

	if ae != "" {
		minInfo, minEnc := s.minEnc(strings.ToLower(ae), fn)

		if minInfo.hash != "" {
			// Copy compressed data into response.
			s.serve(rw, req, fn, minEnc.FileExt, minEnc.ContentEncoding, minInfo, nil)

			return
		}
	}

	// Copy uncompressed data into response.
	uncompressedInfo, uncompressedFound := s.info[fn]
	if uncompressedFound {
		s.serve(rw, req, fn, "", "", uncompressedInfo, nil)

		return
	}

	// Decompress compressed data into response.
	for _, enc := range s.Encodings {
		info, found := s.info[fn+enc.FileExt]
		if !found || enc.Decoder == nil || info.isDir {
			continue
		}

		info.hash += "U"
		info.size = 0
		s.serve(rw, req, fn, enc.FileExt, "", info, enc.Decoder)

		return
	}

	http.NotFound(rw, req)
}

// Encoding describes content encoding.
type Encoding struct {
	// FileExt is an extension of file with compressed content, for example ".gz".
	FileExt string

	// ContentEncoding is encoding name that is used in Accept-Encoding and Content-Encoding
	// headers, for example "gzip".
	ContentEncoding string

	// Decoder is a function that can decode data for an agent that does not accept encoding,
	// can be nil to disable dynamic decompression.
	Decoder func(r io.Reader) (io.Reader, error)

	// Encoder is a function that can encode data
	Encoder func(r io.Reader) ([]byte, error)
}

type fileInfo struct {
	hash    string
	size    int
	content []byte
	isDir   bool
}

// OnError is an option to customize error handling in Server.
func OnError(onErr func(rw http.ResponseWriter, r *http.Request, err error)) func(server *Server) {
	return func(server *Server) {
		server.OnError = onErr
	}
}

// GzipEncoding provides gzip Encoding.
func GzipEncoding() Encoding {
	return Encoding{
		FileExt:         ".gz",
		ContentEncoding: "gzip",
		Decoder: func(r io.Reader) (io.Reader, error) {
			return gzip.NewReader(r)
		},
		Encoder: func(r io.Reader) ([]byte, error) {
			res := bytes.NewBuffer(nil)
			w := gzip.NewWriter(res)

			if _, err := io.Copy(w, r); err != nil {
				return nil, err
			}

			if err := w.Close(); err != nil {
				return nil, err
			}

			return res.Bytes(), nil
		},
	}
}

// EncodeOnInit enables runtime encoding for unencoded files to allow compression
// for uncompressed embedded files.
//
// Enabling this option can degrade startup performance and memory usage in case
// of large embeddings, use with caution.
func EncodeOnInit(server *Server) {
	server.EncodeOnInit = true
}

// localRedirect gives a Moved Permanently response.
// It does not convert relative paths to absolute paths like Redirect does.
//
// Copied go1.17/src/net/http/fs.go:685.
func localRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}

	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)
}

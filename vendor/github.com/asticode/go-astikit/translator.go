package astikit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Translator represents an object capable of translating stuff
type Translator struct {
	defaultLanguage string
	m               *sync.RWMutex // Lock p
	p               map[string]string
	validLanguages  map[string]bool
}

// TranslatorOptions represents Translator options
type TranslatorOptions struct {
	DefaultLanguage string
	ValidLanguages  []string
}

// NewTranslator creates a new Translator
func NewTranslator(o TranslatorOptions) (t *Translator) {
	t = &Translator{
		defaultLanguage: o.DefaultLanguage,
		m:               &sync.RWMutex{},
		p:               make(map[string]string),
		validLanguages:  make(map[string]bool),
	}
	for _, l := range o.ValidLanguages {
		t.validLanguages[l] = true
	}
	return
}

// ParseDir adds translations located in ".json" files in the specified dir
// If ".json" files are located in child dirs, keys will be prefixed with their paths
func (t *Translator) ParseDir(dirPath string) (err error) {
	// Default dir path
	if dirPath == "" {
		if dirPath, err = os.Getwd(); err != nil {
			err = fmt.Errorf("astikit: getwd failed: %w", err)
			return
		}
	}

	// Walk through dir
	if err = filepath.Walk(dirPath, func(path string, info os.FileInfo, e error) (err error) {
		// Check input error
		if e != nil {
			err = fmt.Errorf("astikit: walking %s has an input error for path %s: %w", dirPath, path, e)
			return
		}

		// Only process files
		if info.IsDir() {
			return
		}

		// Only process ".json" files
		if filepath.Ext(path) != ".json" {
			return
		}

		// Parse file
		if err = t.ParseFile(dirPath, path); err != nil {
			err = fmt.Errorf("astikit: parsing %s failed: %w", path, err)
			return
		}
		return
	}); err != nil {
		err = fmt.Errorf("astikit: walking %s failed: %w", dirPath, err)
		return
	}
	return
}

// ParseFile adds translation located in the provided path
func (t *Translator) ParseFile(dirPath, path string) (err error) {
	// Lock
	t.m.Lock()
	defer t.m.Unlock()

	// Open file
	var f *os.File
	if f, err = os.Open(path); err != nil {
		err = fmt.Errorf("astikit: opening %s failed: %w", path, err)
		return
	}
	defer f.Close()

	// Unmarshal
	var p map[string]interface{}
	if err = json.NewDecoder(f).Decode(&p); err != nil {
		err = fmt.Errorf("astikit: unmarshaling %s failed: %w", path, err)
		return
	}

	// Get language
	language := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	// Update valid languages
	t.validLanguages[language] = true

	// Get prefix
	prefix := language
	if dp := filepath.Dir(path); dp != dirPath {
		var fs []string
		for _, v := range strings.Split(strings.TrimPrefix(dp, dirPath), string(os.PathSeparator)) {
			if v != "" {
				fs = append(fs, v)
			}
		}
		prefix += "." + strings.Join(fs, ".")
	}

	// Parse
	t.parse(p, prefix)
	return
}

func (t *Translator) key(prefix, key string) string {
	return prefix + "." + key
}

func (t *Translator) parse(i map[string]interface{}, prefix string) {
	for k, v := range i {
		p := t.key(prefix, k)
		switch a := v.(type) {
		case string:
			t.p[p] = a
		case map[string]interface{}:
			t.parse(a, p)
		}
	}
}

// HTTPMiddleware is the Translator HTTP middleware
func (t *Translator) HTTPMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Store language in context
		if l := r.Header.Get("Accept-Language"); l != "" {
			*r = *r.WithContext(contextWithTranslatorLanguage(r.Context(), t.parseAcceptLanguage(l)))
		}

		// Next handler
		h.ServeHTTP(rw, r)
	})
}

func (t *Translator) parseAcceptLanguage(h string) string {
	// Split on comma
	var qs []float64
	ls := make(map[float64][]string)
	for _, c := range strings.Split(strings.TrimSpace(h), ",") {
		// Empty
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}

		// Split on semi colon
		ss := strings.Split(c, ";")

		// Parse coefficient
		q := float64(1)
		if len(ss) > 1 {
			s := strings.TrimSpace(ss[1])
			if strings.HasPrefix(s, "q=") {
				var err error
				if q, err = strconv.ParseFloat(strings.TrimPrefix(s, "q="), 64); err != nil {
					q = 1
				}
			}
		}

		// Add
		if _, ok := ls[q]; !ok {
			qs = append(qs, q)
		}
		ls[q] = append(ls[q], strings.TrimSpace(ss[0]))
	}

	// Order coefficients
	sort.Float64s(qs)

	// Loop through coefficients in reverse order
	for idx := len(qs) - 1; idx >= 0; idx-- {
		for _, l := range ls[qs[idx]] {
			if _, ok := t.validLanguages[l]; ok {
				return l
			}
		}
	}
	return ""
}

const contextKeyTranslatorLanguage = contextKey("astikit.translator.language")

type contextKey string

func contextWithTranslatorLanguage(ctx context.Context, language string) context.Context {
	return context.WithValue(ctx, contextKeyTranslatorLanguage, language)
}

func translatorLanguageFromContext(ctx context.Context) string {
	v, ok := ctx.Value(contextKeyTranslatorLanguage).(string)
	if !ok {
		return ""
	}
	return v
}

func (t *Translator) language(language string) string {
	if language == "" {
		return t.defaultLanguage
	}
	return language
}

// LanguageCtx returns the translator language from the context, or the default language if not in the context
func (t *Translator) LanguageCtx(ctx context.Context) string {
	return t.language(translatorLanguageFromContext(ctx))
}

// Translate translates a key into a specific language
func (t *Translator) Translate(language, key string) string {
	// Lock
	t.m.RLock()
	defer t.m.RUnlock()

	// Get translation
	k1 := t.key(t.language(language), key)
	v, ok := t.p[k1]
	if ok {
		return v
	}

	// Default translation
	k2 := t.key(t.defaultLanguage, key)
	if v, ok = t.p[k2]; ok {
		return v
	}
	return k1
}

// Translatef translates a key into a specific language with optional formatting args
func (t *Translator) Translatef(language, key string, args ...interface{}) string {
	return fmt.Sprintf(t.Translate(language, key), args...)
}

// TranslateCtx is an alias for TranslateC
func (t *Translator) TranslateCtx(ctx context.Context, key string) string {
	return t.TranslateC(ctx, key)
}

// TranslateC translates a key using the language specified in the context
func (t *Translator) TranslateC(ctx context.Context, key string) string {
	return t.Translate(translatorLanguageFromContext(ctx), key)
}

func (t *Translator) TranslateCf(ctx context.Context, key string, args ...interface{}) string {
	return t.Translatef(translatorLanguageFromContext(ctx), key, args...)
}

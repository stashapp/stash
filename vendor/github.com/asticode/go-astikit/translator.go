package astikit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Translator represents an object capable of translating stuff
type Translator struct {
	m *sync.RWMutex // Lock p
	o TranslatorOptions
	p map[string]string
}

// TranslatorOptions represents Translator options
type TranslatorOptions struct {
	DefaultLanguage string
}

// NewTranslator creates a new Translator
func NewTranslator(o TranslatorOptions) *Translator {
	return &Translator{
		m: &sync.RWMutex{},
		o: o,
		p: make(map[string]string),
	}
}

// ParseDir adds translations located in ".json" files in the specified dir
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

		// Only process first level files
		if info.IsDir() {
			if path != dirPath {
				err = filepath.SkipDir
			}
			return
		}

		// Only process ".json" files
		if filepath.Ext(path) != ".json" {
			return
		}

		// Parse file
		if err = t.ParseFile(path); err != nil {
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
func (t *Translator) ParseFile(path string) (err error) {
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

	// Parse
	t.parse(p, strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
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
			*r = *r.WithContext(contextWithTranslatorLanguage(r.Context(), l))
		}

		// Next handler
		h.ServeHTTP(rw, r)
	})
}

const contextKeyTranslatorLanguage = "astikit.translator.language"

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
		return t.o.DefaultLanguage
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
	k2 := t.key(t.o.DefaultLanguage, key)
	if v, ok = t.p[k2]; ok {
		return v
	}
	return k1
}

// TranslateCtx translates a key using the language specified in the context
func (t *Translator) TranslateCtx(ctx context.Context, key string) string {
	return t.Translate(translatorLanguageFromContext(ctx), key)
}

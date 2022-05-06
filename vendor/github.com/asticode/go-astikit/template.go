package astikit

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

// Templater represents an object capable of storing and parsing templates
type Templater struct {
	layouts   []string
	m         sync.Mutex
	templates map[string]*template.Template
}

// NewTemplater creates a new templater
func NewTemplater() *Templater {
	return &Templater{templates: make(map[string]*template.Template)}
}

// AddLayoutsFromDir walks through a dir and add files as layouts
func (t *Templater) AddLayoutsFromDir(dirPath, ext string) (err error) {
	// Get layouts
	if err = filepath.Walk(dirPath, func(path string, info os.FileInfo, e error) (err error) {
		// Check input error
		if e != nil {
			err = fmt.Errorf("astikit: walking layouts has an input error for path %s: %w", path, e)
			return
		}

		// Only process files
		if info.IsDir() {
			return
		}

		// Check extension
		if ext != "" && filepath.Ext(path) != ext {
			return
		}

		// Read layout
		var b []byte
		if b, err = ioutil.ReadFile(path); err != nil {
			err = fmt.Errorf("astikit: reading %s failed: %w", path, err)
			return
		}

		// Add layout
		t.AddLayout(string(b))
		return
	}); err != nil {
		err = fmt.Errorf("astikit: walking layouts in %s failed: %w", dirPath, err)
		return
	}
	return
}

// AddTemplatesFromDir walks through a dir and add files as templates
func (t *Templater) AddTemplatesFromDir(dirPath, ext string) (err error) {
	// Loop through templates
	if err = filepath.Walk(dirPath, func(path string, info os.FileInfo, e error) (err error) {
		// Check input error
		if e != nil {
			err = fmt.Errorf("astikit: walking templates has an input error for path %s: %w", path, e)
			return
		}

		// Only process files
		if info.IsDir() {
			return
		}

		// Check extension
		if ext != "" && filepath.Ext(path) != ext {
			return
		}

		// Read file
		var b []byte
		if b, err = ioutil.ReadFile(path); err != nil {
			err = fmt.Errorf("astikit: reading template content of %s failed: %w", path, err)
			return
		}

		// Add template
		// We use ToSlash to homogenize Windows path
		if err = t.AddTemplate(filepath.ToSlash(strings.TrimPrefix(path, dirPath)), string(b)); err != nil {
			err = fmt.Errorf("astikit: adding template failed: %w", err)
			return
		}
		return
	}); err != nil {
		err = fmt.Errorf("astikit: walking templates in %s failed: %w", dirPath, err)
		return
	}
	return
}

// AddLayout adds a new layout
func (t *Templater) AddLayout(c string) {
	t.layouts = append(t.layouts, c)
}

// AddTemplate adds a new template
func (t *Templater) AddTemplate(path, content string) (err error) {
	// Parse
	var tpl *template.Template
	if tpl, err = t.Parse(content); err != nil {
		err = fmt.Errorf("astikit: parsing template for path %s failed: %w", path, err)
		return
	}

	// Add template
	t.m.Lock()
	t.templates[path] = tpl
	t.m.Unlock()
	return
}

// DelTemplate deletes a template
func (t *Templater) DelTemplate(path string) {
	t.m.Lock()
	defer t.m.Unlock()
	delete(t.templates, path)
}

// Template retrieves a templates
func (t *Templater) Template(path string) (tpl *template.Template, ok bool) {
	t.m.Lock()
	defer t.m.Unlock()
	tpl, ok = t.templates[path]
	return
}

// Parse parses the content of a template
func (t *Templater) Parse(content string) (o *template.Template, err error) {
	// Parse content
	o = template.New("root")
	if o, err = o.Parse(content); err != nil {
		err = fmt.Errorf("astikit: parsing template content failed: %w", err)
		return
	}

	// Parse layouts
	for idx, l := range t.layouts {
		if o, err = o.Parse(l); err != nil {
			err = fmt.Errorf("astikit: parsing layout #%d failed: %w", idx+1, err)
			return
		}
	}
	return
}

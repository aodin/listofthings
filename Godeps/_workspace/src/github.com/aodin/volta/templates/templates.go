package templates

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Templates holds the parsed template files and any local attributes.
type Templates struct {
	parsed *template.Template
	locals Attrs
}

// Add adds a template to the current parsed templates
func (t *Templates) Add(src string) error {
	if _, err := t.parsed.Parse(src); err != nil {
		return err
	}
	return nil
}

// SetAttr sets a local templates variable. Overwritten variables return an
// error, which can be ignored as needed.
func (t *Templates) SetAttr(key string, value interface{}) (err error) {
	if _, exists := t.locals[key]; exists {
		err = fmt.Errorf("templates: overwriting key %s in attrs", key)
	}
	t.locals[key] = value
	return
}

// Execute will render the given template define name with the given attrs.
// Panics on error, because that's how we roll.
func (t *Templates) Execute(w io.Writer, n string, attrs ...Attrs) {
	data := Attrs{}
	data.Merge(t.locals)
	for _, attr := range attrs {
		data.Merge(attr)
	}
	if err := t.parsed.ExecuteTemplate(w, n, data); err != nil {
		panic(fmt.Sprintf("templates: could execute template %s: %s", n, err))
	}
}

// New creates a new Templates instance by compiling all files recursively in
// the given directory. Panics on error.
func New(path string, attrs ...Attrs) *Templates {
	return create(path, "", "", attrs...)
}

// New creates a new Templates instance by compiling all files recursively in
// the given directory using the given delimiters. Panics on error.
func NewWithDelims(path, openTag, closeTag string, attrs ...Attrs) *Templates {
	return create(path, openTag, closeTag, attrs...)
}

func create(path, openTag, closeTag string, attrs ...Attrs) *Templates {
	data := Attrs{}
	for _, attr := range attrs {
		data.Merge(attr)
	}
	t := &Templates{
		parsed: template.New("template"),
		locals: data,
	}

	// Set alternative delims if both open and close tags are set
	if openTag != "" && closeTag != "" {
		t.parsed = t.parsed.Delims(openTag, closeTag)
	}

	// Recursively walk the given directory and build templates, returning
	// immediately on error.
	err := filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(strings.ToLower(name), ".html") {
			if _, err := t.parsed.ParseFiles(name); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("templates: error while parsing templates: %s", err))
	}
	return t
}

// Empty Creates an empty Templates
func Empty() *Templates {
	return &Templates{
		parsed: template.New("template"),
	}
}

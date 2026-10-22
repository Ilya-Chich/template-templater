package templates

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"sort"
)

// Engine is the filesystem adapter for the renderer port.
type Engine struct {
	tmpl *template.Template
}

// Load parses page, partial, and error templates in dir.
func Load(dir string) (*Engine, error) {
	patterns := []string{
		filepath.Join(dir, "*.tmpl"),
		filepath.Join(dir, "partials", "*.tmpl"),
		filepath.Join(dir, "errors", "*.tmpl"),
	}
	var files []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		files = append(files, matches...)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no templates in %s", dir)
	}
	sort.Strings(files)

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	return &Engine{tmpl: tmpl}, nil
}

// Render executes a define name. It does not write an HTTP response.
func (e *Engine) Render(ctx context.Context, name string, data any) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := e.tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

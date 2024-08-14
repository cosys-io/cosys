package gen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// NewFileAction is an action that create a new file.
type NewFileAction struct {
	path       string
	tmplString string
	ctx        any
	opts       genOptions
}

// NewFile returns a new file action that creates a file at the given path using the given template and context.
func NewFile(path string, tmplString string, ctx any, options ...genOption) *NewFileAction {
	opts := genOptions{
		false,
		false,
		false,
	}

	for _, option := range options {
		option(&opts)
	}

	return &NewFileAction{
		path:       path,
		tmplString: tmplString,
		ctx:        ctx,
		opts:       opts,
	}
}

// Act creates a new file.
func (a NewFileAction) Act() error {
	exists, err := pathExists(a.path)
	if err != nil {
		return err
	}

	if exists {
		if a.opts.deleteIfExists {
			if err := os.RemoveAll(a.path); err != nil {
				return err
			}
		} else {
			if a.opts.skipIfExists {
				return nil
			}
			return fmt.Errorf("file already exists: %s", a.path)
		}
	}

	if !a.opts.genHeadOnly {
		if err := os.MkdirAll(filepath.Dir(a.path), os.ModePerm); err != nil {
			return err
		}
	}

	file, err := os.Create(a.path)
	if err != nil {
		return err
	}

	tmpl, err := template.New("tmpl").Parse(a.tmplString)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	if err = tmpl.Execute(&buffer, a.ctx); err != nil {
		return err
	}

	if _, err = file.WriteString(buffer.String()); err != nil {
		return err
	}

	return nil
}

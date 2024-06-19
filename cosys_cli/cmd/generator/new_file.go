package gen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type NewFileAction struct {
	path    string
	tmplStr string
	ctx     any
	opts    genOptions
}

func NewFile(path string, tmplStr string, ctx any, options ...genOption) *NewFileAction {
	opts := &genOptions{
		false,
		false,
		false,
	}

	for _, option := range options {
		option(opts)
	}

	return &NewFileAction{
		path:    path,
		tmplStr: tmplStr,
		ctx:     ctx,
		opts:    *opts,
	}
}

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

	tmpl, err := template.New("tmpl").Parse(a.tmplStr)
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

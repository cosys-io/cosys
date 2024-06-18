package gen

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"text/template"
)

type ModifyFileAction struct {
	path       string
	patternStr string
	tmplStr    string
	ctx        any
}

func ModifyFile(path, patternStr, tmplStr string, ctx any) *ModifyFileAction {
	return &ModifyFileAction{
		path:       path,
		patternStr: patternStr,
		tmplStr:    tmplStr,
		ctx:        ctx,
	}
}

func (a ModifyFileAction) Act() error {
	exists, err := pathExists(a.path)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("file does not exist: %s", a.path)
	}

	oldFile, err := os.ReadFile(a.path)
	if err != nil {
		return err
	}

	pattern, err := regexp.Compile(a.patternStr)
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

	newFile := pattern.ReplaceAll(oldFile, buffer.Bytes())

	err = os.WriteFile(a.path, newFile, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

package gen

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"text/template"
)

// ModifyFileAction is an action that modifies a file.
type ModifyFileAction struct {
	path          string
	patternString string
	tmplString    string
	ctx           any
}

// ModifyFile returns a modify file action that modifies the file at the given path, replacing all text matching the given regex pattern with the given template and context.
func ModifyFile(path, patternString, tmplString string, ctx any) *ModifyFileAction {
	return &ModifyFileAction{
		path:          path,
		patternString: patternString,
		tmplString:    tmplString,
		ctx:           ctx,
	}
}

// Act modifies a file.
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

	pattern, err := regexp.Compile(a.patternString)
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

	newFile := pattern.ReplaceAll(oldFile, buffer.Bytes())

	err = os.WriteFile(a.path, newFile, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

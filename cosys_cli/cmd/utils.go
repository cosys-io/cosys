package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"text/template"
)

func generateDir(path string, options ...genOption) error {
	opts := &genOptions{
		false,
		false,
		false,
	}

	for _, option := range options {
		option(opts)
	}

	exists, err := pathExists(path)
	if err != nil {
		return err
	}

	if exists {
		if opts.deleteIfExists {
			os.RemoveAll(path)
		} else {
			if opts.skipIfExists {
				return nil
			}
			return fmt.Errorf("file already exists: %s", path)
		}
	}

	if opts.genHeadOnly {
		return os.Mkdir(path, os.ModePerm)
	} else {
		return os.MkdirAll(path, os.ModePerm)
	}
}

func generateFile(path string, tmplStr string, ctx any, options ...genOption) error {
	opts := &genOptions{
		false,
		false,
		false,
	}

	for _, option := range options {
		option(opts)
	}

	exists, err := pathExists(path)
	if err != nil {
		return err
	}

	if exists {
		if opts.deleteIfExists {
			os.RemoveAll(path)
		} else {
			if opts.skipIfExists {
				return nil
			}
			return fmt.Errorf("file already exists: %s", path)
		}
	}

	if !opts.genHeadOnly {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	tmpl, err := template.New("tmpl").Parse(tmplStr)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	if err = tmpl.Execute(&buffer, ctx); err != nil {
		return err
	}

	if _, err = file.WriteString(buffer.String()); err != nil {
		return err
	}

	return nil
}

func modifyFile(path string, patternStr string, tmplStr string, ctx any) error {
	exists, err := pathExists(path)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("file does not exist: %s", path)
	}

	oldFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		return err
	}

	tmpl, err := template.New("tmpl").Parse(tmplStr)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	if err = tmpl.Execute(&buffer, ctx); err != nil {
		return err
	}

	newFile := pattern.ReplaceAll(oldFile, buffer.Bytes())

	err = os.WriteFile(path, newFile, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

type genOptions struct {
	deleteIfExists bool
	skipIfExists   bool
	genHeadOnly    bool
}

type genOption func(*genOptions)

func deleteIfExists(options *genOptions) {
	options.deleteIfExists = true
}

func skipIfExists(options *genOptions) {
	options.skipIfExists = true
}

func genHeadOnly(options *genOptions) {
	options.genHeadOnly = true
}

func getModfile() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		return line[7:], nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("module not found")
}

func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

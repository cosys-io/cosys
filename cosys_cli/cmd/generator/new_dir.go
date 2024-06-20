package gen

import (
	"fmt"
	"os"
)

type NewDirAction struct {
	path string
	opts genOptions
}

func NewDir(path string, options ...genOption) *NewDirAction {
	opts := &genOptions{
		false,
		false,
		false,
	}

	for _, option := range options {
		option(opts)
	}

	return &NewDirAction{
		path: path,
		opts: *opts,
	}
}

func (a NewDirAction) Act() error {
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

	if a.opts.genHeadOnly {
		return os.Mkdir(a.path, os.ModePerm)
	} else {
		return os.MkdirAll(a.path, os.ModePerm)
	}
}

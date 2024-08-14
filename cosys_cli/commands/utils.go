package commands

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

// runCommand runs a command in shell.
func runCommand(command string, options ...runOption) error {
	cfg := &runConfigs{
		dir:    "",
		quiet:  false,
		cancel: nil,
	}

	for _, option := range options {
		option(cfg)
	}

	ctx := context.Background()
	var cmd *exec.Cmd

	args := strings.Split(command, " ")
	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}

	if cfg.cancel == nil {
		if len(args) == 1 {
			cmd = exec.Command(args[0])
		} else {
			cmd = exec.Command(args[0], args[1:]...)
		}
	} else {
		if len(args) == 1 {
			cmd = exec.CommandContext(ctx, args[0])
		} else {
			cmd = exec.CommandContext(ctx, args[0], args[1:]...)
		}
		cmd.Cancel = cfg.cancel(ctx)
	}

	if cfg.dir != "" {
		cmd.Dir = cfg.dir
	}

	if !cfg.quiet {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// pathExists returns whether a path exists.
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

// runConfigs are configurations for running a shell command.
type runConfigs struct {
	dir    string
	quiet  bool
	cancel func(context.Context) func() error
}

// runOption is a configuration for running a shell command.
type runOption func(*runConfigs)

// dir is a configuration for running a shell command,
// specifying which directory to run the command from.
func dir(dir string) runOption {
	return func(cfg *runConfigs) {
		cfg.dir = dir
	}
}

// quiet is a configuration for running a shell command,
// specifying that the command should be run without outputting to stdout and stderr.
func quiet(cfg *runConfigs) {
	cfg.quiet = true
}

// cancel is a configuration for running a shell command,
// providing a context for cancelling the command.
func cancel(cancelFunc func(context.Context) func() error) runOption {
	return func(cfg *runConfigs) {
		cfg.cancel = cancelFunc
	}
}

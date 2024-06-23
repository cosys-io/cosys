package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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

func RunCommand(command string, options ...RunOption) error {
	cfg := &RunConfigs{
		Dir:    "",
		Quiet:  false,
		Cancel: nil,
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

	if cfg.Cancel == nil {
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
		cmd.Cancel = cfg.Cancel(ctx)
	}

	if cfg.Dir != "" {
		cmd.Dir = cfg.Dir
	}

	if !cfg.Quiet {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

type RunConfigs struct {
	Dir    string
	Quiet  bool
	Cancel func(context.Context) func() error
}

type RunOption func(*RunConfigs)

func Dir(dir string) RunOption {
	return func(cfg *RunConfigs) {
		cfg.Dir = dir
	}
}

func Quiet(cfg *RunConfigs) {
	cfg.Quiet = true
}

func Cancel(cancelFunc func(context.Context) func() error) RunOption {
	return func(cfg *RunConfigs) {
		cfg.Cancel = cancelFunc
	}
}

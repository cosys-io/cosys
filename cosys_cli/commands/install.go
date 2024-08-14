package commands

import (
	"fmt"
	gen "github.com/cosys-io/cosys/cosys_cli/generator"
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install [module_names]",
	Short: "Install modules",
	Long:  "Install modules.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := installModules(args); err != nil {
			log.Fatal(err)
		}
	},
}

func installModules(modules []string, options ...RunOption) error {
	cfg := &RunConfigs{
		Dir:    "",
		Quiet:  false,
		Cancel: nil,
	}

	for _, option := range options {
		option(cfg)
	}

	clonePath := filepath.Join(cfg.Dir, ".clone")
	if err := gen.NewDir(clonePath).Act(); err != nil {
		return err
	}
	defer os.RemoveAll(clonePath)

	errors := make(chan error, len(modules))
	for _, module := range modules {
		go func() {
			if err := installModule(module, options...); err != nil {
				errors <- err
			}

			errors <- nil
		}()
	}

	for _ = range len(modules) {
		if err := <-errors; err != nil {
			return err
		}
	}

	return nil
}

var repoRegexp = regexp.MustCompile(`^github\.com\/([a-zA-Z0-9_-]+)\/([a-zA-Z0-9_-]+)\/([a-zA-Z0-9_\-\/]*)$`)

func installModule(module string, options ...RunOption) error {
	cfg := &RunConfigs{
		Dir:    "",
		Quiet:  false,
		Cancel: nil,
	}

	for _, option := range options {
		option(cfg)
	}

	if module[len(module)-1] != '/' {
		module = module + "/"
	}

	matches := repoRegexp.FindStringSubmatch(module)
	if matches == nil {
		return fmt.Errorf("invalid module name: %s", module)
	}

	orgName := matches[1]
	repoName := matches[2]
	path := matches[3]
	moduleName := filepath.Base(path)

	if err := RunCommand(fmt.Sprintf("git clone --depth 1 --filter=blob:none --no-checkout git@github.com:%s/%s.git "+moduleName, orgName, repoName), append(options, Dir(filepath.Join(cfg.Dir, ".clone")))...); err != nil {
		return err
	}

	clonePath := filepath.Join(cfg.Dir, ".clone", moduleName)

	if err := RunCommand("git sparse-checkout set --no-cone "+path, append(options, Dir(clonePath))...); err != nil {
		return err
	}

	if err := RunCommand("git checkout", append(options, Dir(clonePath))...); err != nil {
		return err
	}

	if err := cp.Copy(filepath.Join(cfg.Dir, ".clone", moduleName, path), filepath.Join(cfg.Dir, "modules", moduleName)); err != nil {
		return err
	}

	//if err := RunCommand("go get "+module, options...); err != nil {
	//	return err
	//}

	return nil
}

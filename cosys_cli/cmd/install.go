package cmd

import (
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	sb := strings.Builder{}
	for _, moduleName := range modules {
		sb.WriteString(" modules/" + moduleName)
	}

	if err := RunCommand("git clone --depth 1 --filter=blob:none --no-checkout git@github.com:cosys-io/cosys.git .clone", options...); err != nil {
		return err
	}

	defer os.RemoveAll(filepath.Join(cfg.Dir, ".clone"))

	if err := RunCommand("git sparse-checkout set --no-cone"+sb.String(), append(options, Dir(filepath.Join(cfg.Dir, ".clone")))...); err != nil {
		return err
	}

	if err := RunCommand("git checkout", append(options, Dir(filepath.Join(cfg.Dir, ".clone")))...); err != nil {
		return err
	}

	for _, moduleName := range modules {
		if err := cp.Copy(filepath.Join(cfg.Dir, ".clone", "modules", moduleName), filepath.Join(cfg.Dir, "modules", moduleName)); err != nil {
			return err
		}
	}

	//if err := RunCommand("go mod tidy", options...); err != nil {
	//	return err
	//}

	return nil
}

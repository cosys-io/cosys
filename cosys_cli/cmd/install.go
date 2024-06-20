package cmd

import (
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
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
		if err := installModules("", args...); err != nil {
			log.Fatal(err)
		}
	},
}

func installModules(dir string, modules ...string) error {
	cmdArgs := []string{
		"sparse-checkout",
		"set",
		"--no-cone",
	}
	for _, moduleName := range modules {
		cmdArgs = append(cmdArgs, "modules/"+moduleName)
	}

	if err := RunCommand(dir, "git", "clone", "--depth", "1", "--filter=blob:none", "--no-checkout", "git@github.com:cosys-io/cosys.git", ".clone"); err != nil {
		return err
	}

	defer os.RemoveAll(filepath.Join(dir, ".clone"))

	if err := RunCommand(filepath.Join(dir, ".clone"), "git", cmdArgs...); err != nil {
		return err
	}

	if err := RunCommand(filepath.Join(dir, ".clone"), "git", "checkout"); err != nil {
		return err
	}

	for _, moduleName := range modules {
		if err := cp.Copy(filepath.Join(dir, ".clone", "modules", moduleName), filepath.Join(dir, "modules", moduleName)); err != nil {
			return err
		}
	}

	if err := RunCommand(dir, "go", "mod", "tidy"); err != nil {
		return err
	}

	return nil
}

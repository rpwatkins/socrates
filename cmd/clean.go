package cmd

import (
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean deletes all build files.",
	Long:  `The clean command deletes all the build files in the build folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		clean()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func clean() {

	cwd, err := os.Getwd()
	if err != nil {
		log.Error("could not get curent directory")
		os.Exit(1)
	}

	buildDir := path.Join(cwd, "build")
	if err := removeContents(buildDir); err != nil {
		log.Error("build folder could not be cleaned")
		os.Exit(1)
	}

}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

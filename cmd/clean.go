package cmd

import (
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "<clean> deletes all build files.",
	Long:  `The <clean> command deletes all the build files in the builds folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		clean(afero.NewOsFs())
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func clean(fs afero.Fs) {

	buildDir := path.Join("build")

	exists, err := afero.Exists(fs, buildDir)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if !exists {
		log.Warning("no build folder exists.")
		os.Exit(1)
	}
	if err := removeContents(fs, buildDir); err != nil {
		log.Error(err)
		log.Error("build folder could not be cleaned")
		os.Exit(1)
	}

}

func removeContents(fs afero.Fs, dir string) error {
	names, err := afero.ReadDir(fs, dir)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

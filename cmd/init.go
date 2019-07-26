package cmd

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/gobuffalo/packr/v2"
	log "github.com/sirupsen/logrus"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init creates a new manuscript.",
	Long:  `The init command creates a new manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		initProject()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initProject() {

	// default file system
	var fs = afero.NewOsFs()

	// get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	// confirm current working directory is empty
	empty, err := afero.IsEmpty(fs, cwd)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	if !empty {
		log.Error("Folder not empty. A Socrates project can only be initialized in an empty folder.")
		os.Exit(1)
	}

	log.Info("Initializing Socrates project.")

	box := packr.New("assets", "./templates")

	for _, v := range InitPaths() {
		if err := fs.Mkdir(filepath.Join(cwd, v), 0755); err != nil {
			log.Error(err.Error())
		}
	}

	for k, v := range InitFileMap() {
		// get file from box
		file, err := box.Find(k)
		if err != nil {
			log.Error(err.Error())
		}
		// copy file to destination
		if err := afero.WriteReader(fs, filepath.Join(cwd, v, k), bytes.NewReader(file)); err != nil {
			log.Error(err.Error())
		}
	}
	log.Infof("Socrates project created at %s", cwd)

}

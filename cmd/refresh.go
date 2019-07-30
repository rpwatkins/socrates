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

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "<refresh> adds missing items to an existing manuscript.",
	Long:  `The <refresh> command searches for missing items and adds them to an existing manuscript. Existing files are not over-written.`,
	Run: func(cmd *cobra.Command, args []string) {
		refreshProject()
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}

func refreshProject() {

	// default file system
	var fs = afero.NewOsFs()

	// get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	// confirm current working directory is not empty
	empty, err := afero.IsEmpty(fs, cwd)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	if empty {
		log.Error("This folder is empty. Run <init> instead to create a new Socrates project.")
		os.Exit(1)
	}

	log.Info("Refreshing Socrates project.")

	box := packr.New("assets", "./templates")

	for _, v := range InitPaths() {
		exists, err := afero.DirExists(fs, filepath.Join(cwd, v))
		if err != nil {
			log.Error(err)
		}
		if exists {
			log.Warningf("%s folder exists", v)
		} else {
			if err := fs.Mkdir(filepath.Join(cwd, v), 0755); err != nil {
				log.Error(err.Error())
			}
			log.Infof("%s folder created", v)
		}
	}

	for k, v := range InitFileMap() {
		// get file from box
		exists, err := afero.Exists(fs, filepath.Join(v, k))
		if err != nil {
			log.Error(err.Error())
		}
		if exists {
			log.Warningf("%s file exists", filepath.Join(v, k))
		} else {
			file, err := box.Find(k)
			if err != nil {
				log.Error(err.Error())
			}
			// copy file to destination
			if err := afero.WriteReader(fs, filepath.Join(cwd, v, k), bytes.NewReader(file)); err != nil {
				log.Error(err.Error())
			}
			log.Infof("%s file created", filepath.Join(v, k))
		}

	}
	log.Infof("Socrates project at %s refreshed", cwd)

}

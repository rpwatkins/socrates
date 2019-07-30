package cmd

import (
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "<init> creates a new manuscript.",
	Long:  `The <init> command creates a new manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		// default file system
		var fs = afero.NewOsFs()
		initProject(fs)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initProject(fs afero.Fs) {

	// confirm current working directory is empty
	empty, err := afero.IsEmpty(fs, ".")
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	if !empty {
		log.Error("Folder not empty. A Socrates project can only be initialized in an empty folder.")
		os.Exit(1)
	}

	log.Info("Initializing Socrates project.")
	if err := writeFS(fs); err != nil {
		log.Error(err)
		log.Error("initilization failed.")
		os.Exit(1)
	}
	log.Info("Socrates project created.")
}

package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version displays the current version of Socrates.",
	Long:  `The version command displays the current Socrates version.`,
	Run: func(cmd *cobra.Command, args []string) {
		version()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func version() {
	log.Infof("Socrates v0.7.0")
}

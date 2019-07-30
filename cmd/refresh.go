package cmd

import (
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "<refresh> adds missing items to an existing manuscript.",
	Long:  `The <refresh> command searches for missing items and adds them to an existing manuscript. Existing files are not over-written.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := writeFS(afero.NewOsFs()); err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}

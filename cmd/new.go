package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [path]",
	Short: "new creates a new asciidoc file at the path entered.",
	Long: `The new command creates new asciidoc files for a manuscript. It will create a new asciidoc file at the path entered. For example:

		socrates new parts/part_01/chapters/chapter_02.adoc

	This will create chapter_02.adoc in the parts/part_01/chapters folder.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		newDoc(afero.NewOsFs(), args)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func newDoc(fs afero.Fs, args []string) {
	path := args[0]
	// check if exists
	exists, err := afero.Exists(fs, path)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if exists {
		log.Errorf("%s already exists..", path)
		os.Exit(1)
	}
	// copy file to destination
	if err := afero.WriteReader(fs, path, bytes.NewReader([]byte(""))); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	yellow := color.New(color.FgYellow).SprintFunc()

	log.Infof("%s created.", path)
	fmt.Print("\nThe following include directive should be added to where you want the include to appear in your manuscript:\n\n")
	inc := fmt.Sprintf("include::%s[]\n\n", path)
	fmt.Printf("%s", yellow(inc))
	fmt.Print("The path in the inclue directive may need to be edited depending on the file for which it is intended.")
}

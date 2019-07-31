package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

var includeCmd = &cobra.Command{
	Use:   "inc [part name] [chapter name] [name]",
	Short: "inc adds an include file to the part/chapter entered.",
	Long:  `The inc command adds a new include file to a part/chapter using the names entered.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		include(afero.NewOsFs(), args)
	},
}

func init() {
	rootCmd.AddCommand(includeCmd)
}

func include(fs afero.Fs, args []string) {
	part := args[0]
	chapter := args[1]
	name := args[2]

	// new file path
	fileName := fmt.Sprintf("include_%s_%s.adoc", chapter, name)
	path := filepath.Join("src", "parts", part)

	// check if already exists
	exists, err := afero.Exists(fs, filepath.Join(path, fileName))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if exists {
		log.Warningf("%s already exists.", fileName)
		os.Exit(1)
	} else {
		// create empty file
		file := []byte("")
		// write the file to disk
		if err := afero.WriteFile(fs, filepath.Join(path, fileName), file, 0644); err != nil {
			log.Error(err)
			os.Exit(1)
		}
		log.Infof("%s created in %s folder", fileName, path)
		fmt.Printf("\nCopy and paste the following include directive in the proper place in %s.adoc in %s.", chapter, part)
		fmt.Printf("\n\n     include::%s[]\n\n", fileName)
	}
}

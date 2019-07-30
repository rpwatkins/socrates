package cmd

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check analyzes the integrity of a master document.",
	Long:  `The check command analyzes the include directives in a master document for missing files or incorrect file names.`,
	Run: func(cmd *cobra.Command, args []string) {
		check(afero.NewOsFs())
	},
}

type master struct {
	Attributes map[string]string
	Includes   []string
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func check(fs afero.Fs) {
	// load master.adoc
	m, err := loadMaster(fs)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	// create a list of paths to check
	paths := []string{}
	// add location of referecnes.bib
	paths = append(paths, strings.TrimSpace(m.Attributes["bibliography-database"]))
	// parse include directives for broken references.bib
	for _, v := range m.Includes {
		// split at "::"
		res := strings.Split(v, "::")
		path := strings.Split(res[1], "[")
		paths = append(paths, filepath.Join("src", strings.TrimSpace(path[0])))
	}
	noErr := true
	for _, v := range paths {
		exists, err := afero.Exists(fs, v)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		if !exists {
			noErr = false
			log.Warningf("%s file missing.", v)
		} else {
			log.Infof("%s found", v)
		}
	}
	if !noErr {
		log.Warning("some files are missing.")
	} else {
		log.Info("all files found.")
	}

}

func loadMaster(fs afero.Fs) (master, error) {
	master := master{}

	file, err := afero.ReadFile(fs, filepath.Join("src", "master.adoc"))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(file)))
	lMap := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) >= 1 && line[:1] == ":" {
			// load attribute
			atts := strings.Split(line, ":")
			lMap[atts[1]] = atts[2]

		}
		master.Attributes = lMap
		if len(line) >= 7 && line[:7] == "include" {
			// load include statement
			master.Includes = append(master.Includes, line)
		}
	}
	return master, nil

}

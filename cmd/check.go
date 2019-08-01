package cmd

import (
	"bufio"
	"fmt"
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
		missing, err := check(afero.NewOsFs())
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		if len(missing) == 0 {
			log.Info("all included files found.")
		} else {
			for _, m := range missing {
				log.Warning(m)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func check(fs afero.Fs) ([]string, error) {

	missing := []string{}
	paths := parseMaster(fs, "master.adoc")
	childPaths := []string{}

	for _, p := range paths {
		childPaths = append(childPaths, parseChild(fs, p)...)
	}
	paths = append(paths, childPaths...)

	for _, v := range paths {
		exists, err := afero.Exists(fs, v)
		if err != nil {
			return nil, err
		}
		if !exists {
			missing = append(missing, fmt.Sprintf("%s file missing.", v))
		}
	}
	return missing, nil

}

func parseMaster(fs afero.Fs, file string) []string {

	paths := []string{}
	attributes := make(map[string]string)

	content, err := afero.ReadFile(fs, file)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) >= 1 && line[:1] == ":" {
			// load attribute
			atts := strings.Split(line, ":")
			attributes[atts[1]] = atts[2]
		}

		if len(line) >= 7 && line[:7] == "include" {
			res := strings.Split(line, "::")
			path := strings.Split(res[1], "[")
			p := strings.TrimSpace(path[0])
			paths = append(paths, p)
		}
	}
	paths = append(paths, strings.TrimSpace(attributes["bibliography-database"]))

	return paths

}

func parseChild(fs afero.Fs, file string) []string {
	extension := filepath.Ext(file)
	parentPath := file[0 : len(file)-len(extension)]
	paths := []string{}

	exists, err := afero.Exists(fs, file)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if exists {
		content, err := afero.ReadFile(fs, file)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(strings.NewReader(string(content)))

		for scanner.Scan() {
			line := scanner.Text()
			if len(line) >= 9 && line[:9] == "include::" {

				includeParts := strings.Split(line, "::")
				includePath := strings.Split(includeParts[1], "[")[0]
				paths = append(paths, filepath.Join(parentPath, includePath))
			}
		}
		// recurse paths
		childPaths := []string{}
		for _, p := range paths {
			childPaths = append(childPaths, parseChild(fs, p)...)
		}
		paths = append(paths, childPaths...)

	}

	return paths
}

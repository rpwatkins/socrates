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
		if _, err := check(afero.NewOsFs()); err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func check(fs afero.Fs) (bool, error) {
	// load master.adoc

	paths := parseMaster(fs, "master.adoc")
	childPaths := []string{}
	for _, p := range paths {
		childPaths = append(childPaths, parseChild(fs, p)...)
	}
	paths = append(paths, childPaths...)

	noErr := true
	for _, v := range paths {
		exists, err := afero.Exists(fs, v)
		if err != nil {
			return true, err
		}
		if !exists {
			noErr = false
			log.Warningf("%s file missing.", v)
		} else {
			log.Infof("%s found.", v)
		}
	}
	if !noErr {
		log.Warning("some files are missing.")
		return true, nil
	} else {
		log.Info("all included files found.")
		return false, nil
	}

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
			p := filepath.Join("src", strings.TrimSpace(path[0]))
			paths = append(paths, p)
		}
	}
	paths = append(paths, strings.TrimSpace(attributes["bibliography-database"]))

	return paths

}

func parseChild(fs afero.Fs, file string) []string {
	// recurse includes
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

				// get chapter name from parent
				pathParts := strings.Split(file, "/")
				chapter := pathParts[len(pathParts)-2]

				p := filepath.Join("src", "parts", chapter, includePath)
				paths = append(paths, p)
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

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
		check(afero.NewOsFs())
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

// helper functionm used before builds begin
func runValidation(fs afero.Fs) {

	_, missing, err := validate(fs)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if len(missing) > 0 {
		for _, m := range missing {
			log.Error(m)
		}
		log.Errorf("build failed due to invalid include directive(s).")
		os.Exit(1)
	}
}

// checks to make sure all included files can be found.
func check(fs afero.Fs) {
	found, missing, err := validate(fs)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if Verbose {
		for _, m := range found {
			log.Info(m)
		}
	}
	for _, m := range missing {
		log.Warning(m)
	}
	if len(missing) == 0 {
		log.Infof("all %d included files found.", len(found))
	} else {
		log.Warningf("%d included files missing.", len(missing))
	}
}

func validate(fs afero.Fs) ([]string, []string, error) {

	missing := []string{}
	found := []string{}

	paths := parseMaster(fs, "master.adoc")
	for k := range paths {
		childPaths := parseChild(fs, k)
		for j, u := range childPaths {
			paths[j] = u
		}
	}

	// check for all includes
	for k, v := range paths {
		exists, err := afero.Exists(fs, k)
		if err != nil {
			return nil, nil, err
		}
		if !exists {
			missing = append(missing, fmt.Sprintf("%s missing. Correct the include directive in %s.", k, v))
		} else {
			found = append(found, fmt.Sprintf("%s found.", k))
		}
	}
	return found, missing, nil
}

func parseMaster(fs afero.Fs, file string) map[string]string {

	paths := make(map[string]string)
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
			child := strings.TrimSpace(path[0])
			paths[child] = file
		}
	}
	paths[strings.TrimSpace(attributes["bibliography-database"])] = "."
	return paths

}

func parseChild(fs afero.Fs, file string) map[string]string {
	extension := filepath.Ext(file)
	parentPath := file[0 : len(file)-len(extension)]
	paths := make(map[string]string)

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
				paths[filepath.Join(parentPath, includePath)] = parentPath
			}
		}
		// recurse paths
		childPaths := make(map[string]string)
		for _, p := range paths {
			innerPaths := parseChild(fs, p)
			for k, v := range innerPaths {
				paths[k] = v
			}
		}
		for k, v := range childPaths {
			paths[k] = v
		}
	}
	return paths
}

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/plouc/textree"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list displays a tree diagram of a project's include directive structure. Those displayed in r3ed are missing.",
	Long:  `The list command displays a tree diagram of a project's include directive structure. Those displayed in r3ed are missing.`,
	Run: func(cmd *cobra.Command, args []string) {
		list(afero.NewOsFs())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

type include struct {
	Path     string
	Found    bool
	Includes []include
}

func list(fs afero.Fs) {

	incs := listMaster(fs, "master.adoc")

	green := color.New(color.FgHiGreen).SprintFunc()
	root := textree.NewNode(green("master.adoc"))
	displayIncludes(incs, root)
	o := textree.NewRenderOptions()
	root.Render(os.Stdout, o)

}

func displayIncludes(includes []include, parent *textree.Node) {
	for _, inc := range includes {
		p := strings.Split(inc.Path, "/")
		name := p[len(p)-1]
		red := color.New(color.FgRed).SprintFunc()
		blue := color.New(color.FgBlue).SprintFunc()
		if !inc.Found {
			name = red(fmt.Sprintf("%s (missing: %s)", name, inc.Path))
		} else {
			name = blue(name)
		}

		newNode := textree.NewNode(name)
		displayIncludes(inc.Includes, newNode)
		parent.Append(newNode)
	}
}

func listMaster(fs afero.Fs, file string) []include {

	// open master.adoc
	content, err := afero.ReadFile(fs, file)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// list of includes
	incs := []include{}
	// scan file line by line
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) >= 7 && line[:7] == "include" {
			res := strings.Split(line, "::")
			path := strings.Split(res[1], "[")
			child := strings.TrimSpace(path[0])
			inc := include{
				Path:  child,
				Found: false,
			}
			// check if exists
			exists, err := afero.Exists(fs, inc.Path)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			inc.Found = exists
			// get child includes
			inc.Includes = listChild(fs, inc.Path)
			incs = append(incs, inc)

		}
	}

	return incs

}

func listChild(fs afero.Fs, file string) []include {
	// get parent path from file
	n := filepath.Base(file)
	parentPath := strings.Replace(file, n, "", 1)
	incs := []include{}

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
				res := strings.Split(line, "::")
				path := strings.Split(res[1], "[")
				child := strings.TrimSpace(path[0])
				inc := include{
					Path:  filepath.Join(parentPath, child),
					Found: false,
				}
				// check if exists
				exists, err := afero.Exists(fs, inc.Path)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
				inc.Found = exists
				// recurse
				inc.Includes = listChild(fs, inc.Path)
				incs = append(incs, inc)

			}
		}
	}
	return incs
}

func runValidation(fs afero.Fs) []string {

	missing := []string{}
	incs := listMaster(fs, "master.adoc")
	for _, i := range incs {
		if !i.Found {
			missing = append(missing, i.Path)
		}
	}
	return missing
}

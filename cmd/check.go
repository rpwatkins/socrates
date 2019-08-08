package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"

	"github.com/plouc/textree"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"mvdan.cc/xurls/v2"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check validates a project's block macros and includes for missing files and tests URLs.",
	Long:  `The check command validates a project's block macros and include directives for missing files and it also test URLs.`,
	Run: func(cmd *cobra.Command, args []string) {
		check(afero.NewOsFs())
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

type include struct {
	Path     string
	Found    bool
	Includes []include
	Kind     string
	LineNum  int
}

func check(fs afero.Fs) {

	incs := checkMaster(fs, "master.adoc")
	f, m := flattenAndSortByMissingIncludes(incs)
	// prepare summary
	missingCount := len(m)
	includeCount := 0
	attributeCount := 0
	imageCount := 0
	urlCount := 0
	diagramCount := 0
	for _, i := range f {
		switch i.Kind {
		case "include":
			includeCount += 1
		case "attribute":
			attributeCount += 1
		case "image":
			imageCount += 1
		case "url":
			urlCount += 1
		case "diagram":
			diagramCount += 1
		}
	}

	// show tree
	red := color.New(color.FgRed).SprintFunc()         // missing
	blue := color.New(color.FgBlue).SprintFunc()       // includes
	yellow := color.New(color.FgYellow).SprintFunc()   // attributes
	cyan := color.New(color.FgCyan).SprintFunc()       // images
	magenta := color.New(color.FgMagenta).SprintFunc() // urls
	green := color.New(color.FgHiGreen).SprintFunc()   // diagrams

	fmt.Print("\nSUMMARY:   ")
	fmt.Print(red(fmt.Sprintf("%d missing   ", missingCount)))
	fmt.Print(blue(fmt.Sprintf("%d include   ", includeCount)))
	fmt.Print(yellow(fmt.Sprintf("%d attribute   ", attributeCount)))
	fmt.Print(cyan(fmt.Sprintf("%d image   ", imageCount)))
	fmt.Print(magenta(fmt.Sprintf("%d url   ", urlCount)))
	fmt.Print(green(fmt.Sprintf("%d diagram\n", diagramCount)))

	// prepare tree display
	root := textree.NewNode("master.adoc")
	// get all child nodes
	display(incs, root)
	o := textree.NewRenderOptions()
	root.Render(os.Stdout, o)

}

func display(includes []include, parent *textree.Node) {
	red := color.New(color.FgRed).SprintFunc()         // missing
	blue := color.New(color.FgBlue).SprintFunc()       // includes
	yellow := color.New(color.FgYellow).SprintFunc()   // attributes
	cyan := color.New(color.FgCyan).SprintFunc()       // images
	magenta := color.New(color.FgMagenta).SprintFunc() // urls
	green := color.New(color.FgHiGreen).SprintFunc()   // diagrams

	for _, inc := range includes {
		name := filepath.Base(inc.Path)

		if !inc.Found {
			if inc.Kind == "attribute" {
				name = red(fmt.Sprintf("%s (line %d: attribute file missing %s)", name, inc.LineNum, inc.Path))
			} else {
				name = red(fmt.Sprintf("%s (line %d: file missing %s)", name, inc.LineNum, inc.Path))
			}
		} else {
			switch inc.Kind {
			case "include":
				name = blue(name)
			case "image":
				name = cyan(name)
			case "url":
				name = magenta(name)
			case "diagram":
				name = green(name)
			case "attribute":
				name = yellow(name)
			}
		}
		newNode := textree.NewNode(name)
		display(inc.Includes, newNode)
		parent.Append(newNode)
	}

}

func checkMaster(fs afero.Fs, file string) []include {

	imagePath, err := getImagePath(fs, "master.adoc")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
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
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		// check for includes
		if strings.HasPrefix(line, "include::") {
			res := strings.Split(line, "::")
			path := strings.Split(res[1], "[")
			child := strings.TrimSpace(path[0])
			if !includesContains(incs, child) {
				inc := include{
					Path:    child,
					Found:   false,
					Kind:    "include",
					LineNum: lineNum,
				}
				// check if exists
				exists, err := afero.Exists(fs, inc.Path)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
				inc.Found = exists
				inc.LineNum = lineNum
				// get child includes
				inc.Includes = checkChild(fs, inc.Path, imagePath)
				incs = append(incs, inc)
			}
		}
		if strings.HasPrefix(line, ":") {
			parts := strings.Split(line, ":")
			k := strings.TrimSpace(parts[1])
			v := strings.TrimSpace(parts[2])

			if k == "bibliography-database" {
				exists, err := afero.Exists(fs, v)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
				i := include{
					Path:    strings.TrimSpace(v),
					Found:   exists,
					LineNum: lineNum,
					Kind:    "attribute",
				}
				incs = append(incs, i)
			}
			pdfPaths := []string{
				"front-cover-image",
				"page-background-image",
				"page-background-image-recto",
				"title-logo-image",
				"title-page-background-image",
			}
			for _, path := range pdfPaths {
				if k == path {
					exists, err := afero.Exists(fs, filepath.Join(imagePath, v))
					if err != nil {
						log.Error(err)
						os.Exit(1)
					}
					i := include{
						Path:    filepath.Join(imagePath, v),
						Found:   exists,
						LineNum: lineNum,
						Kind:    "attribute",
					}
					incs = append(incs, i)
				}
			}
		}
		lineNum += 1
	}
	return incs
}

func checkChild(fs afero.Fs, file string, imagePath string) []include {
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
		lineNum := 1

		for scanner.Scan() {
			// get includes
			line := scanner.Text()

			if strings.HasPrefix(line, "include::") {

				res := strings.Split(line, "::")
				path := strings.TrimSpace(strings.Split(res[1], "[")[0])

				if !includesContains(incs, path) {
					inc := checkItem(fs, filepath.Join(parentPath, path))
					inc.Kind = "include"
					inc.LineNum = lineNum
					// recurse
					inc.Includes = checkChild(fs, inc.Path, imagePath)
					incs = append(incs, inc)
				}
			}
			// images
			if strings.HasPrefix(line, "image::") {
				parts := strings.Split(line, "::")
				path := strings.Split(parts[1], "[")[0]
				// check if exists
				inc := checkItem(fs, filepath.Join(imagePath, path))
				inc.Kind = "image"
				inc.LineNum = lineNum
				incs = append(incs, inc)
			}

			macros := []string{
				"a2s::",
				"actdiag::",
				"blockdiag::",
				"ditaa::",
				"erd::",
				"graphviz::",
				"meme::",
				"mermaid::",
				"msc::",
				"nomnoml::",
				"nwdiag::",
				"packetdiag::",
				"plantuml::",
				"rackdiag::",
				"seqdiag::",
				"shaape::",
				"svgbob::",
				"syntrax::",
				"umlet::",
				"vega::",
				"vegalite::",
				"wavedrom::",
			}
			for _, d := range macros {
				if strings.HasPrefix(line, d) {
					parts := strings.Split(line, "::")
					path := strings.Split(parts[1], "[")[0]
					inc := checkItem(fs, path)
					inc.Kind = "diagram"
					inc.LineNum = lineNum
					incs = append(incs, inc)
				}
			}

			urls := xurls.Strict().FindAllString(line, -1)
			for _, url := range urls {
				if strings.Contains(url, "[") {
					url = strings.Split(url, "[")[0]
				}
				inc := checkURL(url)
				inc.Kind = "url"
				inc.LineNum = lineNum
				incs = append(incs, inc)
			}

			// inline images TODO:
			regex := regexp.MustCompile(`image:[^:](.+?)\[`)
			matches := regex.FindAll([]byte(line), -1)
			for _, v := range matches {
				p := string(v)
				path := strings.Split(p, ":")[1]
				path = strings.Split(path, "[")[0]

				inc := checkItem(fs, filepath.Join(imagePath, path))
				inc.Kind = "image"
				inc.LineNum = lineNum
				incs = append(incs, inc)
			}
			lineNum += 1
		}
	}
	return incs
}

// check for item and return include with status
func checkItem(fs afero.Fs, path string) include {
	inc := include{}
	inc.Path = path
	exists, err := afero.Exists(fs, inc.Path)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	inc.Found = exists
	return inc
}

// check URL works
func checkURL(path string) include {
	inc := include{}
	inc.Path = path
	resp, err := http.Get(path)
	if err != nil {
		inc.Found = false
	} else {
		if resp.Status == "200 OK" {
			inc.Found = true
		} else {
			inc.Found = false
		}
	}
	return inc
}

// prevents duplication of includes
func includesContains(includes []include, path string) bool {
	found := false
	for _, v := range includes {
		if v.Path == path {
			found = true
		}
	}
	return found
}

// flattens the hierarchy of include returned by listMaster into two flat lists of missing/dound items
func flattenAndSortByMissingIncludes(includes []include) ([]include, []include) {
	found := []include{}
	missing := []include{}
	for _, v := range includes {
		if v.Found {
			found = append(found, v)
		} else {
			missing = append(missing, v)
		}
		f, m := flattenAndSortByMissingIncludes(v.Includes)
		missing = append(missing, m...)
		found = append(found, f...)
	}
	return found, missing
}

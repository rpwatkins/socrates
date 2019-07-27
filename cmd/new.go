package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [item]",
	Short: "new creates a new element of a manuscript.",
	Long: `The new command creates elements of a manuscript. Supported values"

	Single instance elements:

	abstract
	bibliography
	colophon
	dedication
	glossary
	index
	preface

	Autonumbered elements
	part
	chapter
	appendix
	`,
	Run: func(cmd *cobra.Command, args []string) {
		initProject()
	},
}

// new part command
var newPartCmd = &cobra.Command{
	Use:   "part",
	Short: "create a new manuscript part.",
	Long:  `The <new part> command creates a new auto-numbered part for a manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		newPart()
	},
}

// new chapter command
var newChapterCmd = &cobra.Command{
	Use:   "chapter",
	Short: "<chapter> creates a new chapter in the part entered.",
	Long:  `The <new chapter> command creates a new auto-numbered chapter for (a part of) a manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		newChapter()
	},
}

// new appendix command
var newAppendixCmd = &cobra.Command{
	Use:   "appendix",
	Short: "create a new appendix.",
	Long:  `The <new appendix> command creates a new autonumbered appendix for a manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		newAppendix()
	},
}

// new abstract command
var newAbstractCmd = &cobra.Command{
	Use:   "abstract",
	Short: "<abstract> creates a new abstract.",
	Long:  `The <new abstract> command creates a new abstract for a manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		newItem("abstract")
	},
}

// new bibiliography command
var newBibliographyCmd = &cobra.Command{
	Use:   "bibliography",
	Short: "<bibliography> creates a new bibliography.",
	Long:  `The <new bib> command creates a new bibliography for a manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		newItem("bibliography")
	},
}

// new bibliography command
var newColophonCmd = &cobra.Command{
	Use:   "colophon",
	Short: "<colophon> creates a new colophon.",
	Long:  `The <new colophon> command creates a new colophon for a manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		newItem("colophon")
	},
}

// new dedication command
var newDedicationCmd = &cobra.Command{
	Use:   "dedication",
	Short: "<dedication> creates a new dedication.",
	Long:  `The <new dedication> command creates a new dedication for a manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		newItem("dedication")
	},
}

// new glossary command
var newGlossaryCmd = &cobra.Command{
	Use:   "glossary",
	Short: "<glossary> creates a new glossary.",
	Long:  `The <new glossaru> command creates a new glossary for a manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		newItem("glossary")
	},
}

// new index command
var newIndexCmd = &cobra.Command{
	Use:   "index",
	Short: "<index> creates a new index.",
	Long:  `The <new index> command creates a new index for a manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		newItem("index")
	},
}

// new preface command
var newPrefaceCmd = &cobra.Command{
	Use:   "preface",
	Short: "<preface> creates a new index.",
	Long:  `The <new preface> command creates a new preface for a manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		newItem("preface")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// subcommands -- part, chapter, bibliography, appendix, index, glossary, colophon, dedication, appendix, abstract, preface
	newCmd.AddCommand(newPartCmd)
	newCmd.AddCommand(newChapterCmd)
	newCmd.AddCommand(newAppendixCmd)
	newCmd.AddCommand(newAbstractCmd)
	newCmd.AddCommand(newBibliographyCmd)
	newCmd.AddCommand(newColophonCmd)
	newCmd.AddCommand(newDedicationCmd)
	newCmd.AddCommand(newGlossaryCmd)
	newCmd.AddCommand(newIndexCmd)
	newCmd.AddCommand(newPrefaceCmd)
}

func newPart() {

	// default file system
	var fs = afero.NewOsFs()

	// find existing parts to get new part number
	files, err := afero.ReadDir(fs, filepath.Join("src", "chapters"))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	parts := []int{}

	for _, f := range files {
		file := f.Name()
		if file[2:7] == "_part" {
			value, err := strconv.Atoi(file[:2])
			if err != nil {
				log.Error(err)
			}
			parts = append(parts, value)
		}
	}

	// get most recent part number
	var largerNumber, temp int
	for _, element := range parts {
		if element > temp {
			temp = element
			largerNumber = temp
		}
	}
	newNum := largerNumber + 1
	newName := ""
	if newNum < 10 {
		newName = fmt.Sprintf("0%d_part.adoc", newNum)
	} else {
		newName = fmt.Sprintf("%d_part.adoc", newNum)
	}

	// check if exists
	p := filepath.Join("src", "chapters", newName)
	exists, err := afero.Exists(fs, p)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if exists {
		log.Errorf("%s already exists", newName)
		os.Exit(1)
	}

	// get part template and dave it as new part in chapters folder
	box := packr.New("assets", "./templates")
	file, err := box.Find("part_01.adoc")
	if err != nil {
		log.Error(err.Error())
	}

	// run through plush and add chapter title
	ctx := plush.NewContext()
	ctx.Set("number", NumberToWord(newNum))

	s, err := plush.Render(string(file), ctx)
	if err != nil {
		log.Fatal(err)
	}
	s2 := []byte(s)

	// copy file to destination
	if err := afero.WriteReader(fs, filepath.Join("src", "chapters", newName), bytes.NewReader(s2)); err != nil {
		log.Error(err.Error())
	}
	log.Infof("%s created at src/chapters", newName)
	fmt.Printf(`Please add the following inlude directive to the proper place in your master.adoc file.

	include::chapters/%s[]`, newName)

}

func newAppendix() {

	// default file system
	var fs = afero.NewOsFs()

	// find existing parts to get new part number
	files, err := afero.ReadDir(fs, filepath.Join("src", "back_matter"))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	parts := []int{}

	for _, f := range files {
		file := f.Name()
		if file[2:11] == "_appendix" {
			value, err := strconv.Atoi(file[:2])
			if err != nil {
				log.Error(err)
			}
			parts = append(parts, value)
		}
	}

	// get most recent part number
	var largerNumber, temp int
	for _, element := range parts {
		if element > temp {
			temp = element
			largerNumber = temp
		}
	}

	newNum := largerNumber + 1
	newName := ""
	if newNum < 10 {
		newName = fmt.Sprintf("0%d_appendix.adoc", newNum)
	} else {
		newName = fmt.Sprintf("%d_appendix.adoc", newNum)
	}

	// check if exists
	p := filepath.Join("src", "back_matter", newName)
	exists, err := afero.Exists(fs, p)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if exists {
		log.Errorf("%s already exists", newName)
		os.Exit(1)
	}

	// get part template and dave it as new part in chapters folder
	box := packr.New("assets", "./templates")
	file, err := box.Find("appendix.adoc.plush")
	if err != nil {
		log.Error(err.Error())
	}

	// run through plush and add chapter title
	ctx := plush.NewContext()
	ctx.Set("number", NumberToWord(newNum))

	s, err := plush.Render(string(file), ctx)
	if err != nil {
		log.Fatal(err)
	}
	s2 := []byte(s)

	// copy file to destination
	if err := afero.WriteReader(fs, filepath.Join("src", "back_matter", newName), bytes.NewReader(s2)); err != nil {
		log.Error(err.Error())
	}
	log.Infof("%s created at src/back_matter", newName)
	fmt.Printf(`Please add the following inlude directive to the proper place in your master.adoc file.

	include::back_matter/%s[]`, newName)

}

func newChapter() {

	// default file system
	var fs = afero.NewOsFs()

	// find existing parts to get new part number
	files, err := afero.ReadDir(fs, filepath.Join("src", "chapters"))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	parts := []int{}

	for _, f := range files {
		file := f.Name()
		if file[2:10] == "_chapter" {
			value, err := strconv.Atoi(file[:2])
			if err != nil {
				log.Error(err)
			}
			parts = append(parts, value)
		}
	}

	// get most recent part number
	var largerNumber, temp int
	for _, element := range parts {
		if element > temp {
			temp = element
			largerNumber = temp
		}
	}
	newNum := largerNumber + 1
	newName := ""
	newFolder := ""
	if newNum < 10 {
		newName = fmt.Sprintf("0%d_chapter.adoc", newNum)
		newFolder = fmt.Sprintf("0%d_chapter", newNum)
	} else {
		newName = fmt.Sprintf("%d_chapter.adoc", newNum)
		newFolder = fmt.Sprintf("%d_chapter", newNum)
	}

	// check if exists
	p := filepath.Join("src", "chapters", newFolder, newName)
	exists, err := afero.Exists(fs, p)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if exists {
		log.Errorf("%s already exists", newName)
		os.Exit(1)
	}

	// get part template and dave it as new part in chapters folder
	box := packr.New("assets", "./templates")
	file, err := box.Find("chapter.adoc.plush")
	if err != nil {
		log.Error(err.Error())
	}
	// run through plush and add chapter title
	ctx := plush.NewContext()
	ctx.Set("number", NumberToWord(newNum))

	s, err := plush.Render(string(file), ctx)
	if err != nil {
		log.Fatal(err)
	}
	s2 := []byte(s)

	// copy file to destination
	if err := afero.WriteReader(fs, filepath.Join("src", "chapters", newFolder, newName), bytes.NewReader(s2)); err != nil {
		log.Error(err.Error())
	}

	log.Infof("%s created at src/chapters/%s", newName, newFolder)
	fmt.Printf(`Please add the following inlude directive to the proper place in your master.adoc file.

	include::chapters/%s[]`, newName)
}

func newItem(name string) {

	// default file system
	var fs = afero.NewOsFs()

	// set save path
	path := ""
	includePath := ""
	newName := fmt.Sprintf("%s.adoc", name)

	if name == "abstract" || name == "dedication" || name == "preface" {
		path = filepath.Join("src", "front_matter")
		includePath = "front_matter"
	} else {
		path = filepath.Join("src", "back_matter")
		includePath = "back_matter"
	}

	// check if exists
	exists, err := afero.Exists(fs, filepath.Join(path, newName))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if exists {
		log.Errorf("%s already exists", name)
		os.Exit(1)
	}
	// get part template and dave it as new part in chapters folder
	box := packr.New("assets", "./templates")
	file, err := box.Find(fmt.Sprintf("%s.adoc", name))
	if err != nil {
		log.Error(err.Error())
	}

	// copy file to destination
	if err := afero.WriteReader(fs, filepath.Join(path, newName), bytes.NewReader(file)); err != nil {
		log.Error(err.Error())
	}

	log.Infof("%s created at %s", name+".adoc", path)
	fmt.Printf(`
	Please add the following inlude directive to the proper place in your master.adoc file.

	include::%s/%s[]`, includePath, newName)

}

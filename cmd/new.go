package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [item]",
	Short: "new creates a new element of a manuscript.",
	Long: `The new command creates elements of a manuscript. Supported [item] values:"

	Single instance elements:

	abstract
	bibliography
	colophon
	dedication
	glossary
	index
	preface

	Auto-numbered elements

	part
	chapter
	appendix
	`,
}

// new part command
var newPartCmd = &cobra.Command{
	Use:   "part [name]",
	Short: "create a new manuscript part.",
	Long:  `The <new part> command creates a new manuscript part using the name entered.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// default file system
		newPart(afero.NewOsFs(), args)
	},
}

// new chapter command
var newChapterCmd = &cobra.Command{
	Use:   "chapter [name] [part name]",
	Short: "<chapter> creates a new chapter in the part entered.",
	Long:  `The <new chapter [name] [part name]> command creates a new chapter to manuscript part.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// default file system
		newChapter(afero.NewOsFs(), args)
	},
}

// new appendix command
var newAppendixCmd = &cobra.Command{
	Use:   "appendix [name]",
	Short: "create a new appendix.",
	Long:  `The <new appendix [name]> command creates a new appendix for a manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		// default file system
		newAppendix(afero.NewOsFs(), args)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// subcommands -- part, chapter, bibliography, appendix, index, glossary, colophon, dedication, appendix, abstract, preface
	newCmd.AddCommand(newPartCmd)
	newCmd.AddCommand(newChapterCmd)
	newCmd.AddCommand(newAppendixCmd)

}

func newPart(fs afero.Fs, args []string) {
	name := args[0]
	// check if exists
	path := filepath.Join("src", "parts", name)
	exists, err := afero.Exists(fs, filepath.Join(path, name+".adoc"))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if exists {
		log.Errorf("%s already exists", name+".adoc")
		os.Exit(1)
	}
	// crerate folder
	if err := fs.Mkdir(path, 0755); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	include := fmt.Sprintf("\n// TODO: move me\n\ninclude::parts/%s/%s.adoc[]", name, name)

	if err := createItem("part", name, path, include, fs); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func newChapter(fs afero.Fs, args []string) {
	name := args[0]
	part := args[1]
	// check if exists
	path := filepath.Join("src", "parts", part)
	log.Debug(path)
	exists, err := afero.Exists(fs, filepath.Join(path, name+".adoc"))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if exists {
		log.Errorf("%s already exists", name+".adoc")
		os.Exit(1)
	}
	include := fmt.Sprintf("\n// TODO: move me\n\ninclude::parts/%s/%s.adoc[]", part, name)

	if err := createItem("chapter", name, path, include, fs); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func newAppendix(fs afero.Fs, args []string) {
	name := args[0]
	// check if exists
	path := filepath.Join("src", "back_matter")
	exists, err := afero.Exists(fs, filepath.Join(path, name+".adoc"))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if exists {
		log.Errorf("%s already exists", name+".adoc")
		os.Exit(1)
	}
	include := fmt.Sprintf("\n// TODO: move me\n\ninclude::back_matter/%s.adoc[]", name)

	if err := createItem("appendix", name, path, include, fs); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func createItem(item, name, dest, include string, fs afero.Fs) error {

	// get part template and dave it as new part in chapters folder
	box := packr.New("assets", "./templates")
	file, err := box.Find(fmt.Sprintf("%s.adoc.plush", item))
	if err != nil {
		return err
	}

	// run through plush and add chapter title
	ctx := plush.NewContext()
	ctx.Set("title", name)

	s, err := plush.Render(string(file), ctx)
	if err != nil {
		return err
	}
	s2 := []byte(s)

	// copy file to destination
	if err := afero.WriteReader(fs, filepath.Join(dest, name+".adoc"), bytes.NewReader(s2)); err != nil {
		return err
	}
	// add new item to end of master.adoc
	master, err := afero.ReadFile(fs, filepath.Join("src", "master.adoc"))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	m := string(master)
	n := fmt.Sprintf("%s\n\n%s", m, include)

	err2 := afero.WriteFile(fs, filepath.Join("src", "master.adoc"), []byte(n), 0777)
	if err2 != nil {
		log.Error(err2)
		os.Exit(1)
	}
	log.Infof("%s.adoc created at %s", name, dest)
	fmt.Printf(`The following include directive has been added to the end of master.adoc:
		%s
		
Please move it to the correct place in master.adoc.`, include)
	return nil
}

package cmd

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init creates a new manuscript.",
	Long:  `The init command creates a new manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		// default file system
		var fs = afero.NewOsFs()
		initProject(fs)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initProject(fs afero.Fs) {

	// confirm current working directory is empty
	empty, err := afero.IsEmpty(fs, ".")
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	if !empty {
		log.Error("Folder not empty. A Socrates project can only be initialized in an empty folder.")
		os.Exit(1)
	}

	// write folders and files
	if err := writeFS(fs); err != nil {
		log.Error(err)
		log.Error("initilization failed.")
		os.Exit(1)
	}
	log.Info("Socrates project created.")
}

func InitPaths() []string {
	return []string{
		"src",
		filepath.Join("src", "front_matter"),
		filepath.Join("src", "back_matter"),
		filepath.Join("src", "images"),
		filepath.Join("src", "assets"),
		filepath.Join("src", "resources"),
		filepath.Join("src", "resources", "pdfstyles"),
		filepath.Join("src", "parts"),
		filepath.Join("src", "parts", "01_part"),
	}
}

func InitFileMap() map[string]string {
	// create a map from a file name to a destination fold
	src := "src"
	back := filepath.Join(src, "back_matter")
	front := filepath.Join(src, "front_matter")
	m := make(map[string]string)
	m["appendix.adoc.plush"] = back
	m["bibliography.adoc"] = back
	m["colophon.adoc"] = back
	m["glossary.adoc"] = back
	m["index.adoc"] = back
	m["preface.adoc"] = front
	m["dedication.adoc"] = front
	m["abstract.adoc"] = front
	m["master.adoc"] = src
	m["references.bib"] = src
	m["chapter.adoc.plush"] = filepath.Join(src, "parts", "01_part")
	m["default-theme.yml"] = filepath.Join(src, "resources", "pdfstyles")
	m["part.adoc.plush"] = filepath.Join(src, "parts", "01_part")
	return m
}

func writeFS(fs afero.Fs) error {
	box := packr.New("assets", "./templates")

	for _, v := range InitPaths() {
		exists, err := afero.DirExists(fs, v)
		if err != nil {
			return err
		}
		if exists {
			log.Warningf("%s folder exists", v)
		} else {
			if err := fs.Mkdir(v, 0755); err != nil {
				return err
			}
			log.Infof("%s folder created", v)
		}
	}

	for k, v := range InitFileMap() {
		exists, err := afero.Exists(fs, filepath.Join(v, k))
		if err != nil {
			return err
		}
		if exists {
			log.Warningf("%s file exists", filepath.Join(v, k))
		} else {
			// get file from box
			file, err := box.Find(k)
			if err != nil {
				return err
			}
			if k[len(k)-5:] == "plush" {
				// run through plush with number = 1
				title := ""
				if k[:8] == "appendix" {
					title = "Appendix"
				} else if k[:7] == "chapter" {
					title = "Chapter"
				} else if k[0:4] == "part" {
					title = "Part"
				}

				ctx := plush.NewContext()
				ctx.Set("title", title)

				s, err := plush.Render(string(file), ctx)
				if err != nil {
					return err
				}
				s2 := []byte(s)

				// get name of file with .plush extension
				extension := filepath.Ext(k)
				name := "01_" + k[0:len(k)-len(extension)]

				// copy file to destination
				if err := afero.WriteReader(fs, filepath.Join(v, name), bytes.NewReader(s2)); err != nil {
					return err
				}
			} else {
				// copy file to destination
				if err := afero.WriteReader(fs, filepath.Join(v, k), bytes.NewReader(file)); err != nil {
					return err
				}

			}
			log.Infof("%s file created.", v)
		}
	}
	return nil
}

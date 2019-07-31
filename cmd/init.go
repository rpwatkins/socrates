package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

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

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "<refresh> adds missing items to an existing manuscript.",
	Long:  `The <refresh> command searches for missing items and adds them to an existing manuscript. Existing files are not over-written.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := writeFS(afero.NewOsFs()); err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(refreshCmd)

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
		"front_matter",
		"back_matter",
		"images",
		"assets",
		"resources",
		filepath.Join("resources", "pdfstyles"),
		"parts",
		filepath.Join("parts", "part_01"),
		filepath.Join("parts", "part_01", "chapters"),
	}
}

func InitFileMap() map[string]string {
	m := make(map[string]string)
	m["appendix.adoc.plush"] = "back_matter"
	m["bibliography.adoc"] = "back_matter"
	m["colophon.adoc"] = "back_matter"
	m["glossary.adoc"] = "back_matter"
	m["index.adoc"] = "back_matter"
	m["preface.adoc"] = "front_matter"
	m["dedication.adoc"] = "front_matter"
	m["abstract.adoc"] = "front_matter"
	m["master.adoc"] = "."
	m["references.bib"] = "."
	m["chapter.adoc.plush"] = filepath.Join("parts", "part_01", "chapters")
	m["default-theme.yml"] = filepath.Join("resources", "pdfstyles")
	m["part.adoc.plush"] = filepath.Join("parts", "part_01")
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
			log.Warningf("%s file exists.", filepath.Join(v, k))
		} else {
			// get file from box
			file, err := box.Find(k)
			if err != nil {
				return err
			}
			if k[len(k)-5:] == "plush" {
				// run through plush with number = 1

				extension := filepath.Ext(k)
				oldName := k[0 : len(k)-len(extension)]

				title := ""
				name := ""
				if k[:8] == "appendix" {
					title = "Appendix"
					name = strings.Replace(oldName, "appendix", "appendix_01", 1)
				} else if k[:7] == "chapter" {
					title = "Chapter"
					name = strings.Replace(oldName, "chapter", "chapter_01", 1)
				} else if k[0:4] == "part" {
					title = "Part"
					name = strings.Replace(oldName, "part", "part_01", 1)
				}

				ctx := plush.NewContext()
				ctx.Set("title", title)

				s, err := plush.Render(string(file), ctx)
				if err != nil {
					return err
				}
				s2 := []byte(s)
				// copy file to destination
				if err := afero.WriteReader(fs, filepath.Join(v, name), bytes.NewReader(s2)); err != nil {
					return err
				}
				log.Infof("%s file created.", filepath.Join(v, name))
			} else {
				// copy file to destination
				if err := afero.WriteReader(fs, filepath.Join(v, k), bytes.NewReader(file)); err != nil {
					return err
				}
				log.Infof("%s file created.", filepath.Join(v, k))

			}

		}
	}
	return nil
}

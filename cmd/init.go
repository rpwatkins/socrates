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

var bare bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init creates a new manuscript.",
	Long:  `The init command creates a new manuscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		initProject(afero.NewOsFs())
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&bare, "bare", "b", false, "create a bare project.")
}

func initProject(fs afero.Fs) {

	empty, err := afero.IsEmpty(fs, ".")
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	if !empty {
		log.Error("Folder not empty. A Socrates project can only be initialized in an empty folder.")
		os.Exit(1)
	}

	// confirm current working directory is empty
	if !bare {
		if err := writeFS(fs); err != nil {
			log.Error(err)
			log.Error("initilization failed.")
			os.Exit(1)
		}
	} else {
		if err := writeBare(fs); err != nil {
			log.Error(err)
			log.Error("initilization failed.")
			os.Exit(1)
		}
	}
	log.Info("Socrates project created.")
}

// bare socrates project
func writeBare(fs afero.Fs) error {
	box := packr.New("assets", "./templates")
	file, err := box.Find("master-bare.adoc")
	if err != nil {
		return err
	}
	// copy file to destination
	if err := afero.WriteFile(fs, master, file, 0644); err != nil {
		return err
	}
	file2, err := box.Find("references.bib")
	if err != nil {
		return err
	}
	// copy file to destination
	if err := afero.WriteFile(fs, "references.bib", file2, 0644); err != nil {
		return err
	}
	if Verbose {
		log.Infof("%s created.", master)
		log.Infof("references.bib created.")
	}
	return nil

}

// default init below
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
		filepath.Join("parts", "part_01", "chapters", "chapter_01"),
	}
}

func InitFileMap() map[string]string {
	b := "back_matter"
	f := "front_matter"
	m := make(map[string]string)
	m["appendix.adoc.plush"] = b
	m["bibliography.adoc"] = b
	m["colophon.adoc"] = b
	m["glossary.adoc"] = b
	m["index.adoc"] = b
	m["preface.adoc"] = f
	m["dedication.adoc"] = f
	m["abstract.adoc"] = f
	m[master] = "."
	m["placeholder.jpg"] = "images"
	m["references.bib"] = "."
	m["socrates.toml.plush"] = "."
	m["chapter.adoc.plush"] = filepath.Join("parts", "part_01", "chapters", "chapter_01")
	m["include_01.adoc"] = filepath.Join("parts", "part_01", "chapters", "chapter_01")
	m["include_02.adoc"] = filepath.Join("parts", "part_01", "chapters", "chapter_01")
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
			if Verbose {
				log.Infof("%s folder created.", v)
			}
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
				} else if k == "socrates.toml.plush" {
					cwd, err := os.Getwd()
					if err != nil {
						log.Error(err)
						os.Exit(1)
					}
					title = filepath.Base(cwd)
					name = oldName
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
				if Verbose {
					log.Infof("%s created.", filepath.Join(v, name))
				}
			} else {
				// copy file to destination
				if err := afero.WriteReader(fs, filepath.Join(v, k), bytes.NewReader(file)); err != nil {
					return err
				}
				if Verbose {
					log.Infof("%s created.", filepath.Join(v, k))
				}
			}
		}
	}
	return nil
}

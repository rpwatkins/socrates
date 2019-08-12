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
		filepath.Join("resources", "htmlstyles"),
		filepath.Join("resources", "htmlstyles", "asciidoctor-bs-themes"),
		filepath.Join("resources", "htmlstyles", "asciidoctor-skins"),
		filepath.Join("resources", "htmlstyles", "stylesheet-factory"),
		filepath.Join("resources", "htmlstyles", "stylesheets-bulma"),
	}
}

func InitFileMap() map[string]string {
	b := "back_matter"
	f := "front_matter"
	res := "resources"

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
	m["htmlstyles/asciidoctor-bs-themes/bootstrap_cerulean.css"] = res
	m["htmlstyles/asciidoctor-bs-themes/bootstrap_cerulean.min.css"] = res
	m["htmlstyles/asciidoctor-bs-themes/bootstrap_custom.css"] = res
	m["htmlstyles/asciidoctor-bs-themes/bootstrap_custom.min.css"] = res
	m["htmlstyles/asciidoctor-bs-themes/bootstrap_default_themed.css"] = res
	m["htmlstyles/asciidoctor-bs-themes/bootstrap_default_themed.min.css"] = res
	m["htmlstyles/asciidoctor-bs-themes/bootstrap_default.css"] = res
	m["htmlstyles/asciidoctor-bs-themes/bootstrap_default.min.css"] = res
	m["htmlstyles/asciidoctor-bs-themes/bootstrap_slate.css"] = res
	m["htmlstyles/asciidoctor-bs-themes/bootstrap_slate.min.css"] = res
	m["htmlstyles/asciidoctor-skins/asciidoctor.css"] = res
	m["htmlstyles/asciidoctor-skins/boot-cerulean.css"] = res
	m["htmlstyles/asciidoctor-skins/boot-cosmo.css"] = res
	m["htmlstyles/asciidoctor-skins/boot-cyborg.css"] = res
	m["htmlstyles/asciidoctor-skins/boot-darkly.css"] = res
	m["htmlstyles/asciidoctor-skins/boot-flatly.css"] = res
	m["htmlstyles/asciidoctor-skins/boot-journal.css"] = res
	m["htmlstyles/asciidoctor-skins/boot-lumen.css"] = res
	m["htmlstyles/asciidoctor-skins/boot-paper.css"] = res
	m["htmlstyles/asciidoctor-skins/boot-readable.css"] = res
	m["htmlstyles/asciidoctor-skins/boot-sandstone.css"] = res
	m["htmlstyles/asciidoctor-skins/boot-slate.css"] = res
	m["htmlstyles/asciidoctor-skins/boot-spacelab.css"] = res
	m["htmlstyles/asciidoctor-skins/boot-superhero.css"] = res
	m["htmlstyles/asciidoctor-skins/boot-yeti.css"] = res
	m["htmlstyles/asciidoctor-skins/clean.css"] = res
	m["htmlstyles/asciidoctor-skins/dark.css"] = res
	m["htmlstyles/asciidoctor-skins/fedora.css"] = res
	m["htmlstyles/asciidoctor-skins/gazette.css"] = res
	m["htmlstyles/asciidoctor-skins/italian-pop.css"] = res
	m["htmlstyles/asciidoctor-skins/material-amber.css"] = res
	m["htmlstyles/asciidoctor-skins/material-blue.css"] = res
	m["htmlstyles/asciidoctor-skins/material-brown.css"] = res
	m["htmlstyles/asciidoctor-skins/material-green.css"] = res
	m["htmlstyles/asciidoctor-skins/material-grey.css"] = res
	m["htmlstyles/asciidoctor-skins/material-orange.css"] = res
	m["htmlstyles/asciidoctor-skins/material-pink.css"] = res
	m["htmlstyles/asciidoctor-skins/material-purple.css"] = res
	m["htmlstyles/asciidoctor-skins/material-red.css"] = res
	m["htmlstyles/asciidoctor-skins/material-teal.css"] = res
	m["htmlstyles/asciidoctor-skins/medium.css"] = res
	m["htmlstyles/asciidoctor-skins/monospace.css"] = res
	m["htmlstyles/asciidoctor-skins/notebook.css"] = res
	m["htmlstyles/asciidoctor-skins/plain.css"] = res
	m["htmlstyles/asciidoctor-skins/template.css"] = res
	m["htmlstyles/asciidoctor-skins/tufte.css"] = res
	m["htmlstyles/asciidoctor-skins/ubuntu.css"] = res
	m["htmlstyles/stylesheet-factory/asciidoctor.css"] = res
	m["htmlstyles/stylesheet-factory/colony.css"] = res
	m["htmlstyles/stylesheet-factory/foundation-lime.css"] = res
	m["htmlstyles/stylesheet-factory/foundation-potion.css"] = res
	m["htmlstyles/stylesheet-factory/foundation.css"] = res
	m["htmlstyles/stylesheet-factory/github.css"] = res
	m["htmlstyles/stylesheet-factory/golo.css"] = res
	m["htmlstyles/stylesheet-factory/iconic.css"] = res
	m["htmlstyles/stylesheet-factory/maker.css"] = res
	m["htmlstyles/stylesheet-factory/readthedocs.css"] = res
	m["htmlstyles/stylesheet-factory/riak.css"] = res
	m["htmlstyles/stylesheet-factory/rocket-panda.css"] = res
	m["htmlstyles/stylesheet-factory/rubygems.css"] = res
	m["htmlstyles/stylesheets-bulma/asciidoctor-embedded.css"] = res
	m["htmlstyles/stylesheets-bulma/asciidoctor-embedded.css.map"] = res
	m["htmlstyles/stylesheets-bulma/asciidoctor-embedded.min.css"] = res
	m["htmlstyles/stylesheets-bulma/asciidoctor.css"] = res
	m["htmlstyles/stylesheets-bulma/asciidoctor.css.map"] = res
	m["htmlstyles/stylesheets-bulma/asciidoctor.min.css"] = res
	m["htmlstyles/stylesheets-bulma/colony.css"] = res
	m["htmlstyles/stylesheets-bulma/colony.css.map"] = res
	m["htmlstyles/stylesheets-bulma/colony.min.css"] = res
	m["htmlstyles/stylesheets-bulma/darkly.css"] = res
	m["htmlstyles/stylesheets-bulma/darkly.css.map"] = res
	m["htmlstyles/stylesheets-bulma/darkly.min.css"] = res
	m["htmlstyles/stylesheets-bulma/foundation-lime.css"] = res
	m["htmlstyles/stylesheets-bulma/foundation-lime.css.map"] = res
	m["htmlstyles/stylesheets-bulma/foundation-lime.min.css"] = res
	m["htmlstyles/stylesheets-bulma/foundation-potion.css"] = res
	m["htmlstyles/stylesheets-bulma/foundation-potion.css.map"] = res
	m["htmlstyles/stylesheets-bulma/foundation-potion.min.css"] = res
	m["htmlstyles/stylesheets-bulma/foundation.css"] = res
	m["htmlstyles/stylesheets-bulma/foundation.css.map"] = res
	m["htmlstyles/stylesheets-bulma/foundation.min.css"] = res
	m["htmlstyles/stylesheets-bulma/github.css"] = res
	m["htmlstyles/stylesheets-bulma/github.css.map"] = res
	m["htmlstyles/stylesheets-bulma/github.min.css"] = res
	m["htmlstyles/stylesheets-bulma/golo.css"] = res
	m["htmlstyles/stylesheets-bulma/golo.css.map"] = res
	m["htmlstyles/stylesheets-bulma/golo.min.css"] = res
	m["htmlstyles/stylesheets-bulma/iconic.css"] = res
	m["htmlstyles/stylesheets-bulma/iconic.css.map"] = res
	m["htmlstyles/stylesheets-bulma/iconic.min.css"] = res
	m["htmlstyles/stylesheets-bulma/maker.css"] = res
	m["htmlstyles/stylesheets-bulma/maker.css.map"] = res
	m["htmlstyles/stylesheets-bulma/maker.min.css"] = res
	m["htmlstyles/stylesheets-bulma/readthedocs.css"] = res
	m["htmlstyles/stylesheets-bulma/readthedocs.css.map"] = res
	m["htmlstyles/stylesheets-bulma/readthedocs.min.css"] = res
	m["htmlstyles/stylesheets-bulma/riak.css"] = res
	m["htmlstyles/stylesheets-bulma/riak.css.map"] = res
	m["htmlstyles/stylesheets-bulma/riak.min.css"] = res
	m["htmlstyles/stylesheets-bulma/rocket-panda.css"] = res
	m["htmlstyles/stylesheets-bulma/rocket-panda.css.map"] = res
	m["htmlstyles/stylesheets-bulma/rocket-panda.min.css"] = res
	m["htmlstyles/stylesheets-bulma/rubygems.css"] = res
	m["htmlstyles/stylesheets-bulma/rubygems.css.map"] = res
	m["htmlstyles/stylesheets-bulma/rubygems.min.css"] = res

	return m
}

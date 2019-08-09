package cmd

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	assert := assert.New(t)
	// use in memory file system
	fs := afero.NewMemMapFs()
	// initial a project
	initProject(fs)
	// check for folders
	for _, v := range InitPaths() {
		exists, err := afero.DirExists(fs, v)
		assert.NoError(err)
		assert.True(exists)
	}
	// check files
	files := []string{
		filepath.Join(master),
		filepath.Join("references.bib"),
		filepath.Join("socrates.toml"),
		filepath.Join("front_matter", "abstract.adoc"),
		filepath.Join("front_matter", "dedication.adoc"),
		filepath.Join("front_matter", "preface.adoc"),
		filepath.Join("back_matter", "appendix_01.adoc"),
		filepath.Join("back_matter", "bibliography.adoc"),
		filepath.Join("back_matter", "colophon.adoc"),
		filepath.Join("back_matter", "glossary.adoc"),
		filepath.Join("back_matter", "index.adoc"),
		filepath.Join("parts", "part_01", "chapters", "chapter_01", "chapter_01.adoc"),
		filepath.Join("resources", "pdfstyles", "default-theme.yml"),
		filepath.Join("parts", "part_01", "part_01.adoc"),
	}

	for _, v := range files {

		exists, err := afero.Exists(fs, v)
		if !exists {
			fmt.Printf("checking: %s\n", v)
		}
		assert.NoError(err)
		assert.True(exists)

	}
}

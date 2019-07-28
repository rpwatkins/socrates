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
		filepath.Join("src", "master.adoc"),
		filepath.Join("src", "references.bib"),
		filepath.Join("src", "front_matter", "abstract.adoc"),
		filepath.Join("src", "front_matter", "dedication.adoc"),
		filepath.Join("src", "front_matter", "preface.adoc"),
		filepath.Join("src", "back_matter", "01_appendix.adoc"),
		filepath.Join("src", "back_matter", "bibliography.adoc"),
		filepath.Join("src", "back_matter", "colophon.adoc"),
		filepath.Join("src", "back_matter", "glossary.adoc"),
		filepath.Join("src", "back_matter", "index.adoc"),
		filepath.Join("src", "parts", "01_part", "01_chapter.adoc"),
		filepath.Join("src", "resources", "pdfstyles", "default-theme.yml"),
		filepath.Join("src", "parts", "01_part", "01_part.adoc"),
	}

	for _, v := range files {
		fmt.Printf("%s\n", v)
		exists, err := afero.Exists(fs, v)
		assert.NoError(err)
		assert.True(exists)
	}

}

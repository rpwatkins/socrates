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
	files := InitFileMap()
	for k, v := range files {
		f := k
		// check name
		if k[len(k)-5:] == "plush" {
			// run through plush with number = 1
			extension := filepath.Ext(k)
			f = k[0 : len(k)-len(extension)]
		}

		// rename files
		if f == "chapter.adoc" {
			f = "chapter_01.adoc"
		}
		if f == "appendix.adoc" {
			f = "appendix_01.adoc"
		}
		if f == "part.adoc" {
			f = "part_01.adoc"
		}

		exists, err := afero.Exists(fs, filepath.Join(v, f))

		assert.NoError(err)
		assert.True(exists)

		if !exists {
			fmt.Printf("%s not found", filepath.Join(v, f))
		}
	}
}

func Test_InitBare(t *testing.T) {

	assert := assert.New(t)
	// use in memory file system
	fs := afero.NewMemMapFs()

	err := writeBare(fs)
	assert.NoError(err)

	exists, err := afero.Exists(fs, "master.adoc")
	assert.NoError(err)
	assert.True(exists)

}

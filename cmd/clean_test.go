package cmd

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestClean(t *testing.T) {
	assert := assert.New(t)
	// use in memory file system
	fs := afero.NewMemMapFs()
	// create some dummy folders and files in build folder
	folders := []string{
		"build",
		filepath.Join("build", "html"),
		filepath.Join("build", "pdf"),
	}
	files := []string{
		filepath.Join("build", "html", "test.adoc"),
		filepath.Join("build", "pdf", "test.adoc"),
	}
	for _, v := range folders {
		if err := fs.MkdirAll(v, 0755); err != nil {
			fmt.Print(err)
		}
	}
	for _, v := range files {
		if err := afero.WriteFile(fs, v, []byte(""), 0644); err != nil {
			fmt.Print(err)
		}
	}
	clean(fs)
	for _, v := range folders {
		exists, _ := afero.DirExists(fs, v)
		assert.True(exists)
	}
	for _, v := range files {
		exists, _ := afero.Exists(fs, v)
		assert.True(exists)
	}

}

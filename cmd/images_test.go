package cmd

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestImages(t *testing.T) {
	assert := assert.New(t)
	// use in memory file system
	fs := afero.NewMemMapFs()
	// initial a project
	initProject(fs)

	// create some image files
	files := []string{
		filepath.Join("images", "image.jpg"),
		filepath.Join("images", "image2.png"),
		filepath.Join("images", "image3.jpg"),
	}
	// create the files
	for _, v := range files {
		if err := afero.WriteFile(fs, v, []byte(""), 0644); err != nil {
			fmt.Print(err)
		}
	}
	if err := fs.MkdirAll(filepath.Join("build", "html"), 0755); err != nil {
		fmt.Print(err)
	}
	if err := CopyFolder("images", filepath.Join("build", "html", "images"), fs); err != nil {
		fmt.Print(err)
	}

	for _, v := range files {
		exists, err := afero.Exists(fs, filepath.Join("build", "html", v))
		if err != nil {
			fmt.Print(err)
		}
		assert.True(exists)
	}
}

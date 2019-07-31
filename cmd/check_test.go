package cmd

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {

	assert := assert.New(t)
	// use in memory file system
	fs := afero.NewMemMapFs()
	// initial a project
	initProject(fs)

	c1 := "include::includes/include1.adoc"
	// create some nested include files
	if err := afero.WriteFile(fs, filepath.Join("parts", "part_01", "chapters", "chapter_02.adoc"), []byte(c1), 0644); err != nil {
		fmt.Print(err)
	}
	// ctreate included file
	if err := afero.WriteFile(fs, filepath.Join("parts", "part_01", "chapters", "includes", "include1.adoc"), []byte(""), 0644); err != nil {
		fmt.Print(err)
	}
	good, err := check(fs)
	assert.NoError(err)
	assert.True(good)
}

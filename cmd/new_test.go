package cmd

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	assert := assert.New(t)
	// use in memory file system
	fs := afero.NewMemMapFs()
	// initial a project
	initProject(fs)
	path := filepath.Join("parts", "part_01", "chapters", "chapter_02.adoc")
	newDoc(fs, []string{path})
	exists, _ := afero.Exists(fs, path)
	assert.True(exists)

}

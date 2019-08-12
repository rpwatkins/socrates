package cmd

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_DefaultBuild(t *testing.T) {
	// this tests the default init. The default as a total of twelve includes (with two in chapter_01.adoc)
	assert := assert.New(t)
	// use in memory file system
	fs := afero.NewMemMapFs()
	// initial a project
	initProject(fs)
	// build html
	buildHTML(fs)
	// check for build file
	exists, err := afero.Exists(fs, filepath.Join("build", "html", "file.html"))
	assert.True(exists)
	assert.NoError(err)

}

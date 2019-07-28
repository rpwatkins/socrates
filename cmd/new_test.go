package cmd

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/assert"
)

func Test_NewPart(t *testing.T) {
	assert := assert.New(t)
	// create in mem file system
	fs := afero.NewMemMapFs()
	// init a new project in file system
	initProject(fs)

	newPart(fs, []string{
		"test_part",
	})

	exists, err := afero.Exists(fs, filepath.Join("src", "parts", "test_part", "test_part.adoc"))
	assert.NoError(err)
	assert.True(exists)

	// check title
	newPart, err := afero.ReadFile(fs, filepath.Join("src", "parts", "test_part", "test_part.adoc"))
	assert.NoError(err)
	assert.True(strings.Contains(string(newPart), "test_part"))
	// check master.adoc
	master, err := afero.ReadFile(fs, filepath.Join("src", "master.adoc"))
	assert.NoError(err)
	assert.True(strings.Contains(string(master), "test_part"))

}

func Test_NewChapter(t *testing.T) {
	assert := assert.New(t)
	// create in mem file system
	fs := afero.NewMemMapFs()
	// init a new project in file system
	initProject(fs)

	newChapter(fs, []string{
		"test_chapter",
		"01_part",
	})

	exists, err := afero.Exists(fs, filepath.Join("src", "parts", "01_part", "test_chapter.adoc"))
	assert.NoError(err)
	assert.True(exists)

	// check title
	newChap, err := afero.ReadFile(fs, filepath.Join("src", "parts", "01_part", "test_chapter.adoc"))
	assert.NoError(err)
	assert.True(strings.Contains(string(newChap), "test_chapter"))
	// check master.adoc
	master, err := afero.ReadFile(fs, filepath.Join("src", "master.adoc"))
	assert.NoError(err)
	assert.True(strings.Contains(string(master), "include::parts/01_part/test_chapter"))

}

func Test_NewAppendix(t *testing.T) {
	assert := assert.New(t)
	// create in mem file system
	fs := afero.NewMemMapFs()
	// init a new project in file system
	initProject(fs)

	newAppendix(fs, []string{
		"test_appendix",
	})

	exists, err := afero.Exists(fs, filepath.Join("src", "back_matter", "test_appendix.adoc"))
	assert.NoError(err)
	assert.True(exists)

	// check title
	newPart, err := afero.ReadFile(fs, filepath.Join("src", "back_matter", "test_appendix.adoc"))
	assert.NoError(err)
	assert.True(strings.Contains(string(newPart), "test_appendix"))
	// check master.adoc
	master, err := afero.ReadFile(fs, filepath.Join("src", "master.adoc"))
	assert.NoError(err)
	assert.True(strings.Contains(string(master), "include::back_matter/test_appendix"))

}

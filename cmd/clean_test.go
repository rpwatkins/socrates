package cmd

import (
	"bytes"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestClean(t *testing.T) {
	assert := assert.New(t)
	// use in memory file system
	fs := afero.NewMemMapFs()
	// initial a project
	initProject(fs)

// save some files in the build folder
	file := []byte(`some test for a file`)

	for _, v := range []string{
		"html",
		"pdf",
		"fopub",
	} {
		if err := afero.WriteReader(fs, filepath.Join("src", "build", v, "output.txt"), bytes.NewReader(file)); err != nil {
			log.Error(err)
		}
	}

	clean(fs)

	count, err := afero.ReadDir(fs, filepath.Join("src", "build"))
	if err != nil {
		log.Error(err)
	}
	assert.Equal(0, len(count))

}

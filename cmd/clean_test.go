package cmd

import (
	"bytes"
	"fmt"
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

	count, err := afero.ReadDir(fs, "./src")
	if err != nil {
		log.Error(err)
	}
	assert.Equal(8, len(count))

	fmt.Print(len(count))

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

	// count2, err := afero.ReadDir(fs, "./")
	// if err != nil {
	// 	log.Error(err)
	// }
	// log.Infof("%d dirs", count2)
	// assert.Equal(8, len(count2))

}

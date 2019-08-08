package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestCheck_runValidation(t *testing.T) {

	// this tests the default init. The default as a total of twelve includes (with two in chapter_01.adoc)
	assert := assert.New(t)
	// use in memory file system
	fs := afero.NewMemMapFs()
	// initial a project
	initProject(fs)

	m := runValidation(fs)
	assert.Equal(0, len(m))

	// now we delete everything except master and test master.adoc
	for _, v := range []string{
		"front_matter",
		"back_matter",
		"parts",
		"images",
		"asets",
		"resoucres",
	} {
		if err := fs.RemoveAll(v); err != nil {
			fmt.Print(err)
		}
	}

	m2 := runValidation(fs)
	assert.Equal(9, len(m2))

}

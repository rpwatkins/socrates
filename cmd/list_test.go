package cmd

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {

	assert := assert.New(t)
	// use in memory file system
	fs := afero.NewMemMapFs()
	// initial a project
	initProject(fs)

	inc := "include::includes/include1.adoc"
	// open chapter_01.adoc and add include
	file, err := afero.ReadFile(fs, filepath.Join("parts", "part_01", "chapters", "chapter_01.adoc"))
	if err != nil {
		fmt.Print(err)
	}
	contents := string(file)
	newContents := fmt.Sprintf("%s\n\n%s", contents, inc)

	// ctreate included file
	if err := afero.WriteFile(fs, filepath.Join("parts", "part_01", "chapters", "includes", "include1.adoc"), []byte(newContents), 0644); err != nil {
		fmt.Print(err)
	}

	missing := runValidation(fs)
	for _, m := range missing {
		fmt.Printf("%s\n", m)
	}
	fmt.Printf("missing count = %d", len(missing))
	assert.NoError(err)
	assert.True(len(missing) == 0)
}

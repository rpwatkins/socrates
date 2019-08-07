package cmd

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestDisplay(t *testing.T) {
	fs := afero.NewMemMapFs()
	initProject(fs)
	// get a list of all the includes, images, etc.
	check(fs)
}

func TestCheck_InitPasses(t *testing.T) {

	// this tests the default init. The default as a total of twelve includes (with two in chapter_01.adoc)
	assert := assert.New(t)
	// use in memory file system
	fs := afero.NewMemMapFs()
	// initial a project
	initProject(fs)
	// get a list of all the includes, images, etc.
	incs := listMaster(fs, "master.adoc")
	// none missing
	f, m := flattenAndSortByMissingIncludes(incs)
	// twelve found
	assert.Equal(12, len(f))
	assert.Equal(0, len(m))
	for _, v := range m {
		fmt.Printf("%s\n", v.Path)
	}

}

func TestCheck_InitFails(t *testing.T) {

	assert := assert.New(t)
	// use in memory file system
	fs := afero.NewMemMapFs()
	// initial a project
	initProject(fs)
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
	// delete refences.bib
	if err := fs.Remove("references.bib"); err != nil {
		fmt.Print(err)
	}

	incs := listMaster(fs, "master.adoc")
	f, m := flattenAndSortByMissingIncludes(incs)
	// twelve found
	assert.Equal(0, len(f))
	assert.Equal(10, len(m))
}

func Test_flattenAndSortByMissingIncludes(t *testing.T) {

	assert := assert.New(t)

	child1 := include{
		Path:  "child1",
		Found: true,
	}
	parent1 := include{
		Path:  "parent1",
		Found: true,
		Includes: []include{
			child1,
		},
	}
	child2 := include{
		Path:  "child2",
		Found: false,
	}
	parent2 := include{
		Path:  "parent2",
		Found: false,
		Includes: []include{
			child2,
		},
	}
	incs := []include{
		parent1,
		parent2,
	}
	f, m := flattenAndSortByMissingIncludes(incs)
	assert.Equal(2, len(f))
	assert.Equal(2, len(m))
	for _, v := range f {
		assert.True(v.Found)
	}
	for _, v := range m {
		assert.False(v.Found)
	}

}

func Test_includeContains(t *testing.T) {
	assert := assert.New(t)

	i := include{
		Path: "root",
	}
	incs := []include{
		i,
	}
	res := includesContains(incs, i.Path)
	assert.True(res)

	res2 := includesContains(incs, "not-root")
	assert.False(res2)

}

func Test_checkItem(t *testing.T) {
	assert := assert.New(t)
	fs := afero.NewMemMapFs()
	// initial a project
	initProject(fs)

	i3 := checkItem(fs, "front_matter/preface.adoc")
	assert.True(i3.Found)
	i4 := checkItem(fs, "not-present.adoc")
	assert.False(i4.Found)
}

func Test_checkURL(t *testing.T) {
	assert := assert.New(t)

	url1 := "http://www.google.com"
	inc := checkURL(url1)
	assert.True(inc.Found)
	url2 := "http://drf.dwwdasffdf.com"
	inc2 := checkURL(url2)
	assert.False(inc2.Found)
}

func Test_ImageBlockMacro(t *testing.T) {
	assert := assert.New(t)
	fs := afero.NewMemMapFs()
	// initial a project

	c := filepath.Join("parts", "part_01", "chapters", "chapter_01", "chapter_01.adoc")

	initProject(fs)
	file, err := afero.ReadFile(fs, c)
	if err != nil {
		fmt.Print(err)
	}
	img := fmt.Sprint("image::test.png[]")
	contents := string(file)
	new := fmt.Sprintf("%s\n\n%s", contents, img)
	// write it back to fs
	if err := afero.WriteFile(fs, c, []byte(new), 0644); err != nil {
		fmt.Print(err)
	}
	// then check for fail
	incs := listMaster(fs, "master.adoc")
	// none missing
	f, m := flattenAndSortByMissingIncludes(incs)
	// twelve found
	assert.Equal(12, len(f))
	assert.Equal(1, len(m))
	// create file
	if err := afero.WriteFile(fs, filepath.Join("images", "test.png"), []byte(""), 0644); err != nil {
		fmt.Print(err)
	}

	newIncs := listMaster(fs, "master.adoc")
	found, missing := flattenAndSortByMissingIncludes(newIncs)

	// check success
	assert.Equal(13, len(found))
	assert.Equal(0, len(missing))

}

func Test_URL(t *testing.T) {
	assert := assert.New(t)
	fs := afero.NewMemMapFs()
	// initial a project

	c := filepath.Join("parts", "part_01", "chapters", "chapter_01", "chapter_01.adoc")

	initProject(fs)
	file, err := afero.ReadFile(fs, c)
	if err != nil {
		fmt.Print(err)
	}
	url1 := "http://www.google.com"
	url2 := "http://derf.sdewssedsadfftre.com"

	contents := string(file)
	new := fmt.Sprintf("%s\n\n%s\n%s", contents, url1, url2)
	// write it back to fs
	if err := afero.WriteFile(fs, c, []byte(new), 0644); err != nil {
		fmt.Print(err)
	}
	// then check for fail
	incs := listMaster(fs, "master.adoc")
	// none missing
	f, m := flattenAndSortByMissingIncludes(incs)

	// check success
	assert.Equal(13, len(f))
	assert.Equal(1, len(m))

	for _, v := range f {
		fmt.Printf("%s\n", v.Path)
	}
	for _, v := range m {
		fmt.Printf("\n\n%s\n", v.Path)
	}

}

func Test_diagram(t *testing.T) {

	assert := assert.New(t)
	fs := afero.NewMemMapFs()
	// initial a project

	c := filepath.Join("parts", "part_01", "chapters", "chapter_01", "chapter_01.adoc")

	initProject(fs)
	file, err := afero.ReadFile(fs, c)
	if err != nil {
		fmt.Print(err)
	}

	diagram := "mermaid"
	// add mermaid diagram
	d := fmt.Sprintf("%s::parts/part_01/chapters/chapter_01/test.mmd[]", diagram)
	contents := string(file)
	new := fmt.Sprintf("%s\n\n%s", contents, d)
	// write it back to fs
	if err := afero.WriteFile(fs, c, []byte(new), 0644); err != nil {
		fmt.Print(err)
	}
	// then check for fail
	incs := listMaster(fs, "master.adoc")
	// none missing
	f, m := flattenAndSortByMissingIncludes(incs)
	// twelve found
	assert.Equal(12, len(f))
	assert.Equal(1, len(m))

	// create file
	if err := afero.WriteFile(fs, filepath.Join("parts", "part_01", "chapters", "chapter_01", "test.mmd"), []byte(""), 0644); err != nil {
		fmt.Print(err)
	}

	newIncs := listMaster(fs, "master.adoc")
	found, missing := flattenAndSortByMissingIncludes(newIncs)

	// check success
	assert.Equal(13, len(found))
	assert.Equal(0, len(missing))

}

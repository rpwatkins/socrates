package cmd

import "path/filepath"

func InitPaths() []string {
	return []string{
		"src",
		filepath.Join("src", "front_matter"),
		filepath.Join("src", "back_matter"),
		filepath.Join("src", "images"),
		filepath.Join("src", "assets"),
		filepath.Join("src", "resources"),
		filepath.Join("src", "resources", "pdfstyles"),
		filepath.Join("src", "chapters"),
		filepath.Join("src", "chapters", "01_chapter"),
	}
}

func InitFileMap() map[string]string {
	// create a map from a file name to a destination fold
	src := "src"
	back := filepath.Join(src, "back_matter")
	front := filepath.Join(src, "front_matter")
	m := make(map[string]string)
	m["appendix.adoc.plush"] = back
	m["bibliography.adoc"] = back
	m["colophon.adoc"] = back
	m["glossary.adoc"] = back
	m["index.adoc"] = back
	m["preface.adoc"] = front
	m["dedication.adoc"] = front
	m["abstract.adoc"] = front
	m["master.adoc"] = src
	m["references.bib"] = src
	m["chapter.adoc.plush"] = filepath.Join(src, "chapters", "01_chapter")
	m["default-theme.yml"] = filepath.Join(src, "resources", "pdfstyles")
	m["part.adoc.plush"] = filepath.Join(src, "chapters")
	return m
}

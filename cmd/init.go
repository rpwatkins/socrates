package cmd

import (
  "os"
  "path"
  "path/filepath"

  "github.com/spf13/afero"
  "github.com/spf13/cobra"

  log "github.com/sirupsen/logrus"
)

var initCmd = &cobra.Command{
  Use:   "init",
  Short: "init creates a new manuscript.",
  Long:  `The init command creates a new manuscript.`,
  Run: func(cmd *cobra.Command, args []string) {
    initProject()
  },
}

func init() {
  rootCmd.AddCommand(initCmd)
}

func initProject() {

  // default filesystem
  var fs = afero.NewOsFs()

  // get the current working directory
  cwd, err := os.Getwd()
  if err != nil {
    log.Error(err.Error())
    os.Exit(1)
  }
  // confirm current working directory is empty
  empty, err := afero.IsEmpty(fs, cwd)
  if err != nil {
    log.Error(err.Error())
    os.Exit(1)
  }
  if !empty {
    log.Error("Folder not empty. A Socrates project can only be initialized in an empty folder.")
    os.Exit(1)
  }

  log.Info("Initializing Socrates project.")

  // create src folder/
  if err := fs.Mkdir(filepath.Join(cwd, "src"), 0755); err != nil {
    log.Error(err.Error())
  }
  src := path.Join(cwd, "src")

  // create src/chapters
  if err := fs.Mkdir(filepath.Join(src, "chapters"), 0755); err != nil {
    log.Error(err.Error())
  }

  // create /src/assets
  if err := fs.Mkdir(filepath.Join(src, "assets"), 0755); err != nil {
    log.Error(err.Error())
  }

  // create images folder
  if err := fs.Mkdir(filepath.Join(src, "images"), 0755); err != nil {
    log.Error(err.Error())
  }

  // create resources folder
  if err := fs.Mkdir(filepath.Join(src, "resources"), 0755); err != nil {
    log.Error(err.Error())
  }

  // create pdfstyles folder
  if err := fs.Mkdir(filepath.Join(src, "resources", "pdfstyles"), 0755); err != nil {
    log.Error(err.Error())
  }

  // create back_matter folder
  if err := fs.Mkdir(filepath.Join(src, "back_matter"), 0755); err != nil {
    log.Error(err.Error())
  }

  createMaster(fs)
  createChapter(fs)
  createBib(fs)
  createPDFStyles(fs)
  createBibliography(fs)

  log.Infof("Socrates project created at %s", cwd)

}

func createMaster(fs afero.Fs) {
  content := []byte(`:doctitle: Against Banning Student-Faculty Relationships
:subtitle:
:description:
:keywords:
:author: Rod Watkins
:authorinitials: RPW
:revnumber: 0.1.0
:email:  rodwatkins@outlook.com
:copyright: CC BY-NC-ND 4.0
:revdate: {docdate}
:doctype: article
:toc: left
:toclevels: 3
:stem:
:experimental:
:lang: en
:secnums:
:icons: font
:icon-set: fi
:source-highlighter: rouge
:bibliography-database: references.bib
:bibliography-style: chicago-author-date
:imagesdir: images
:imagesoutdir: {docdir}/images

include::chapters/chapter_01.adoc[]

include::back_matter/bibliography.adoc[]`)
  path := path.Join("src", "master.adoc")
  // create master.adoc and place content in file
  if err := afero.WriteFile(fs, path, content, 0644); err != nil {
    log.Error(err)
    os.Exit(1)
  }
}

func createChapter(fs afero.Fs) {
  content := []byte(`chapter one`)
  path := path.Join("src", "chapters", "chapter_01.adoc")

  // create master.adoc and place content in file
  if err := afero.WriteFile(fs, path, content, 0644); err != nil {
    log.Error(err)
    os.Exit(1)
  }
}

func createBib(fs afero.Fs) {
  content := []byte(``)
  path := path.Join("src", "references.bib")

  // create references.bib file
  if err := afero.WriteFile(fs, path, content, 0644); err != nil {
    log.Error(err)
    os.Exit(1)
  }
}

func createBibliography(fs afero.Fs) {
  content := []byte(`<<<
[bibliography]
== Bibliography

bibliography::[]`)
  path := path.Join("src", "back_matter", "bibliography.adoc")

  // create references.bib file
  if err := afero.WriteFile(fs, path, content, 0644); err != nil {
    log.Error(err)
    os.Exit(1)
  }
}

func createPDFStyles(fs afero.Fs) {

  content := []byte(`font:
catalog:
  # Noto Serif supports Latin, Latin-1 Supplement, Latin Extended-A, Greek, Cyrillic, Vietnamese & an assortment of symbols
  Noto Serif:
    normal: notoserif-regular-subset.ttf
    bold: notoserif-bold-subset.ttf
    italic: notoserif-italic-subset.ttf
    bold_italic: notoserif-bold_italic-subset.ttf
  # M+ 1mn supports ASCII and the circled numbers used for conums
  M+ 1mn:
    normal: mplus1mn-regular-ascii-conums.ttf
    bold: mplus1mn-bold-ascii.ttf
    italic: mplus1mn-italic-ascii.ttf
    bold_italic: mplus1mn-bold_italic-ascii.ttf
  # M+ 1p supports Latin, Latin-1 Supplement, Latin Extended, Greek, Cyrillic, Vietnamese, Japanese & an assortment of symbols
  # It also provides arrows for ->, <-, => and <= replacements in case these glyphs are missing from font
  M+ 1p Fallback:
    normal: mplus1p-regular-fallback.ttf
    bold: mplus1p-regular-fallback.ttf
    italic: mplus1p-regular-fallback.ttf
    bold_italic: mplus1p-regular-fallback.ttf
fallbacks:
  - M+ 1p Fallback
page:
  background_color: ffffff
  layout: portrait
  margin: [1in, 1in, 1in, 1in]
  size: Letter
base:
  align: justify
  # color as hex string (leading # is optional)
  font_color: 111111
  # color as RGB array
  #font_color: [51, 51, 51]
  # color as CMYK array (approximated)
  #font_color: [0, 0, 0, 0.92]
  #font_color: [0, 0, 0, 92%]
  font_family: Noto Serif
  # choose one of these font_size/line_height_length combinations
  #font_size: 14
  #line_height_length: 20
  #font_size: 11.25
  #line_height_length: 18
  #font_size: 11.2
  #line_height_length: 16
  font_size: 10.5
  line_height_length: 15
  # correct line height for Noto Serif metrics
  #line_height_length: 12
  #font_size: 11.25
  #line_height_length: 18
  line_height: $base_line_height_length / $base_font_size
  font_size_large: round($base_font_size * 1.25)
  font_size_small: round($base_font_size * 0.85)
  font_size_min: $base_font_size * 0.75
  font_style: normal
  border_color: eeeeee
  border_radius: 4
  border_width: 0.5
# FIXME vertical_rhythm is weird; we should think in terms of ems
#vertical_rhythm: $base_line_height_length * 2 / 3
# correct line height for Noto Serif metrics (comes with built-in line height)
vertical_rhythm: $base_line_height_length
horizontal_rhythm: $base_line_height_length
# QUESTION should vertical_spacing be block_spacing instead?
vertical_spacing: $vertical_rhythm
link:
  font_color: 428bca
# literal is currently used for inline monospaced in prose and table cells
literal:
  font_color: b12146
  font_family: M+ 1mn
menu_caret_content: " <font size=\"1.15em\"><color rgb=\"b12146\">\u203a</color></font> "
heading:
  #font_color: 181818
  font_color: $base_font_color
  font_family: M+ 1p Fallback
  font_style: bold
  # h1 is used for part titles (book doctype only)
  h1_font_size: floor($base_font_size * 2.6)
  # h2 is used for chapter titles (book doctype only)
  h2_font_size: floor($base_font_size * 2.15)
  h3_font_size: round($base_font_size * 1.7)
  h4_font_size: $base_font_size_large
  h5_font_size: $base_font_size
  h6_font_size: $base_font_size_small
  #line_height: 1.4
  # correct line height for Noto Serif metrics (comes with built-in line height)
  line_height: 1
  margin_top: $vertical_rhythm * 0.4
  margin_bottom: $vertical_rhythm * 0.9
title_page:
  align: right
  logo:
    top: 10%
  title:
    top: 55%
    font_size: $heading_h1_font_size
    font_color: 999999
    line_height: 0.9
  subtitle:
    font_size: $heading_h3_font_size
    font_style: bold_italic
    line_height: 1
  authors:
    margin_top: $base_font_size * 1.25
    font_size: $base_font_size_large
    font_color: 181818
  revision:
    margin_top: $base_font_size * 1.25
block:
  margin_top: 0
  margin_bottom: $vertical_rhythm
caption:
  align: left
  font_style: italic
  # FIXME perhaps set line_height instead of / in addition to margins?
  margin_inside: $vertical_rhythm / 3
  #margin_inside: $vertical_rhythm / 4
  margin_outside: 0
lead:
  font_size: $base_font_size_large
  line_height: 1.4
abstract:
  font_color: 5c6266
  font_size: $lead_font_size
  line_height: $lead_line_height
  font_style: italic
  first_line_font_style: bold
admonition:
  border_color: $base_border_color
  border_width: $base_border_width
  padding: [0, $horizontal_rhythm, 0, $horizontal_rhythm]
#  icon:
#    tip:
#      name: fa-lightbulb-o
#      stroke_color: 111111
#      size: 24
blockquote:
  font_color: $base_font_color
  font_size: $base_font_size_small
  border_color: $base_border_color
  border_width: 5
  padding: [$vertical_rhythm / 2, $horizontal_rhythm, $vertical_rhythm / -2, $horizontal_rhythm + $blockquote_border_width / 2]
  cite_font_size: $base_font_size_small
  cite_font_color: 999999
# code is used for source blocks (perhaps change to source or listing?)
code:
  font_color: $base_font_color
  font_family: $literal_font_family
  font_size: ceil($base_font_size)
  padding: $code_font_size
  line_height: 1.25
  background_color: f5f5f5
  border_color: cccccc
  border_radius: $base_border_radius
  border_width: 0.75
conum:
  font_family: M+ 1mn
  font_color: $literal_font_color
  font_size: $base_font_size
  line_height: 4 / 3
example:
  border_color: $base_border_color
  border_radius: $base_border_radius
  border_width: 0.75
  background_color: transparent
  # FIXME reenable margin bottom once margin collapsing is implemented
  padding: [$vertical_rhythm, $horizontal_rhythm, 0, $horizontal_rhythm]
image:
  align: left
prose:
  margin_top: 0
  margin_bottom: $vertical_rhythm
sidebar:
  border_color: $page_background_color
  border_radius: $base_border_radius
  border_width: $base_border_width
  background_color: eeeeee
  # FIXME reenable margin bottom once margin collapsing is implemented
  padding: [$vertical_rhythm, $vertical_rhythm * 1.25, 0, $vertical_rhythm * 1.25]
  title:
    align: center
    font_color: $heading_font_color
    font_family: $heading_font_family
    font_size: $heading_h4_font_size
    font_style: $heading_font_style
thematic_break:
  border_color: $base_border_color
  border_style: solid
  border_width: $base_border_width
  margin_top: $vertical_rhythm * 0.5
  margin_bottom: $vertical_rhythm * 1.5
description_list:
  term_font_style: italic
  term_spacing: $vertical_rhythm / 4
  description_indent: $horizontal_rhythm * 1.25
outline_list:
  indent: $horizontal_rhythm * 1.5
  # NOTE item_spacing applies to list items that do not have complex content
  item_spacing: $vertical_rhythm / 2
  #marker_font_color: 404040
table:
  background_color: $page_background_color
  #head_background_color: <hex value>
  #head_font_color: $base_font_color
  head_font_style: bold
  even_row_background_color: f9f9f9
  #odd_row_background_color: <hex value>
  foot_background_color: f0f0f0
  border_color: dddddd
  border_width: $base_border_width
  # HACK accounting for line-height
  cell_padding: [3, 3, 6, 3]
toc:
  dot_leader_color: dddddd
  #dot_leader_content: '. '
  indent: $horizontal_rhythm
  line_height: 1.4
# NOTE In addition to footer, header is also supported
footer:
  font_size: $base_font_size_small
  font_color: $base_font_color
  # NOTE if background_color is set, background and border will span width of page
  border_color: dddddd
  border_width: 0.25
  height: $base_line_height_length * 2.5
  line_height: 1
  padding: [$base_line_height_length / 2, 1, 0, 1]
  vertical_align: top
  #image_vertical_align: <alignment> or <number>
  # additional attributes for content:
  # * {page-count}
  # * {page-number}
  # * {document-title}
  # * {document-subtitle}
  # * {chapter-title}
  # * {section-title}
  # * {section-or-chapter-title}
  recto_content:
    #right: '{section-or-chapter-title} | {page-number}'
    #right: '{document-title} | {page-number}'
    right: '{page-number}'
    #center: '{page-number}'
  verso_content:
    #left: '{page-number} | {chapter-title}'
    left: '{page-number}'
    #center: '{page-number}'`)
  path := path.Join("src", "resources", "pdfstyles", "default-theme.yml")

  // create references.bib file
  if err := afero.WriteFile(fs, path, content, 0644); err != nil {
    log.Error(err)
    os.Exit(1)
  }

}

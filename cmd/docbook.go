package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

var docbookCmd = &cobra.Command{
	Use:   "docbook",
	Short: "docbook compiles a set of asciidoc files into a docbook5 file.",
	Long:  `The docbook command compiles a docbook5 file from a set of asciidoc files.`,
	Run: func(cmd *cobra.Command, args []string) {
		buildDocbook(afero.NewOsFs())
	},
}

func init() {
	rootCmd.AddCommand(docbookCmd)
}

func buildDocbook(fs afero.Fs) {

	// include check
	missing := runValidation(fs)
	if len(missing) > 0 {
		log.Error("build failed.")
		os.Exit(1)
	}
	// buildPDF creates a manuscript from a master.adoc file in the current directory
	// destination is the build folder under the cwd
	cwd, err := os.Getwd()
	if err != nil {
		log.Error("could not get current directory")
		os.Exit(1)
	}
	source := filepath.Join(cwd, "master.adoc")
	dest := filepath.Join(cwd, "build", "docbook")
	out := path.Base(cwd)

	command := AD
	args := []string{
		source,
		"--out-file=" + out + ".xml",
		"--require=asciidoctor-diagram",
		"--require=asciidoctor-bibliography",
		"--backend=docbook5",
		"--quiet",
		"-a imagesdir=images",
		"-a imagesoutdir=" + filepath.Join("build", "docbook", "images"),
		"--destination-dir=" + dest,
	}
	cmd := exec.Command(command, args...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err2 := cmd.Run()
	if err != nil {
		log.Error(err2)
		log.Error(outb.String())
		log.Error(errb.String())

		log.Errorf("%s DocBook could not be built", source)
		os.Exit(1)
	}
	log.Infof("%s DocBook build succeeded.", out)

}

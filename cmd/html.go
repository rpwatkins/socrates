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

var htmlCmd = &cobra.Command{
	Use:   "html",
	Short: "html compiles a set of asciidoc files into a pdf manuscript.",
	Long:  `The html command compiles an html page from a set of asciidoc files.`,
	Run: func(cmd *cobra.Command, args []string) {
		buildHTML(afero.NewOsFs())
	},
}

func init() {
	rootCmd.AddCommand(htmlCmd)
}

func buildHTML(fs afero.Fs) {

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
	source := filepath.Join("master.adoc")
	dest := filepath.Join("build", "html")
	out := path.Base(cwd)

	command := AD
	args := []string{
		source,
		"--out-file=" + out + ".html",
		"--require=asciidoctor-diagram",
		"--require=asciidoctor-bibliography",
		"--backend=html5",
		"--quiet",
		"-a max-width=55em",
		"-a imagesoutdir=" + filepath.Join("build", "html", "images"),
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

		log.Errorf("%s HTML could not be built", source)
		os.Exit(1)
	}

	log.Infof("%s.html build succeeded.", out)

}

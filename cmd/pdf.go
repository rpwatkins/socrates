package cmd

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var pdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "pdf compiles a set of asciidoc files into a pdf manuscript.",
	Long:  `The pdf command compiles a pdf manuscript from a set of asciidoc files.`,
	Run: func(cmd *cobra.Command, args []string) {
		buildPDF()
	},
}

func init() {
	rootCmd.AddCommand(pdfCmd)
}

func buildPDF() {

	// buildPDF creates a manuscript from a master.adoc file in the current directory
	// destination is the build folder under the cwd
	cwd, err := os.Getwd()
	if err != nil {
		log.Error("could not get curent directory")
		os.Exit(1)
	}
	source := filepath.Join("src", "master.adoc")
	dest := filepath.Join("build", "pdf")
	styles := filepath.Join("src", "resources", "pdfstyles")
	out := path.Base(cwd)

	command := "asciidoctor"
	args := []string{
		source,
		"--out-file=" + out,
		"--require=asciidoctor-diagram",
		"--require=asciidoctor-pdf",
		"--require=asciidoctor-bibliography",
		"--backend=docbook5",
		"--quiet",
		"-a pdf-stylesdir=" + styles,
		"-a pdf-style=default",
		"--destination-dir=" + dest,
	}
	cmd := exec.Command(command, args...)
	if err := cmd.Run(); err != nil {
		log.WithFields(log.Fields{
			"source": source,
		}).Errorf("%s PDF could not be built", source)
		os.Exit(1)
	}

}

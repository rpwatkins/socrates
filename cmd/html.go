package cmd

import (
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

	missing, err := check(fs)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if missing {
		log.Error("build failed. File(s) referenced by an include directives missing.")
		os.Exit(1)
	}

	// buildPDF creates a manuscript from a master.adoc file in the current directory
	// destination is the build folder under the cwd
	cwd, err := os.Getwd()
	if err != nil {
		log.Error("could not get current directory")
		os.Exit(1)
	}
	source := filepath.Join("src", "master.adoc")
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
		"-a data-uri",
		"--destination-dir=" + dest,
	}
	cmd := exec.Command(command, args...)
	if err := cmd.Run(); err != nil {
		log.Error(err)
		log.Errorf("%s HTML page could not be built", source)
		os.Exit(1)
	}
	log.Infof("%s.html build succeeded.", out)

}

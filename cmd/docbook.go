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
	if err := cmd.Run(); err != nil {
		log.WithFields(log.Fields{
			"source": source,
		}).Errorf("%s docbook could not be built", source)
		os.Exit(1)
	}

	if err := CopyFolder(filepath.Join("images"), filepath.Join("build", "docbook", "images"), fs); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Infof("%s.xml build succeeded.", out)

}

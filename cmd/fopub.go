package cmd

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var fopubCmd = &cobra.Command{
	Use:   "fopub",
	Short: "<fopub> compiles a set of asciidoc files into a pdf manuscript.",
	Long:  `The <fopub> command compiles a pdf manuscript from a set of asciidoc files using the fopub docbook toolchain.`,
	Run: func(cmd *cobra.Command, args []string) {
		buildFopub()
	},
}

func init() {
	rootCmd.AddCommand(fopubCmd)
}

func buildFopub() {

	// buildPDF creates a manuscript from a master.adoc file in the current directory
	// destination is the build folder under the cwd
	cwd, err := os.Getwd()
	if err != nil {
		log.Error("could not get current directory")
		os.Exit(1)
	}
	source := filepath.Join("src", "master.adoc")
	dest := filepath.Join("build", "fopub")
	out := path.Base(cwd)

	command := AD
	args := []string{
		source,
		"--out-file=" + out + ".xml",
		"--require=asciidoctor-diagram",
		"--require=asciidoctor-bibliography",
		"--backend=docbook5",
		"--quiet",
		"-a imagesdir=../../src/images",
		"--destination-dir=" + dest,
	}
	cmd := exec.Command(command, args...)
	if err := cmd.Run(); err != nil {
		log.WithFields(log.Fields{
			"source": source,
		}).Errorf("%s docbook file could not be built", source)
		os.Exit(1)
	}

	command2 := "fopub"
	args2 := []string{
		"build/fopub/" + out + ".xml",
	}
	cmd2 := exec.Command(command2, args2...)
	if err := cmd2.Run(); err != nil {
		log.WithFields(log.Fields{
			"source": source,
		}).Errorf("%s fopub pdf could not be built", source)
		os.Exit(1)
	}

}

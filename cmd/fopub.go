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

var fopubCmd = &cobra.Command{
	Use:   "fopub",
	Short: "fopub compiles a set of asciidoc files into a pdf manuscript.",
	Long:  `The fopub command compiles a pdf manuscript from a set of asciidoc files using the fopub docbook toolchain.`,
	Run: func(cmd *cobra.Command, args []string) {
		buildFopub(afero.NewOsFs())
	},
}

func init() {
	rootCmd.AddCommand(fopubCmd)
}

func buildFopub(fs afero.Fs) {

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
	dest := filepath.Join(cwd, "build", "fopub")
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
		log.Error(outb.String())

		log.Errorf("%s PDF could not be built", source)
		os.Exit(1)
	}

	command2 := "fopub"
	args2 := []string{
		"build/fopub/" + out + ".xml",
	}
	cmd2 := exec.Command(command2, args2...)
	var out2b, err2b bytes.Buffer
	cmd2.Stdout = &out2b
	cmd2.Stderr = &err2b
	err3 := cmd.Run()
	if err != nil {
		log.Error(err3)
		log.Error(out2b.String())
		log.Error(out2b.String())

		log.Errorf("%s PDF could not be built", source)
		os.Exit(1)
	}
	log.Infof("%s.pdf build succeeded.", out)

}

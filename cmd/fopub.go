package cmd

import (
	"fmt"
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

	missing, err := check(fs)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if len(missing) > 0 {
		log.Error("build failed. The following included file(s) could not be found.")
		for _, m := range missing {
			log.Warning(m)
		}
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
	if err := cmd.Run(); err != nil {
		log.Error(err)
		log.WithFields(log.Fields{
			"source": source,
		}).Errorf("%s docbook file could not be built", source)
		os.Exit(1)
	}

	if err := CopyFolder(filepath.Join("images"), filepath.Join("build", "docbook", "images"), fs); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	command2 := "fopub"
	args2 := []string{
		"build/fopub/" + out + ".xml",
	}
	cmd2 := exec.Command(command2, args2...)
	if err := cmd2.Run(); err != nil {
		fmt.Print(err)
		log.WithFields(log.Fields{
			"source": source,
		}).Errorf("%s fopub pdf could not be built", source)
		os.Exit(1)
	}
	log.Infof("%s.pdf build succeeded.", out)

}

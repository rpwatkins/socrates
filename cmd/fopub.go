package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	fopubCmd.PersistentFlags().StringP("output", "o", "output File Name (no extension)", "The name to be used for the output of the build commands: docbook, html, fopub, pdf")
	fopubCmd.PersistentFlags().Bool("timestamp", false, "Add the build timestamp to the output file name (default=false")

	if err := viper.BindPFlag("output", fopubCmd.PersistentFlags().Lookup("output")); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("timestamp", fopubCmd.PersistentFlags().Lookup("timestamp")); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func buildFopub(fs afero.Fs) {

	// include check
	missing := runValidation(fs)
	if len(missing) > 0 {
		log.Error("build failed. Some files are missing. Please run check.")
		os.Exit(1)
	}
	// buildPDF creates a manuscript from a master.adoc file in the current directory
	// destination is the build folder under the cwd

	source := master
	dest := filepath.Join("build", "fopub")
	out := viper.Get("output").(string)
	if viper.Get("timestamp").(bool) {
		out = fmt.Sprintf("%s--%s", out, time.Now().Format("2006-01-02-15-04-05"))
	}

	command := AD
	args := []string{
		source,
		"--out-file=" + out + ".xml",
		"--require=asciidoctor-diagram",
		"--require=asciidoctor-bibliography",
		"--backend=docbook5",
		"--quiet",
		"--destination-dir=" + dest,
	}
	cmd := exec.Command(command, args...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Error(err)
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
	err2 := cmd.Run()
	if err2 != nil {
		log.Error(err2)
		log.Error(out2b.String())
		log.Error(out2b.String())

		log.Errorf("%s PDF could not be built", source)
		os.Exit(1)
	}
	log.Infof("%s.pdf build succeeded.", out)

}

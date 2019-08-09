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

var pdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "pdf compiles a set of asciidoc files into a pdf manuscript.",
	Long:  `The pdf command compiles a pdf manuscript from a set of asciidoc files.`,
	Run: func(cmd *cobra.Command, args []string) {
		buildPDF(afero.NewOsFs())
	},
}

func init() {
	rootCmd.AddCommand(pdfCmd)

	pdfCmd.PersistentFlags().StringP("output", "o", "output File Name (no extension)", "The name to be used for the output of the build commands: docbook, html, fopub, pdf")
	pdfCmd.PersistentFlags().Bool("timestamp", false, "Add the build timestamp to the output file name (default=false")

	if err := viper.BindPFlag("output", pdfCmd.PersistentFlags().Lookup("output")); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("timestamp", pdfCmd.PersistentFlags().Lookup("timestamp")); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func buildPDF(fs afero.Fs) {

	// include check
	missing := runValidation(fs)
	if len(missing) > 0 {
		log.Error("build failed.")
		os.Exit(1)
	}
	// buildPDF creates a manuscript from a master.adoc file in the current directory
	// destination is the build folder under the cwd

	source := master
	dest := filepath.Join("build", "pdf")
	out := viper.Get("output").(string)
	if viper.Get("timestamp").(bool) {
		out = fmt.Sprintf("%s--%s", out, time.Now().Format("2006-01-02-15-04-05"))
	}
	command := AD
	args := []string{
		source,
		"--out-file=" + out + ".pdf",
		"--require=asciidoctor-diagram",
		"--require=asciidoctor-pdf",
		"--require=asciidoctor-bibliography",
		"--backend=pdf",
		"-a data-uri",
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
		log.Error(errb.String())

		log.Errorf("%s PDF could not be built", source)
		os.Exit(1)
	}
	log.Infof("%s.pdf build succeeded.", out)

}

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

	htmlCmd.PersistentFlags().StringP("output", "o", "output File Name (no extension)", "The name to be used for the output of the build commands: docbook, html, fopub, pdf")
	htmlCmd.PersistentFlags().Bool("timestamp", false, "Add the build timestamp to the output file name (default=false")

	if err := viper.BindPFlag("output", htmlCmd.PersistentFlags().Lookup("output")); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("timestamp", htmlCmd.PersistentFlags().Lookup("timestamp")); err != nil {
		log.Error(err)
		os.Exit(1)
	}
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
	source := master
	dest := filepath.Join("build", "html")
	out := viper.Get("output").(string)
	if viper.Get("timestamp").(bool) {
		out = fmt.Sprintf("%s--%s", out, time.Now().Format("2006-01-02-15-04-05"))
	}

	command := AD
	args := []string{
		source,
		"--out-file=" + out + ".html",
		"--require=asciidoctor-diagram",
		"--require=asciidoctor-bibliography",
		"--backend=html5",
		"--quiet",
		"-a max-width=55em",
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

		log.Errorf("%s HTML could not be built", source)
		os.Exit(1)
	}
	log.Infof("%s.html build succeeded.", out)
}

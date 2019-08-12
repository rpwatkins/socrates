package cmd

import (
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

	// file name options
	htmlCmd.PersistentFlags().StringP("output", "o", "output file name", "the name to be used for the output file")
	htmlCmd.PersistentFlags().Bool("timestamp", false, "add the build timestamp to the output file name")
	htmlCmd.PersistentFlags().Bool("skip", false, "skip validation")

	// read from config file
	if err := viper.BindPFlag("output", htmlCmd.PersistentFlags().Lookup("output")); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("timestamp", htmlCmd.PersistentFlags().Lookup("timestamp")); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("skip", htmlCmd.PersistentFlags().Lookup("skip")); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func buildHTML(fs afero.Fs) {

	// check file references
	if !viper.Get("skip").(bool) {
		missing := runValidation(fs)
		if len(missing) > 0 {
			log.Error("build failed. Some files are missing. Run check to view.")
			os.Exit(1)
		}
	}

	// get output file name from config
	out := viper.Get("output").(string)
	if viper.Get("timestamp").(bool) {
		out = fmt.Sprintf("%s--%s", out, time.Now().Format("2006-01-02-15-04-05"))
	}

	// configure build command
	command := AD
	args := []string{
		master,
		"--out-file=" + out + ".html",
		"--require=asciidoctor-diagram",
		"--require=asciidoctor-bibliography",
		"--backend=html5",
		"-a max-width=55em",
		"--destination-dir=" + filepath.Join("build", "html"),
	}
	// run asciidoctor
	result, err := exec.Command(command, args...).CombinedOutput()
	// display any mesages from asciidoctor (errors or warnings)
	r := string(result)
	if r != "" {
		fmt.Print(r)
	}
	if err != nil {
		log.Error(err)
		log.Errorf("%s HTML could not be built", master)
		os.Exit(1)
	}
}

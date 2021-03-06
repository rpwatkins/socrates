package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

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

	docbookCmd.Flags().StringP("output", "o", "output File Name (no extension)", "The name to be used for the output of the build commands: docbook, html, fopub, pdf")
	docbookCmd.Flags().Bool("timestamp", false, "Add the build timestamp to the output file name")
	docbookCmd.Flags().Bool("skip", false, "skip validation")

	if err := viper.BindPFlag("output", docbookCmd.Flags().Lookup("output")); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("timestamp", docbookCmd.Flags().Lookup("timestamp")); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("skip", docbookCmd.Flags().Lookup("skip")); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func buildDocbook(fs afero.Fs) {

	// include check
	if !viper.Get("skip").(bool) {
		missing := runValidation(fs)
		if len(missing) > 0 {
			log.Error("build failed")
			os.Exit(1)
		}
	}

	out := viper.Get("output").(string)
	if viper.Get("timestamp").(bool) {
		out = fmt.Sprintf("%s--%s", out, time.Now().Format("2006-01-02-15-04-05"))
	}

	command := AD
	args := []string{
		master,
		"--out-file=" + out + ".xml",
		"--require=asciidoctor-diagram",
		"--require=asciidoctor-bibliography",
		"--backend=docbook5",
		"--destination-dir=" + filepath.Join("build", "docbook"),
	}
	result, err := exec.Command(command, args...).CombinedOutput()
	r := string(result)
	if r != "" {
		fmt.Print(r)
	}
	if err != nil {
		log.Error(err)
		log.Errorf("%s DocBook could not be built", master)
		os.Exit(1)
	}

}

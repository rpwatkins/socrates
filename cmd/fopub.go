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
	dest := filepath.Join("build", "fopub")
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
		"--destination-dir=" + dest,
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

	command2 := "fopub"
	args2 := []string{
		"build/fopub/" + out + ".xml",
	}
	result2, err2 := exec.Command(command2, args2...).CombinedOutput()
	r2 := string(result2)

	if r2 != "" {
		fmt.Print(r2)
	}
	if err2 != nil {
		log.Error(err2)
		log.Errorf("%s PDF could not be built", master)
		os.Exit(1)
	}
}

package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// html
// fopub
//

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "socrates [command]",
	Short: "Generates manuscripts from asciidoc files.",
	Long:  `Socrates is a CLI app that generates html and pdf manuscripts from Asciidoctor files.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Error("command failed")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {

	// // get the current working directory
	// cwd, err := os.Getwd()
	// CheckErr(err)

	// // get project name based on folder name
	// projectName := path.Base(cwd)

	// // Search config in project directory with name reid (without extension).
	// viper.AddConfigPath(cwd)
	// viper.SetConfigType("toml")
	// viper.SetConfigName("reid")

	// // get environment variables
	// viper.AutomaticEnv()

	// // If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	Info("Using config file: " + viper.ConfigFileUsed())
	// 	viper.SetDefault("site", "../site")
	// 	viper.SetDefault("project", projectName)
	// } else {
	// 	Info("No config file found.")
	// }
}

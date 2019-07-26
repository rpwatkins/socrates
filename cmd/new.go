package cmd

import (
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [item type] [item name] [parent path]",
	Short: "new creates a new element of a manuscript.",
	Long:  `The new command creates elements of a manuscript. Types include: parts, chapters, sections, diagrams.`,
	Run: func(cmd *cobra.Command, args []string) {
		initProject()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// subcommands -- chapter, section, diagram, bibliography, appendix, index, glossary, colophon, part, dedication, appendix, acknowledgements, abstract, preface
}

/*
Copyright © 2022 Samuel Pearce
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "4bod-asm",
	Short: "Assembles 4BOD assembly language to binary instructions",
	Long: `This takes files in 4BOD assembly language and reduces it
to individual 4BOD instructions. Output can be given in various formats
such as binary, ascii-bitmaps, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: actually compile
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}

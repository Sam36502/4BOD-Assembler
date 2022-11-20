/*
Copyright © 2022 Samuel Pearce
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Sam36502/4BOD-Assembler/src/backend"
	"github.com/Sam36502/4BOD-Assembler/src/frontend"
	"github.com/spf13/cobra"
)

const (
	FLAG_OUTPUT = "output"
	FLAG_FRONT  = "front-end"
	FLAG_BACK   = "back-end"
	FLAG_LIST   = "list"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "4bod-asm",
	Short: "Assembles 4BOD assembly language to binary instructions",
	Long: `This takes files in 4BOD assembly language and reduces it
to individual 4BOD instructions. Output can be given in various formats
such as binary, ascii-bitmaps, etc. Use -l to see what formats are available`,
	Run: func(cmd *cobra.Command, args []string) {

		// If list flag set, display ends and quit
		if cmd.Flags().Changed(FLAG_LIST) {
			fmt.Println("Front-Ends:")
			for name, fe := range frontend.AvailableFrontends {
				fmt.Printf("  %5s - %s\n", name, fe.GetDescription())
			}

			fmt.Println("Back-Ends:")
			for name, be := range backend.AvailableBackends {
				fmt.Printf("  %5s - %s\n", name, be.GetDescription())
			}

			return
		}

		// Get args/flags
		if len(args) != 1 {
			fmt.Println("Exactly one argument (input filename) is required")
			return
		}
		infile := args[0]
		outfile, err := cmd.Flags().GetString(FLAG_OUTPUT)
		if err != nil {
			fmt.Println("Invalid output file specified")
			return
		}
		feName, err := cmd.Flags().GetString(FLAG_FRONT)
		front, exists := frontend.AvailableFrontends[feName]
		if err != nil || !exists {
			fmt.Println("Invalid front-end specified")
			return
		}
		beName, err := cmd.Flags().GetString(FLAG_BACK)
		back, exists := backend.AvailableBackends[beName]
		if err != nil || !exists {
			fmt.Println("Invalid back-end specified")
			return
		}

		// Assemble file
		program, err := front.ParseFile(infile)
		if err != nil {
			// TODO: Make front-ends return an array of errors
			//       and not return until the whole file is scanned
			fmt.Println("Failed to parse input file:")
			fmt.Println(err)
			return
		}
		err = back.GenerateFile(program, outfile)
		if err != nil {
			fmt.Println("Failed to generate output file:")
			fmt.Println(err)
			return
		}

		fmt.Printf("Output to file '%s'\n", outfile)
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
	rootCmd.Flags().StringP(FLAG_OUTPUT, "o", "a.out", "The file to output to (default 'a.out')")

	rootCmd.Flags().StringP(FLAG_FRONT, "f", frontend.FE_ASSEMBLY, "The front-end to use to parse the input file (default asm)")
	rootCmd.Flags().StringP(FLAG_BACK, "b", backend.BE_BINARY, "The back-end to use to generate the output file (default bin)")

	rootCmd.Flags().BoolP(FLAG_LIST, "l", false, "Lists out all available front and back ends")
}

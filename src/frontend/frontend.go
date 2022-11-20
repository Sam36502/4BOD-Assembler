/*
 *		Handles reading in assembly files
 */
package frontend

import "github.com/Sam36502/4BOD-Assembler/src/fbod"

type FrontendOptions struct {
	// None yet
}

type Frontend interface {
	SetOptions(opts FrontendOptions)
	GetOptions() FrontendOptions
	GetDescription() string
	ParseString(data string) (fbod.Program, error)
	ParseFile(filename string) (fbod.Program, error)
}

// Frontend Names
const (
	FE_ASSEMBLY = "asm"
)

var AvailableFrontends = map[string]Frontend{
	FE_ASSEMBLY: NewAssemblyFrontend(),
}

func GetDefaultOptions() FrontendOptions {
	return FrontendOptions{
		// None yet
	}
}

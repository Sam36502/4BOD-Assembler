/*
 *		Handles reading in assembly files
 */
package backend

import (
	"io/fs"

	"github.com/Sam36502/4BOD-Assembler/src/fbod"
)

type BackendOptions struct {
	FileMode fs.FileMode
}

type Backend interface {
	SetOptions(opts BackendOptions)
	GetOptions() BackendOptions
	GetDescription() string
	GenerateFile(prog fbod.Program, filename string) error
}

// Backend Names
const (
	BE_BINARY = "bin"
)

var AvailableBackends = map[string]Backend{
	BE_BINARY: NewBinaryBackend(),
}

func GetDefaultOptions() BackendOptions {
	return BackendOptions{
		FileMode: fs.FileMode(0650),
	}
}

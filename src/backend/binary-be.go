package backend

import (
	"io/ioutil"

	"github.com/Sam36502/4BOD-Assembler/src/fbod"
)

type BinaryBackend struct {
	opts BackendOptions
}

var _ Backend = (*BinaryBackend)(nil)

func NewBinaryBackend() *BinaryBackend {
	be := BinaryBackend{
		opts: GetDefaultOptions(),
	}
	return &be
}

func (be *BinaryBackend) SetOptions(opts BackendOptions) {
	be.opts = opts
}

func (be *BinaryBackend) GetOptions() BackendOptions {
	return be.opts
}

func (be *BinaryBackend) GetDescription() string {
	return "Binary file output. See README.md for binary file spec."
}

func (be *BinaryBackend) GenerateFile(prog fbod.Program, filename string) error {
	data := []byte{}
	for i := 0; i < fbod.FBOD_PROG_SIZE*fbod.FBOD_PAGE_SIZE; i += 2 {
		ins := prog[i/2]
		data[i] = byte(ins.Instruction)
		data[i+1] = byte((ins.Arg1 << 4) | ins.Arg2)
	}

	err := ioutil.WriteFile(filename, data, be.opts.FileMode)
	return err
}

/* Might be useful someday
func LoadBinary(filename string) (Program, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	prog := Program{}
	for i := 0; i < FBOD_PROG_SIZE*FBOD_PAGE_SIZE; i += 2 {
		upper := data[i]
		lower := data[i+1]
		prog[i/2] = Instruction{
			Instruction: upper % 16,
			Arg1:        (lower >> 4) % 16,
			Arg2:        lower % 16,
		}
	}

	return prog, nil
}
*/

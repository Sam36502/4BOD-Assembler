package fbod

import "io/ioutil"

func LoadBinary(filename string) (Program, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	prog := Program{}
	for i := 0; i < FBOD_PROG_SIZE*FBOD_PAGE_SIZE; i += 2 {
		upper := data[i]
		lower := data[i+1]
		prog[i] = Instruction{
			instruction: upper % 16,
			arg1:        (lower >> 4) % 16,
			arg2:        lower % 16,
		}
	}

	return prog, nil
}

func SaveBinary(filename string, prog Program) error {
}

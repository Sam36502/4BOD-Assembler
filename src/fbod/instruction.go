package fbod

type Instruction struct {
	Instruction byte
	Arg1        byte
	Arg2        byte
}

type Program []Instruction

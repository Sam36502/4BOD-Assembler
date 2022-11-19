package fbod

type Instruction struct {
	instruction byte
	arg1        byte
	arg2        byte
}

type Program []Instruction

/*
 *	Front-End implementation that reads 4BOD assembly (`.4sm`)
 */
package frontend

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/Sam36502/4BOD-Assembler/src/fbod"
)

type AssemblyFrontend struct {
	opts FrontendOptions
}

// Assembly language Regexes
const (
	ASM_REGEX_COMMENT    = ";.*$"
	ASM_REGEX_WHITESPACE = "[ \t]+"
)

// Assembly language symbols
const (
	ASM_STR_NEWLINE  = "\n"
	ASM_STR_SPECIAL  = "#"
	ASM_CHAR_SPACE   = ' '
	ASM_CHAR_TAB     = '\t'
	ASM_CHAR_COMMENT = ';'
)

// Special Commands
const (
	ASM_SPEC_VAR = "var"
	ASM_SPEC_LBL = "label"
)

// Assembly language op-codes
const (
	ASM_OP_NOP = 0x0
	ASM_OP_MVA = 0x1
	ASM_OP_MVM = 0x2
	ASM_OP_STA = 0x3
	ASM_OP_INA = 0x4
	ASM_OP_INC = 0x5
	ASM_OP_CLS = 0x6
	ASM_OP_SHL = 0x7
	ASM_OP_SHR = 0x8
	ASM_OP_RDP = 0x9
	ASM_OP_FLP = 0xA
	ASM_OP_FLG = 0xB
	ASM_OP_JMP = 0xC
	ASM_OP_CEQ = 0xD
	ASM_OP_CGT = 0xE
	ASM_OP_CLT = 0xF
)

// Assembly names
var ASM_NAMES = map[string]byte{
	"NOP": 0x0,
	"MVA": 0x1,
	"MVM": 0x2,
	"STA": 0x3,
	"INA": 0x4,
	"INC": 0x5,
	"CLS": 0x6,
	"SHL": 0x7,
	"SHR": 0x8,
	"RDP": 0x9,
	"FLP": 0xA,
	"FLG": 0xB,
	"JMP": 0xC,
	"CEQ": 0xD,
	"CGT": 0xE,
	"CLT": 0xF,
}

// Number Formats
var ASM_NUMFMT = map[int]string{
	2:  "0b", // Binary prefix
	8:  "0",  // Octal prefix
	10: "",   // Decimal prefix
	16: "0x", // Hexadecimal prefix
}

// Compile-time implementation check
var _ Frontend = (*AssemblyFrontend)(nil)

func NewAssemblyFrontend() *AssemblyFrontend {
	fe := AssemblyFrontend{
		opts: GetDefaultOptions(),
	}
	return &fe
}

func (fe *AssemblyFrontend) SetOptions(opts FrontendOptions) {
	fe.opts = opts
}

func (fe *AssemblyFrontend) GetOptions() FrontendOptions {
	return fe.opts
}

func (fe *AssemblyFrontend) GetDescription() string {
	return "The assembly front-end. See README.md for language spec."
}

func (fe *AssemblyFrontend) ParseString(data string) (fbod.Program, error) {
	prog := fbod.Program{}
	labels := map[string]byte{}
	nextlabel := 0
	vars := map[string]byte{}

	// Preprocess file removing unneccessary formatting
	commentRegex := regexp.MustCompile(ASM_REGEX_COMMENT)
	whitespaceRegex := regexp.MustCompile(ASM_REGEX_WHITESPACE)
	data = string(commentRegex.ReplaceAll([]byte(data), []byte{}))                  // Remove all comments
	data = string(whitespaceRegex.ReplaceAll([]byte(data), []byte{ASM_CHAR_SPACE})) // Reduce multiple whitespace to single space

	// File should be mostly uniform columns of assembly and arguments
	for i, line := range strings.Split(data, ASM_STR_NEWLINE) {

		// Ignore empty lines
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// Split line by space for individual fields
		fields := strings.Split(line, string(ASM_CHAR_SPACE))

		// Check for special commands
		if strings.HasPrefix(fields[0], ASM_STR_SPECIAL) {
			specCmd := strings.Trim(fields[0], ASM_STR_SPECIAL+string(ASM_CHAR_SPACE))
			switch specCmd {

			case ASM_SPEC_VAR:
				if len(fields) != 3 {
					return fbod.Program{}, formatSyntaxError(i, "`#var` special command requires two arguments", nil)
				}
				name := strings.ToLower(fields[1])
				addr, err := parseNumber(fields[2])
				if err != nil {
					return fbod.Program{}, formatSyntaxError(i, "`#var` special command requires numeric second argument", err)
				}
				vars[name] = byte(addr)

			case ASM_SPEC_LBL:
				if len(fields) != 2 {
					return fbod.Program{}, formatSyntaxError(i, "`#label` special command requires one argument", nil)
				}
				name := strings.ToLower(fields[1])
				if nextlabel >= 16 {
					return fbod.Program{}, formatSyntaxError(i, "Exceeded maximum number of labels (16)", nil)
				}
				labels[name] = byte(nextlabel)
				prog = append(prog, fbod.Instruction{
					Instruction: ASM_OP_FLG,
					Arg1:        byte(nextlabel),
				})
				nextlabel++

			default:
				return fbod.Program{}, formatSyntaxError(i, fmt.Sprintf("'#%s' is not a recognised special command", specCmd), nil)

			}
			continue
		}

		// Try to find Assembly Opcodes
		ins := fbod.Instruction{}
		opstr := strings.ToUpper(fields[0])
		found := false
		for opcode, bin := range ASM_NAMES {
			if opstr == opcode {
				found = true
				ins.Instruction = bin
				break
			}
		}
		if !found {
			return fbod.Program{}, formatSyntaxError(i, fmt.Sprintf("Invalid Opcode '%s'", opstr), nil)
		}

		// Handle arguments
		switch ins.Instruction {

		// One address argument
		case ASM_OP_MVA:
			fallthrough
		case ASM_OP_MVM:
			fallthrough
		case ASM_OP_CEQ:
			fallthrough
		case ASM_OP_CGT:
			fallthrough
		case ASM_OP_CLT:
			addr, err := parseNamedNum(fields[1], vars)
			if err != nil {
				return fbod.Program{}, formatSyntaxError(i, "No valid address argument found", err)
			}
			ins.Arg1 = addr

		// One numeric argument
		case ASM_OP_STA:
			num, err := parseNumber(fields[1])
			if err != nil {
				return fbod.Program{}, formatSyntaxError(i, "No valid address argument found", err)
			}
			ins.Arg1 = byte(num)

		// One Label argument
		case ASM_OP_FLG:
			fallthrough
		case ASM_OP_JMP:
			labelnum, err := parseNamedNum(fields[1], labels)
			if err != nil {
				return fbod.Program{}, formatSyntaxError(i, "Invalid label argument", err)
			}
			ins.Arg1 = labelnum

		// Two address arguments
		case ASM_OP_RDP:
			fallthrough
		case ASM_OP_FLP:
			addr, err := parseNamedNum(fields[1], vars)
			if err != nil {
				return fbod.Program{}, formatSyntaxError(i, "(Arg 1) Invalid address argument", err)
			}
			ins.Arg1 = addr
			addr, err = parseNamedNum(fields[2], vars)
			if err != nil {
				return fbod.Program{}, formatSyntaxError(i, "(Arg 2) Invalid address argument", err)
			}
			ins.Arg2 = addr

		}

	}

	return prog, nil
}

func (fe *AssemblyFrontend) ParseFile(filename string) (fbod.Program, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fbod.Program{}, err
	}

	return fe.ParseString(string(data))
}

func formatSyntaxError(lineIndex int, msg string, suppl error) error {
	supplStr := ""
	if suppl != nil {
		supplStr = fmt.Sprintf(": %v", suppl)
	}
	return fmt.Errorf("%0d: Syntax error: %s%s", lineIndex+1, msg, supplStr)
}

func parseNumber(str string) (int, error) {
	str = strings.TrimSpace(str)

	for base, prefix := range ASM_NUMFMT {
		if strings.HasPrefix(str, prefix) {
			if num, err := strconv.ParseUint(strings.TrimPrefix(str, prefix), base, 64); err == nil {
				if num > 16 {
					return -1, fmt.Errorf("'%s' is too large to be a 4-bit number", str)
				}
				return int(num), nil
			}
		}
	}

	return -1, fmt.Errorf("'%s' is not a valid number", str)
}

func parseNamedNum(str string, names map[string]byte) (byte, error) {
	str = strings.ToLower(str)

	// Check if it's a number
	if num, err := parseNumber(str); err == nil {
		return byte(num), nil
	}

	// Check for name
	for name, num := range names {
		if name == str {
			return num, nil
		}
	}
	return 0, fmt.Errorf("no number associated with name '%s'", str)
}

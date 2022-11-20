# 4BOD-Assembler
A basic assembler to make 4BOD programs easier to write

## Assembly Language (`.4sm`)
The most basic feature the assembler can provide is that you can use
names for the instructions instead of memorising its binary opcodes,
but I also wanted to provide a few helpful quality-of-life improvements
common in other assemblers as well including named variables and comments.

### Instructions & Opcodes
Each 4BOD instruction has been assigned a 3-letter opcode to make it easier
to remember what each one does. They are as follows:

| 4BOD   | Opcode | Description |
|--------|--------|-------------|
| `0000` | `NOP`  | Do nothing. |
| `0001` | `MVA`  | Sets the accumulator to the value stored at the address arg1. |
| `0010` | `MVM`  | Writes the accumulator to the address arg1.  |
| `0011` | `STA`  | Sets the accumulator to arg1.  |
| `0100` | `INA`  | Sets the accumulator to the states of the arrow keys in the order D|U|R|L. |
| `0101` | `INC`  | Increments the accumulator, overflowing if it goes above 15.  |
| `0110` | `CLS`  | Clears the screen. (setting all VRAM to 0) |
| `0111` | `SHL`  | Shifts the accumulator left. |
| `1000` | `SHR`  | Shifts the accumulator right.   |
| `1001` | `RDP`  | Sets the accumulator to the state of the pixel at the X and Y coordinates of the values at addresses arg1 and arg2, respectively |
| `1010` | `FLP`  | Flips the state of the pixel at the X and Y coordinates of the values at addresses arg1 and arg2, respectively.    |
| `1011` | `FLG`  | Creates a flag named arg1. |
| `1100` | `JMP`  | Jumps to the flag named after the value stored at the address arg1. (If 2 flags have the same ID the one later in the code is the one jumped to) |
| `1101` | `CEQ`  | Only perform next instruction if the value at address arg1 is equal to the accumulator. |
| `1110` | `CGT`  | Only perform next instruction if the value at address arg1 is greater than the accumulator. |
| `1111` | `CLT`  | Only perform next instruction if the value at address arg1 is less than the accumulator. |

(Instruction descriptions from the [Esolangs Page](https://esolangs.org/wiki/4BOD))

### Special Commands
Special commands are extra things that can be defined in the program
which make programming easier. They are as follows:

#### `#var`
This will essentially name a memory address so it can be used like a variable.
Every time it's called will change its name, so vars can be renamed later in the
source file.

Example:

    #var x 0x0          ; Assigns name 'x' to address 0x0
    #var y 0x1          ; Assigns name 'y' to address 0x1
        FLP     x   y   ; Flips pixel at (x/y)

#### `#label`
This is essentially the same as the `FLG` instruction which adds a jump label to the code,
except that it will let you use names for the labels instead of numbers.

Example:

    #label loop         ; Is identical to...
        FLG     0       ; this
        JMP     loop    ; `JMP` calls also recognise named labels

It's not recommended to mix `#label` & `FLG` in the same script as the label command
simply picks numbers in ascending order, so you could easily accidentally reassign a
label (such as in the above example; `loop` is essentially an alias for 0)


## Binary File Format (`.4bb`)
This way of storing 4BOD instructions is really just storing
the instructions in 2 Bytes; There are 4 bits wasted per instruction,
but for the sake of simplicity and given the limited size of
4BOD binaries it doesn't make sense to pack the data any smaller.

This format is simply a series of 2-Byte instructions as follows:

    Byte 1:
    0000      Top 4 bits are empty
        0000  Bottom 4 bits of first byte stores the instruction opcode

    Byte 2:
    0000      Top 4 bits are argument 1
        0000  Bottom 4 bits are argument 2

This pattern repeats every 2 Bytes for a total of 512 Bytes.
Whether empty instructions after the end of the program is stored
is left unspecified.

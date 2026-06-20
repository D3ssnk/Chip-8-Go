// Package cpu implements the CHIP-8 virtual machine's central processing unit.
// It handles instruction execution, memory management, registers, and the program counter.
package cpu

import (
	"errors"
	"io"
	"os"
)

// font contains the 5x8 bitmap data for hexadecimal characters (0-F).
// Each character is 5 bytes, stored as 80 bytes total.
// These are loaded into CPU memory at addresses 0x000-0x04F during initialization.
var font = [80]uint8{0xF0, 0x90, 0x90, 0x90, 0xF0, //ZERO
	0x20, 0x60, 0x20, 0x20, 0x70, //ONE
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //TWO
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //THREE
	0x90, 0x90, 0xF0, 0x10, 0x10, //FOUR
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //FIVE
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //SIX
	0xF0, 0x10, 0x20, 0x40, 0x40, //SEVEN
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //EIGHT
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //NINE
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80} //F

// CPU represents the CHIP-8 processor with its memory, registers, and I/O devices.
type CPU struct {
	// Memory is the 4096-byte RAM. Programs start at address 0x200.
	memory [4096]uint8
	// Registers are 16 general-purpose 8-bit registers (V0-VF).
	// VF is reserved as a flag register for carry/borrow operations.
	registers [16]uint8
	// I is the 16-bit index register used for memory operations.
	i uint16
	// PC is the program counter, pointing to the next instruction to execute.
	pc uint16
	// SP is the stack pointer, indexing into the call stack.
	sp uint8
	// Stack holds up to 16 return addresses for subroutine calls.
	stack [16]uint16
	// DelayTimer is an 8-bit timer that decrements at 60Hz when non-zero.
	delayTimer uint8
	// SoundTimer is an 8-bit timer that decrements at 60Hz; produces sound when non-zero.
	soundTimer uint8
	// Display is the 64x32 pixel screen (stored as 32 rows of 8-byte width).
	display [32][8]uint8
	// Keypad holds the state of the 16 hexadecimal input keys (0x0-0xF).
	keypad [16]bool
}

// loadROM reads a CHIP-8 ROM file and loads its contents into CPU memory starting at address 0x200.
// This is where CHIP-8 programs are expected to begin execution.
func (cpu *CPU) LoadROM(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return errors.New("Couldn't open file.")
	}

	defer file.Close()

	content, err := io.ReadAll(file)

	if err != nil {
		return errors.New("Couldn't read from file.")
	}

	if len(content) > 3584 {
		return errors.New("Rom file is too large")
	}

	copy(cpu.memory[0x200:], content[:])
	return nil
}

func (cpu *CPU) Fetch() (uint16, error) {
	if cpu.pc >= 0xFFF {
		return 0, errors.New("Program counter exceeds memory")
	}
	var instruction uint16
	instruction = uint16(cpu.memory[cpu.pc])
	cpu.pc++
	instruction <<= 8
	instruction |= uint16(cpu.memory[cpu.pc])
	cpu.pc++

	return instruction, nil
}

// NewCPU initializes a new CHIP-8 CPU with default values.
// The program counter starts at 0x200 (standard CHIP-8 program start).
// The font data is loaded into memory at addresses 0x000-0x04F.
func NewCPU() CPU {
	var cpu CPU
	// Load the built-in font data into the reserved memory area.
	copy(cpu.memory[:], font[:])
	// Set the program counter to the standard CHIP-8 program start address.
	cpu.pc = 0x200
	return cpu
}

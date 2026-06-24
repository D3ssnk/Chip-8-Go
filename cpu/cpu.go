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
	display [32][64]bool
	// Keypad holds the state of the 16 hexadecimal input keys (0x0-0xF).
	keypad [16]bool
	// WaitngForKey is a bool that flags if the program is waiting for a key press.
	waitingForKey bool
}

// LoadROM reads a CHIP-8 ROM file and loads its contents into CPU memory starting at address 0x200.
// This is where CHIP-8 programs are expected to begin execution.
// It returns an error if the file exceeds the 3584-byte limit of available memory.
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

// Fetch reads the next 16-bit instruction from memory by combining two consecutive bytes.
// It advances the Program Counter (PC) by 2 and returns an error if the PC exceeds 0xFFF.
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

// GetDisplay returns the current state of the 64x32 monochrome display matrix.
func (cpu *CPU) GetDisplay() [32][64]bool {
	return cpu.display
}

// SetKeypad updates the state of a specific key on the keypad (pressed or unpressed).
func (cpu *CPU) SetKeypad(index int, state bool) {
	cpu.keypad[index] = state
}

// SetWaitingForKeypad toggles the CPU flag indicating whether the system is halting for user input.
func (cpu *CPU) SetWaitingForKeypad(state bool) {
	cpu.waitingForKey = state
}

// GeWaitingForKeypad returns whether the CPU is currently halted waiting for user input.
func (cpu *CPU) GeWaitingForKeypad() bool {
	return cpu.waitingForKey
}

// DecrementDelayTimer reduces the delay timer by 1, simulating the standard 60Hz countdown.
func (cpu *CPU) DecrementDelayTimer() {
	cpu.delayTimer--
}

// GetDelayTimer returns the current value of the delay timer.
func (cpu *CPU) GetDelayTimer() uint8 {
	return cpu.delayTimer
}

// DecrementSoundTimer reduces the sound timer by 1, simulating the standard 60Hz countdown.
func (cpu *CPU) DecrementSoundTimer() {
	cpu.soundTimer--
}

// GetSoundTimer returns the current value of the sound timer.
func (cpu *CPU) GetSoundTimer() uint8 {
	return cpu.soundTimer
}

// ResetKeypad clears all active keypresses, setting the keypad array back to its default false state.
func (cpu *CPU) ResetKeypad() {
	cpu.keypad = [16]bool{}
}

// GetKeypad returns the full boolean array representing the current state of all 16 keys.
func (cpu *CPU) GetKeypad() [16]bool {
	return cpu.keypad
}

// Execute decodes and executes a 16-bit CHIP-8 instruction.
// Instructions are decoded by examining the first nibble to determine the opcode category,
// then further decoded using secondary nibbles as needed.
func (cpu *CPU) Execute(instruction uint16) error {
	firstNibble := instruction >> 12
	secondNibble := (instruction >> 8) & 0x0F
	thirdNibble := (instruction >> 4) & 0x0F
	finalNibble := instruction & 0x0F
	lastByte := instruction & 0xFF
	last12Bits := instruction & 0xFFF

	//log.Println(instruction)
	switch firstNibble {
	case 0x0:
		// 0x00E0: Clear display or 0x00EE: Return from subroutine
		switch lastByte {
		case 0xE0:
			// Clear the display
			cpu.clear()
			return nil
		case 0xEE:
			// Return from subroutine (pop PC from stack)
			return cpu.ret()
		}

	case 0x1:
		// 0x1nnn: Jump to address nnn
		return cpu.jump(last12Bits)

	case 0x2:
		// 0x2nnn: Call subroutine at nnn (push PC to stack)
		return cpu.call(last12Bits)

	case 0x3:
		// 0x3xkk: Skip next instruction if Vx == kk
		return cpu.skipIfEqual(secondNibble, lastByte)

	case 0x4:
		// 0x4xkk: Skip next instruction if Vx != kk
		return cpu.skipIfNotEqual(secondNibble, lastByte)

	case 0x5:
		// 0x5xy0: Skip next instruction if Vx == Vy
		return cpu.skipIfEqualReg(secondNibble, thirdNibble)

	case 0x6:
		// 0x6xkk: Set Vx = kk
		return cpu.setReg(secondNibble, lastByte)

	case 0x7:
		// 0x7xkk: Add kk to Vx
		return cpu.addVal(secondNibble, lastByte)

	case 0x8:
		// 0x8xy_: Arithmetic and logic operations
		switch finalNibble {
		case 0x0:
			// 0x8xy0: Set Vx = Vy
			return cpu.setReg(secondNibble, thirdNibble)
		case 0x1:
			// 0x8xy1: Set Vx = Vx OR Vy
			return cpu.orReg(secondNibble, thirdNibble)
		case 0x2:
			// 0x8xy2: Set Vx = Vx AND Vy
			return cpu.andReg(secondNibble, thirdNibble)
		case 0x3:
			// 0x8xy3: Set Vx = Vx XOR Vy
			return cpu.xorReg(secondNibble, thirdNibble)
		case 0x4:
			// 0x8xy4: Add Vy to Vx (set VF = carry)
			return cpu.addReg(secondNibble, thirdNibble)
		case 0x5:
			// 0x8xy5: Subtract Vy from Vx (set VF = NOT borrow)
			return cpu.subReg(secondNibble, thirdNibble)
		case 0x6:
			// 0x8xy6: Shift Vx right by 1 (set VF = LSB of Vx)
			return cpu.shiftRight(secondNibble)
		case 0x7:
			// 0x8xy7: Set Vx = Vy - Vx (set VF = NOT borrow)
			return cpu.subRegReverse(secondNibble, thirdNibble)
		case 0xE:
			// 0x8xyE: Shift Vx left by 1 (set VF = MSB of Vx)
			return cpu.shiftLeft(secondNibble)
		}

	case 0x9:
		// 0x9xy0: Skip next instruction if Vx != Vy
		return cpu.skipIfNotEqualReg(secondNibble, thirdNibble)

	case 0xA:
		// 0xAnnn: Set I = nnn
		return cpu.setI(last12Bits)

	case 0xB:
		// 0xBnnn: Jump to address nnn + V0
		return cpu.jumpV0(last12Bits)

	case 0xC:
		// 0xCxkk: Set Vx = random byte AND kk
		return cpu.randReg(secondNibble, lastByte)

	case 0xD:
		// 0xDxyn: Draw sprite at (Vx, Vy) with height n (set VF = collision)
		xcord := cpu.registers[secondNibble]
		ycord := cpu.registers[thirdNibble]
		return cpu.draw(uint16(xcord), uint16(ycord), finalNibble)

	case 0xE:
		// 0xEx9E and 0xExA1: Keyboard operations
		switch lastByte {
		case 0x9E:
			// 0xEx9E: Skip next instruction if key Vx is pressed
			return cpu.skipIfKeyPressed(uint16(cpu.registers[secondNibble]))
		case 0xA1:
			// 0xExA1: Skip next instruction if key Vx is not pressed
			return cpu.skipIfKeyNotPressed(uint16(cpu.registers[secondNibble]))
		}

	case 0xF:
		// 0xFx__: Timer, memory, and I/O operations
		switch lastByte {
		case 0x07:
			// 0xFx07: Set Vx = delay timer value
			return cpu.getDelayTimer(secondNibble)

		case 0x0A:
			// 0xFx0A: Wait for key press, store in Vx
			key := uint16(0)
			if cpu.keypad == [16]bool{} && cpu.waitingForKey {
				cpu.waitingForKey = true
				cpu.pc -= 2
				return nil
			}
			for i, arg := range cpu.keypad {
				if arg {
					key = uint16(i)
					break
				}
			}
			return cpu.waitForKeyPress(secondNibble, key)

		case 0x15:
			// 0xFx15: Set delay timer = Vx
			return cpu.setDelayTimer(secondNibble)

		case 0x18:
			// 0xFx18: Set sound timer = Vx
			return cpu.setSoundTimer(secondNibble)

		case 0x1E:
			// 0xFx1E: Add Vx to I
			return cpu.addI(secondNibble)

		case 0x29:
			// 0xFx29: Set I = location of sprite for digit Vx
			return cpu.setIToFont(secondNibble)

		case 0x33:
			// 0xFx33: Store BCD representation of Vx at I, I+1, I+2
			return cpu.storeBCD(secondNibble)

		case 0x55:
			// 0xFx55: Store V0 to Vx in memory starting at I
			return cpu.storeRegisters(secondNibble)

		case 0x65:
			// 0xFx65: Load V0 to Vx from memory starting at I
			return cpu.loadRegisters(secondNibble)
		}

	default:
		return errors.New("Unknown opcode")
	}

	return nil
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

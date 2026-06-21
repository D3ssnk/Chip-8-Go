// Package cpu implements tests for CHIP-8 opcodes.
package cpu

import (
	"testing"
)

// TestSetI verifies that setI (0xAnnn) correctly sets the index register (I)
// to the specified 12-bit address.
func TestSetI(t *testing.T) {
	cpu := NewCPU()

	targetAddress := uint16(0x300)

	err := cpu.setI(targetAddress)
	if err != nil {
		t.Errorf("Expected no error on setI, got %v", err)
	}

	// Verify I register was updated
	if cpu.i != targetAddress {
		t.Errorf("Expected I register to be 0x%X, got 0x%X", targetAddress, cpu.i)
	}
}

// TestSetIOutOfBounds verifies that setI returns an error when the address
// exceeds the maximum valid CHIP-8 memory address (0xFFF).
func TestSetIOutOfBounds(t *testing.T) {
	cpu := NewCPU()

	// Set an initial value to ensure it isn't overwritten
	initialI := uint16(0x200)
	cpu.i = initialI

	outOfBoundsAddress := uint16(0x1000)

	err := cpu.setI(outOfBoundsAddress)
	if err == nil {
		t.Errorf("Expected error for out of bounds address, got none")
	}

	if err.Error() != "Address is out of bounds" {
		t.Errorf("Expected error message 'Address is out of bounds', got '%v'", err.Error())
	}

	// Verify I register was NOT modified
	if cpu.i != initialI {
		t.Errorf("Expected I register to remain 0x%X, got 0x%X", initialI, cpu.i)
	}
}

// TestSetIEdgeCases verifies that setI correctly handles the minimum (0x000)
// and maximum (0xFFF) valid memory addresses.
func TestSetIEdgeCases(t *testing.T) {
	cpu := NewCPU()

	// Test minimum boundary (0x000)
	err := cpu.setI(0x000)
	if err != nil {
		t.Errorf("Expected no error for address 0x000, got %v", err)
	}
	if cpu.i != 0x000 {
		t.Errorf("Expected I register to be 0x000, got 0x%X", cpu.i)
	}

	// Test maximum boundary (0xFFF)
	err = cpu.setI(0xFFF)
	if err != nil {
		t.Errorf("Expected no error for address 0xFFF, got %v", err)
	}
	if cpu.i != 0xFFF {
		t.Errorf("Expected I register to be 0xFFF, got 0x%X", cpu.i)
	}
}

// TestAddI verifies that addI (0xFx1E) correctly adds the value of a
// register to the index register (I).
func TestAddI(t *testing.T) {
	cpu := NewCPU()

	cpu.i = uint16(0x200)
	registerIndex := uint16(0x5)
	
	// Set Vx to 0x42
	cpu.registers[registerIndex] = 0x42

	err := cpu.addI(registerIndex)
	if err != nil {
		t.Errorf("Expected no error on addI, got %v", err)
	}

	// Verify I register was updated (0x200 + 0x42 = 0x242)
	expectedI := uint16(0x242)
	if cpu.i != expectedI {
		t.Errorf("Expected I register to be 0x%X, got 0x%X", expectedI, cpu.i)
	}
}

// TestAddIInvalidRegister verifies that addI returns an error
// when the register index exceeds 0xF.
func TestAddIInvalidRegister(t *testing.T) {
	cpu := NewCPU()

	invalidRegister := uint16(0x10)

	err := cpu.addI(invalidRegister)
	if err == nil {
		t.Errorf("Expected error for invalid register, got none")
	}

	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestAddIOutOfBounds verifies that addI returns an error when the addition
// would cause the I register to exceed the valid memory boundary (0xFFF).
func TestAddIOutOfBounds(t *testing.T) {
	cpu := NewCPU()

	originalI := uint16(0xF00)
	cpu.i = originalI
	
	registerIndex := uint16(0x2)
	// 0xF00 + 0xFF = 0xFFF (valid), but 0xF00 + 0x100 would be 0x1000 (invalid)
	// Since registers are uint8, max is 0xFF. Let's start I even higher.
	
	cpu.i = uint16(0xFF0)
	cpu.registers[registerIndex] = 0x20 // 0xFF0 + 0x20 = 0x1010 (Out of bounds)

	err := cpu.addI(registerIndex)
	if err == nil {
		t.Errorf("Expected error when I register goes out of bounds, got none")
	}

	if err.Error() != "Address out of bounds" {
		t.Errorf("Expected error message 'Address out of bounds', got '%v'", err.Error())
	}

	// Verify the I register was NOT modified
	if cpu.i != uint16(0xFF0) {
		t.Errorf("Expected I register to remain 0xFF0, got 0x%X", cpu.i)
	}
}

// TestAddIEdgeCases verifies addI correctly handles adding zero 
// and reaching the exact maximum memory boundary (0xFFF).
func TestAddIEdgeCases(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x3)

	// Test adding 0
	cpu.i = uint16(0x300)
	cpu.registers[registerIndex] = 0x00
	
	err := cpu.addI(registerIndex)
	if err != nil {
		t.Errorf("Expected no error when adding 0, got %v", err)
	}
	if cpu.i != 0x300 {
		t.Errorf("Expected I register to remain 0x300, got 0x%X", cpu.i)
	}

	// Test reaching exact maximum boundary (0xFFF)
	cpu.i = uint16(0xF00)
	cpu.registers[registerIndex] = 0xFF
	
	err = cpu.addI(registerIndex)
	if err != nil {
		t.Errorf("Expected no error when reaching 0xFFF boundary, got %v", err)
	}
	if cpu.i != 0xFFF {
		t.Errorf("Expected I register to be 0xFFF, got 0x%X", cpu.i)
	}
}

// TestSetIToFont verifies that setIToFont (0xFx29) correctly sets the I register
// to the memory address of the hex font character stored in Vx.
func TestSetIToFont(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	fontChar := uint8(0xA) // Font 'A'
	cpu.registers[registerIndex] = fontChar

	err := cpu.setIToFont(registerIndex)
	if err != nil {
		t.Errorf("Expected no error on setIToFont, got %v", err)
	}

	// Each font sprite is 5 bytes. The offset for 'A' (10) should be 50 (0x32).
	expectedI := uint16(fontChar * 5)
	if cpu.i != expectedI {
		t.Errorf("Expected I register to be 0x%X, got 0x%X", expectedI, cpu.i)
	}
}

// TestSetIToFontInvalidRegister verifies that setIToFont returns an error
// when the register index exceeds 0xF.
func TestSetIToFontInvalidRegister(t *testing.T) {
	cpu := NewCPU()

	invalidRegister := uint16(0x10)

	err := cpu.setIToFont(invalidRegister)
	if err == nil {
		t.Errorf("Expected error for invalid register, got none")
	}

	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestSetIToFontInvalidFont verifies that setIToFont returns an error
// when the value stored in the register is larger than 0xF (not a valid hex font).
func TestSetIToFontInvalidFont(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x2)
	invalidFontChar := uint8(0x10) // CHIP-8 only has fonts for 0x0 through 0xF
	cpu.registers[registerIndex] = invalidFontChar

	err := cpu.setIToFont(registerIndex)
	if err == nil {
		t.Errorf("Expected error for invalid font character, got none")
	}

	if err.Error() != "Not a font" {
		t.Errorf("Expected error message 'Not a font', got '%v'", err.Error())
	}
}

// TestSetIToFontEdgeCases verifies setIToFont works correctly at the absolute
// minimum (0) and maximum (F) valid font character boundaries.
func TestSetIToFontEdgeCases(t *testing.T) {
	cpu := NewCPU()
	registerIndex := uint16(0x3)

	// Test minimum boundary: Font '0'
	cpu.registers[registerIndex] = 0x0
	err := cpu.setIToFont(registerIndex)
	if err != nil {
		t.Errorf("Expected no error for font 0, got %v", err)
	}
	if cpu.i != 0x000 {
		t.Errorf("Expected I register to be 0x000 for font 0, got 0x%X", cpu.i)
	}

	// Test maximum boundary: Font 'F' (15)
	cpu.registers[registerIndex] = 0xF
	err = cpu.setIToFont(registerIndex)
	if err != nil {
		t.Errorf("Expected no error for font F, got %v", err)
	}
	
	expectedI := uint16(15 * 5) // 75 (0x4B)
	if cpu.i != expectedI {
		t.Errorf("Expected I register to be 0x%X for font F, got 0x%X", expectedI, cpu.i)
	}
}

// TestSetIToFontIsolation verifies that calculating the font offset
// does not accidentally mutate the source register.
func TestSetIToFontIsolation(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x7)
	fontChar := uint8(0xB)
	cpu.registers[registerIndex] = fontChar

	err := cpu.setIToFont(registerIndex)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the register still holds the base font character
	if cpu.registers[registerIndex] != fontChar {
		t.Errorf("Expected register[%d] to remain 0x%X, got 0x%X", registerIndex, fontChar, cpu.registers[registerIndex])
	}
}

// TestStoreBCD verifies that storeBCD (0xFx33) correctly extracts the hundreds,
// tens, and units digits of a register and stores them in sequential memory.
func TestStoreBCD(t *testing.T) {
	cpu := NewCPU()

	cpu.i = uint16(0x300)
	registerIndex := uint16(0x5)
	
	// 234 -> Hundreds: 2, Tens: 3, Units: 4
	cpu.registers[registerIndex] = 234 

	err := cpu.storeBCD(registerIndex)
	if err != nil {
		t.Errorf("Expected no error on storeBCD, got %v", err)
	}

	// Verify memory locations
	if cpu.memory[cpu.i] != 2 {
		t.Errorf("Expected hundreds digit to be 2, got %d", cpu.memory[cpu.i])
	}
	if cpu.memory[cpu.i+1] != 3 {
		t.Errorf("Expected tens digit to be 3, got %d", cpu.memory[cpu.i+1])
	}
	if cpu.memory[cpu.i+2] != 4 {
		t.Errorf("Expected units digit to be 4, got %d", cpu.memory[cpu.i+2])
	}
}

// TestStoreBCDEdgeCases verifies storeBCD correctly handles the minimum (0)
// and maximum (255) possible values of an 8-bit register.
func TestStoreBCDEdgeCases(t *testing.T) {
	cpu := NewCPU()
	registerIndex := uint16(0x2)
	cpu.i = uint16(0x400)

	// Test absolute minimum: 0
	cpu.registers[registerIndex] = 0
	err := cpu.storeBCD(registerIndex)
	if err != nil {
		t.Errorf("Expected no error for BCD 0, got %v", err)
	}
	if cpu.memory[cpu.i] != 0 || cpu.memory[cpu.i+1] != 0 || cpu.memory[cpu.i+2] != 0 {
		t.Errorf("Expected memory to hold [0, 0, 0], got [%d, %d, %d]", 
			cpu.memory[cpu.i], cpu.memory[cpu.i+1], cpu.memory[cpu.i+2])
	}

	// Test absolute maximum: 255
	cpu.registers[registerIndex] = 255
	err = cpu.storeBCD(registerIndex)
	if err != nil {
		t.Errorf("Expected no error for BCD 255, got %v", err)
	}
	if cpu.memory[cpu.i] != 2 || cpu.memory[cpu.i+1] != 5 || cpu.memory[cpu.i+2] != 5 {
		t.Errorf("Expected memory to hold [2, 5, 5], got [%d, %d, %d]", 
			cpu.memory[cpu.i], cpu.memory[cpu.i+1], cpu.memory[cpu.i+2])
	}
}

// TestStoreBCDInvalidRegister verifies that storeBCD returns an error
// when the register index exceeds 0xF.
func TestStoreBCDInvalidRegister(t *testing.T) {
	cpu := NewCPU()

	invalidRegister := uint16(0x10)

	err := cpu.storeBCD(invalidRegister)
	if err == nil {
		t.Errorf("Expected error for invalid register, got none")
	}

	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestStoreBCDOutOfBoundsMemory verifies that storeBCD returns an error
// when attempting to write BCD values past the maximum memory boundary (0xFFF).
func TestStoreBCDOutOfBoundsMemory(t *testing.T) {
	cpu := NewCPU()
	registerIndex := uint16(0x0)
	cpu.registers[registerIndex] = 123
	
	// Position I so that I+2 equals 0x1000 (which is out of bounds)
	cpu.i = uint16(0xFFE) 

	err := cpu.storeBCD(registerIndex)
	if err == nil {
		t.Errorf("Expected error for out of bounds memory write, got none")
	}

	if err.Error() != "Address out of bounds" {
		t.Errorf("Expected error message 'Address out of bounds', got '%v'", err.Error())
	}
}
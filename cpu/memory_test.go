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
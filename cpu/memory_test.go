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
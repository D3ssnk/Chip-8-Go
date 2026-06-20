// Package cpu implements tests for the CHIP-8 CPU emulator.
package cpu

import (
	"testing"
)

// TestFetch verifies that the Fetch method correctly reads 16-bit instructions
// from memory and increments the program counter.
func TestFetch(t *testing.T) {
	cpu := NewCPU()

	// Set up test instruction bytes at memory location 0x200 (program start)
	cpu.memory[0x200] = 0x12 // High byte
	cpu.memory[0x201] = 0x34 // Low byte
	cpu.memory[0x202] = 0x56 // Next instruction high byte
	cpu.memory[0x203] = 0x78 // Next instruction low byte

	// Test: First fetch should return 0x1234 with no error
	instruction, err := cpu.fetch()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if instruction != 0x1234 {
		t.Errorf("Expected instruction 0x1234, got 0x%X", instruction)
	}

	// Test: PC should be incremented to 0x202
	if cpu.pc != 0x202 {
		t.Errorf("Expected PC to be 0x202, got 0x%X", cpu.pc)
	}

	// Test: Second fetch should return 0x5678 with no error
	instruction, err = cpu.fetch()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if instruction != 0x5678 {
		t.Errorf("Expected instruction 0x5678, got 0x%X", instruction)
	}

	// Test: PC should be incremented to 0x204
	if cpu.pc != 0x204 {
		t.Errorf("Expected PC to be 0x204, got 0x%X", cpu.pc)
	}
}

// TestFetchEdgeCase verifies Fetch works correctly with edge case values
// including 0x0000 and 0xFFFF.
func TestFetchEdgeCase(t *testing.T) {
	cpu := NewCPU()

	// Test edge case: 0x0000
	cpu.memory[0x200] = 0x00
	cpu.memory[0x201] = 0x00

	instruction, err := cpu.fetch()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if instruction != 0x0000 {
		t.Errorf("Expected instruction 0x0000, got 0x%X", instruction)
	}

	// Test edge case: 0xFFFF
	cpu.memory[0x202] = 0xFF
	cpu.memory[0x203] = 0xFF

	instruction, err = cpu.fetch()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if instruction != 0xFFFF {
		t.Errorf("Expected instruction 0xFFFF, got 0x%X", instruction)
	}
}

// TestFetchBoundaryError verifies that Fetch returns an error when
// the program counter exceeds the memory limit (0xFFF).
func TestFetchBoundaryError(t *testing.T) {
	cpu := NewCPU()

	// Set PC to a position that will exceed 0xFFF after fetching
	cpu.pc = 0xFFF

	instruction, err := cpu.fetch()
	if err == nil {
		t.Errorf("Expected error when PC exceeds memory, got none")
	}
	if instruction != 0 {
		t.Errorf("Expected instruction 0 on error, got 0x%X", instruction)
	}
}

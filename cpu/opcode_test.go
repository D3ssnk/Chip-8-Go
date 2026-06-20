// Package cpu implements tests for CHIP-8 opcodes.
package cpu

import (
	"testing"
)

// TestClear verifies that the clear opcode (0x00E0) properly clears the display.
// The display should be set to all zeros after the clear operation.
func TestClear(t *testing.T) {
	cpu := NewCPU()

	// Fill the display with non-zero values
	for i := 0; i < 32; i++ {
		for j := 0; j < 8; j++ {
			cpu.display[i][j] = 0xFF
		}
	}

	// Verify display is filled
	for i := 0; i < 32; i++ {
		for j := 0; j < 8; j++ {
			if cpu.display[i][j] != 0xFF {
				t.Errorf("Expected display[%d][%d] to be 0xFF before clear, got 0x%X", i, j, cpu.display[i][j])
			}
		}
	}

	// Execute clear opcode
	cpu.clear()

	// Verify display is cleared (all zeros)
	for i := 0; i < 32; i++ {
		for j := 0; j < 8; j++ {
			if cpu.display[i][j] != 0x00 {
				t.Errorf("Expected display[%d][%d] to be 0x00 after clear, got 0x%X", i, j, cpu.display[i][j])
			}
		}
	}
}

// TestClearMultipleTimes verifies that clear can be called multiple times
// and the display remains cleared after each call.
func TestClearMultipleTimes(t *testing.T) {
	cpu := NewCPU()

	for iteration := 0; iteration < 3; iteration++ {
		// Fill display with pattern
		for i := 0; i < 32; i++ {
			for j := 0; j < 8; j++ {
				cpu.display[i][j] = uint8((i + j) % 256)
			}
		}

		// Clear display
		cpu.clear()

		// Verify all pixels are cleared
		for i := 0; i < 32; i++ {
			for j := 0; j < 8; j++ {
				if cpu.display[i][j] != 0x00 {
					t.Errorf("Iteration %d: Expected display[%d][%d] to be 0x00, got 0x%X", iteration, i, j, cpu.display[i][j])
				}
			}
		}
	}
}

// TestReturn verifies that the return opcode (0x00EE) correctly pops a return
// address from the stack and sets the program counter to that address.
func TestReturn(t *testing.T) {
	cpu := NewCPU()

	// Set up stack with a return address
	returnAddress := uint16(0x300)
	cpu.sp = 1
	cpu.stack[1] = returnAddress

	// Execute return opcode
	err := cpu.ret()
	if err != nil {
		t.Errorf("Expected no error on return, got %v", err)
	}

	// Verify PC is set to the return address
	if cpu.pc != returnAddress {
		t.Errorf("Expected PC to be 0x%X, got 0x%X", returnAddress, cpu.pc)
	}

	// Verify stack pointer was decremented
	if cpu.sp != 0 {
		t.Errorf("Expected SP to be 0, got %d", cpu.sp)
	}
}

// TestReturnEmptyStack verifies that ret returns an error when the stack is empty (SP == 0).
func TestReturnEmptyStack(t *testing.T) {
	cpu := NewCPU()

	// Ensure stack pointer is at 0 (empty stack)
	cpu.sp = 0

	// Attempt to execute return opcode
	err := cpu.ret()
	if err == nil {
		t.Errorf("Expected error when returning from empty stack, got none")
	}

	// Verify error message
	if err.Error() != "Stack is empty" {
		t.Errorf("Expected error message 'Stack is empty', got '%v'", err.Error())
	}
}

// TestReturnMultiple verifies that multiple consecutive return operations work correctly.
func TestReturnMultiple(t *testing.T) {
	cpu := NewCPU()

	// Set up stack with multiple return addresses
	addresses := []uint16{0x200, 0x300, 0x400, 0x500}

	for i, addr := range addresses {
		cpu.stack[i+1] = addr
	}
	cpu.sp = uint8(len(addresses))

	// Pop all addresses off the stack
	for i := len(addresses) - 1; i >= 0; i-- {
		err := cpu.ret()
		if err != nil {
			t.Errorf("Iteration %d: Expected no error, got %v", i, err)
		}

		if cpu.pc != addresses[i] {
			t.Errorf("Iteration %d: Expected PC to be 0x%X, got 0x%X", i, addresses[i], cpu.pc)
		}

		if cpu.sp != uint8(i) {
			t.Errorf("Iteration %d: Expected SP to be %d, got %d", i, i, cpu.sp)
		}
	}

	// Verify stack is now empty and next return fails
	err := cpu.ret()
	if err == nil {
		t.Errorf("Expected error when stack is empty after multiple returns, got none")
	}
}

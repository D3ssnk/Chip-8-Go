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

// TestJump verifies that the jump opcode (0x1nnn) correctly sets the program counter
// to the specified address.
func TestJump(t *testing.T) {
	cpu := NewCPU()

	jumpAddress := uint16(0x300)

	// Execute jump opcode
	err := cpu.jump(jumpAddress)
	if err != nil {
		t.Errorf("Expected no error on jump, got %v", err)
	}

	// Verify PC is set to the jump address
	if cpu.pc != jumpAddress {
		t.Errorf("Expected PC to be 0x%X, got 0x%X", jumpAddress, cpu.pc)
	}
}

// TestJumpOutOfBounds verifies that jump returns an error when the address
// exceeds the maximum valid memory address (0xFFF).
func TestJumpOutOfBounds(t *testing.T) {
	cpu := NewCPU()

	// Try to jump to an out-of-bounds address
	outOfBoundsAddress := uint16(0x1000)

	err := cpu.jump(outOfBoundsAddress)
	if err == nil {
		t.Errorf("Expected error when jumping out of bounds, got none")
	}

	// Verify error message
	if err.Error() != "Address is out of bounds" {
		t.Errorf("Expected error message 'Address is out of bounds', got '%v'", err.Error())
	}

	// Verify PC was not changed
	if cpu.pc != 0x200 {
		t.Errorf("Expected PC to remain 0x200, got 0x%X", cpu.pc)
	}
}

// TestJumpBoundaryValues verifies jump works with edge case addresses
// at the boundary of valid memory.
func TestJumpBoundaryValues(t *testing.T) {
	cpu := NewCPU()

	// Test jump to minimum valid address (after font data)
	err := cpu.jump(0x200)
	if err != nil {
		t.Errorf("Expected no error jumping to 0x200, got %v", err)
	}
	if cpu.pc != 0x200 {
		t.Errorf("Expected PC to be 0x200, got 0x%X", cpu.pc)
	}

	// Test jump to maximum valid address
	err = cpu.jump(0xFFF)
	if err != nil {
		t.Errorf("Expected no error jumping to 0xFFF, got %v", err)
	}
	if cpu.pc != 0xFFF {
		t.Errorf("Expected PC to be 0xFFF, got 0x%X", cpu.pc)
	}

	// Test jump to first address beyond bounds
	err = cpu.jump(0x1000)
	if err == nil {
		t.Errorf("Expected error jumping to 0x1000, got none")
	}
}

// TestJumpMultiple verifies that multiple consecutive jumps work correctly.
func TestJumpMultiple(t *testing.T) {
	cpu := NewCPU()

	addresses := []uint16{0x200, 0x400, 0x600, 0x800, 0xA00}

	for _, addr := range addresses {
		err := cpu.jump(addr)
		if err != nil {
			t.Errorf("Expected no error jumping to 0x%X, got %v", addr, err)
		}

		if cpu.pc != addr {
			t.Errorf("Expected PC to be 0x%X, got 0x%X", addr, cpu.pc)
		}
	}
}

// TestCall verifies that the call opcode (0x2nnn) correctly pushes the current PC
// onto the stack and sets PC to the subroutine address.
func TestCall(t *testing.T) {
	cpu := NewCPU()

	callAddress := uint16(0x400)
	originalPC := cpu.pc

	// Execute call opcode
	err := cpu.call(callAddress)
	if err != nil {
		t.Errorf("Expected no error on call, got %v", err)
	}

	// Verify PC is set to the call address
	if cpu.pc != callAddress {
		t.Errorf("Expected PC to be 0x%X, got 0x%X", callAddress, cpu.pc)
	}

	// Verify stack pointer was incremented
	if cpu.sp != 1 {
		t.Errorf("Expected SP to be 1, got %d", cpu.sp)
	}

	// Verify return address was pushed onto stack
	if cpu.stack[1] != originalPC {
		t.Errorf("Expected stack[1] to be 0x%X, got 0x%X", originalPC, cpu.stack[1])
	}
}

// TestCallOutOfBounds verifies that call returns an error when the address
// exceeds the maximum valid memory address (0xFFF).
func TestCallOutOfBounds(t *testing.T) {
	cpu := NewCPU()

	outOfBoundsAddress := uint16(0x1000)

	err := cpu.call(outOfBoundsAddress)
	if err == nil {
		t.Errorf("Expected error when calling out of bounds, got none")
	}

	// Verify error message
	if err.Error() != "Address is out of bounds" {
		t.Errorf("Expected error message 'Address is out of bounds', got '%v'", err.Error())
	}

	// Verify PC and SP were not changed
	if cpu.pc != 0x200 {
		t.Errorf("Expected PC to remain 0x200, got 0x%X", cpu.pc)
	}
	if cpu.sp != 0 {
		t.Errorf("Expected SP to remain 0, got %d", cpu.sp)
	}
}

// TestCallStackFull verifies that call returns an error when the stack is full (SP == 0xF).
func TestCallStackFull(t *testing.T) {
	cpu := NewCPU()

	// Fill the stack to capacity
	cpu.sp = 0xF

	callAddress := uint16(0x400)

	err := cpu.call(callAddress)
	if err == nil {
		t.Errorf("Expected error when stack is full, got none")
	}

	// Verify error message
	if err.Error() != "Stack is full" {
		t.Errorf("Expected error message 'Stack is full', got '%v'", err.Error())
	}

	// Verify PC and SP were not changed
	if cpu.pc != 0x200 {
		t.Errorf("Expected PC to remain 0x200, got 0x%X", cpu.pc)
	}
	if cpu.sp != 0xF {
		t.Errorf("Expected SP to remain 0xF, got %d", cpu.sp)
	}
}

// TestCallMultiple verifies that multiple consecutive calls properly manage the stack
// and can be correctly unwound using return operations.
func TestCallMultiple(t *testing.T) {
	cpu := NewCPU()

	callAddresses := []uint16{0x300, 0x400, 0x500, 0x600}

	// Perform multiple calls
	for i, addr := range callAddresses {
		err := cpu.call(addr)
		if err != nil {
			t.Errorf("Call %d: Expected no error, got %v", i, err)
		}

		// Verify PC is set to call address
		if cpu.pc != addr {
			t.Errorf("Call %d: Expected PC to be 0x%X, got 0x%X", i, addr, cpu.pc)
		}

		// Verify stack pointer is correct
		expectedSP := uint8(i + 1)
		if cpu.sp != expectedSP {
			t.Errorf("Call %d: Expected SP to be %d, got %d", i, expectedSP, cpu.sp)
		}
	}

	// Verify we can unwind the stack with returns
	for i := len(callAddresses) - 1; i >= 0; i-- {
		err := cpu.ret()
		if err != nil {
			t.Errorf("Return %d: Expected no error, got %v", i, err)
		}

		// For each return, verify we get back the original PC
		// (which would have been updated by previous calls)
		if cpu.sp != uint8(i) {
			t.Errorf("Return %d: Expected SP to be %d, got %d", i, i, cpu.sp)
		}
	}
}
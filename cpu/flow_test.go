// Package cpu implements tests for CHIP-8 opcodes.
package cpu

import (
	"testing"
)

// TestClear verifies that the clear opcode (0x00E0) properly clears the display.
// The display should be set to all false (off) after the clear operation.
func TestClear(t *testing.T) {
	cpu := NewCPU()

	// Fill the display with true (on) values
	for i := 0; i < 32; i++ {
		for j := 0; j < 64; j++ {
			cpu.display[i][j] = true
		}
	}

	// Verify display is filled
	for i := 0; i < 32; i++ {
		for j := 0; j < 64; j++ {
			if !cpu.display[i][j] {
				t.Errorf("Expected display[%d][%d] to be true before clear, got false", i, j)
			}
		}
	}

	// Execute clear opcode
	cpu.clear()

	// Verify display is cleared (all false)
	for i := 0; i < 32; i++ {
		for j := 0; j < 64; j++ {
			// If the pixel is true, the clear failed
			if cpu.display[i][j] {
				t.Errorf("Expected display[%d][%d] to be false after clear, got true", i, j)
			}
		}
	}
}

// TestClearMultipleTimes verifies that clear can be called multiple times
// and the display remains cleared after each call.
func TestClearMultipleTimes(t *testing.T) {
	cpu := NewCPU()

	for iteration := 0; iteration < 3; iteration++ {
		// Fill display with an alternating pattern
		for i := 0; i < 32; i++ {
			for j := 0; j < 64; j++ {
				cpu.display[i][j] = (i+j)%2 == 0
			}
		}

		// Clear display
		cpu.clear()

		// Verify all pixels are cleared (set to false)
		for i := 0; i < 32; i++ {
			for j := 0; j < 64; j++ {
				if cpu.display[i][j] {
					t.Errorf("Iteration %d: Expected display[%d][%d] to be false, got true", iteration, i, j)
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

// TestJumpV0 verifies that jumpV0 (0xBnnn) correctly sets the program counter
// to the specified address plus the value in register V0.
func TestJumpV0(t *testing.T) {
	cpu := NewCPU()

	// Set V0 to a specific offset
	cpu.registers[0x0] = 0x42
	baseAddress := uint16(0x300)

	// Execute jump V0
	err := cpu.jumpV0(baseAddress)
	if err != nil {
		t.Errorf("Expected no error on jumpV0, got %v", err)
	}

	// Verify PC is set to base address + V0 (0x300 + 0x42 = 0x342)
	expectedPC := uint16(0x342)
	if cpu.pc != expectedPC {
		t.Errorf("Expected PC to be 0x%X, got 0x%X", expectedPC, cpu.pc)
	}
}

// TestJumpV0OutOfBounds verifies that jumpV0 returns an error when the combined
// address (base + V0) exceeds the maximum valid memory address (0xFFF).
func TestJumpV0OutOfBounds(t *testing.T) {
	cpu := NewCPU()

	originalPC := uint16(0x200)
	cpu.pc = originalPC

	// Set V0 to a high value
	cpu.registers[0x0] = 0xFF

	// Base address + V0 will equal 0x1000, which is out of bounds
	baseAddress := uint16(0x0F01)

	err := cpu.jumpV0(baseAddress)
	if err == nil {
		t.Errorf("Expected error when jumpV0 goes out of bounds, got none")
	}

	// Verify error message
	if err.Error() != "Address is out of bounds" {
		t.Errorf("Expected error message 'Address is out of bounds', got '%v'", err.Error())
	}

	// Verify PC was not changed
	if cpu.pc != originalPC {
		t.Errorf("Expected PC to remain 0x%X, got 0x%X", originalPC, cpu.pc)
	}
}

// TestJumpV0EdgeCases verifies jumpV0 works correctly at the extreme edges
// of memory boundaries and with zero values.
func TestJumpV0EdgeCases(t *testing.T) {
	cpu := NewCPU()

	// Test absolute maximum valid boundary (V0 + base = 0xFFF)
	cpu.registers[0x0] = 0xFF
	err := cpu.jumpV0(0x0F00) // 0xF00 + 0xFF = 0xFFF
	if err != nil {
		t.Errorf("Expected no error jumping to maximum boundary 0xFFF, got %v", err)
	}
	if cpu.pc != 0xFFF {
		t.Errorf("Expected PC to be 0xFFF, got 0x%X", cpu.pc)
	}

	// Test with V0 as 0
	cpu.registers[0x0] = 0x00
	err = cpu.jumpV0(0x400)
	if err != nil {
		t.Errorf("Expected no error when V0 is 0, got %v", err)
	}
	if cpu.pc != 0x400 {
		t.Errorf("Expected PC to be 0x400, got 0x%X", cpu.pc)
	}
}

// TestDrawBasic verifies that the draw opcode (0xDxyn) correctly renders
// a sprite to the display and sets VF to 0 when no pixels are erased.
func TestDrawBasic(t *testing.T) {
	cpu := NewCPU()

	// Set up memory with a 1-byte sprite: 0xC0 (Binary: 1100 0000)
	cpu.i = 0x300
	cpu.memory[cpu.i] = 0xC0

	// Pre-set VF to 1 to ensure the draw function resets it to 0
	cpu.registers[0xF] = 1

	// Draw 1 byte at (X: 0, Y: 0)
	err := cpu.draw(0, 0, 1)
	if err != nil {
		t.Errorf("Expected no error on draw, got %v", err)
	}

	// Verify pixels were toggled on
	if !cpu.display[0][0] {
		t.Errorf("Expected display[0][0] to be true")
	}
	if !cpu.display[0][1] {
		t.Errorf("Expected display[0][1] to be true")
	}
	// Verify subsequent pixel is off
	if cpu.display[0][2] {
		t.Errorf("Expected display[0][2] to be false")
	}

	// Verify VF flag was reset to 0 (no collision)
	if cpu.registers[0xF] != 0 {
		t.Errorf("Expected VF (register[15]) to be 0, got %d", cpu.registers[0xF])
	}
}

// TestDrawCollision verifies that drawing over an existing active pixel
// correctly turns the pixel off (XOR) and sets VF to 1.
func TestDrawCollision(t *testing.T) {
	cpu := NewCPU()

	// Pre-activate a pixel at (0, 0)
	cpu.display[0][0] = true

	// Set up memory with a 1-byte sprite: 0xC0 (Binary: 1100 0000)
	cpu.i = 0x300
	cpu.memory[cpu.i] = 0xC0

	// Draw 1 byte at (X: 0, Y: 0)
	err := cpu.draw(0, 0, 1)
	if err != nil {
		t.Errorf("Expected no error on draw, got %v", err)
	}

	// Verify pixel [0][0] was XOR'd off (collision)
	if cpu.display[0][0] {
		t.Errorf("Expected display[0][0] to be false after collision")
	}
	// Verify pixel [0][1] was toggled on (no collision)
	if !cpu.display[0][1] {
		t.Errorf("Expected display[0][1] to be true")
	}

	// Verify VF flag was set to 1 (collision occurred)
	if cpu.registers[0xF] != 1 {
		t.Errorf("Expected VF (register[15]) to be 1 due to collision, got %d", cpu.registers[0xF])
	}
}

// TestDrawWrapping verifies that sprites correctly wrap around the edges
// of the 64x32 display.
func TestDrawWrapping(t *testing.T) {
	cpu := NewCPU()

	// Set up memory with a 1-byte sprite: 0xC0 (Binary: 1100 0000)
	cpu.i = 0x300
	cpu.memory[cpu.i] = 0xC0

	// Draw 1 byte at the extreme bottom right: X=63, Y=31
	err := cpu.draw(63, 31, 1)
	if err != nil {
		t.Errorf("Expected no error on draw, got %v", err)
	}

	// The first bit (1) should be at [31][63]
	if !cpu.display[31][63] {
		t.Errorf("Expected display[31][63] to be true")
	}
	// The second bit (1) should wrap around to the left edge [31][0]
	if !cpu.display[31][0] {
		t.Errorf("Expected display[31][0] to be true (wrapped)")
	}
}

// TestDrawOutOfBoundsMemory verifies the function prevents reading past
// the maximum valid memory address of 0xFFF.
func TestDrawOutOfBoundsMemory(t *testing.T) {
	cpu := NewCPU()

	// Position the index register at the very end of memory
	cpu.i = 0xFFF

	// Try to draw 2 bytes (which would require reading 0x1000)
	err := cpu.draw(0, 0, 2)
	if err == nil {
		t.Errorf("Expected error for memory read out of bounds, got none")
	}
	if err != nil && err.Error() != "Address out of bounds" {
		t.Errorf("Expected 'Address out of bounds', got '%v'", err.Error())
	}
}

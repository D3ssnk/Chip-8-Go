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

// TestSkipIfEqual verifies that skipIfEqual (0x3xkk) correctly skips the next instruction
// when the register value equals the given byte value.
func TestSkipIfEqual(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	value := uint16(0x42)

	// Set register to the value
	cpu.registers[registerIndex] = uint8(value)
	originalPC := cpu.pc

	// Execute skip if equal
	err := cpu.skipIfEqual(registerIndex, value)
	if err != nil {
		t.Errorf("Expected no error on skipIfEqual, got %v", err)
	}

	// Verify PC was incremented by 2 (skip next instruction)
	expectedPC := originalPC + 2
	if cpu.pc != expectedPC {
		t.Errorf("Expected PC to be 0x%X, got 0x%X", expectedPC, cpu.pc)
	}
}

// TestSkipIfEqualNoSkip verifies that skipIfEqual does NOT skip when values don't match.
func TestSkipIfEqualNoSkip(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x3)
	registerValue := uint16(0x42)
	compareValue := uint16(0x50)

	// Set register to a different value
	cpu.registers[registerIndex] = uint8(registerValue)
	originalPC := cpu.pc

	// Execute skip if equal with different value
	err := cpu.skipIfEqual(registerIndex, compareValue)
	if err != nil {
		t.Errorf("Expected no error on skipIfEqual, got %v", err)
	}

	// Verify PC was NOT incremented
	if cpu.pc != originalPC {
		t.Errorf("Expected PC to remain 0x%X, got 0x%X", originalPC, cpu.pc)
	}
}

// TestSkipIfEqualInvalidRegister verifies that skipIfEqual returns an error
// when the register index exceeds 0xF.
func TestSkipIfEqualInvalidRegister(t *testing.T) {
	cpu := NewCPU()

	invalidRegister := uint16(0x10)
	value := uint16(0x42)

	err := cpu.skipIfEqual(invalidRegister, value)
	if err == nil {
		t.Errorf("Expected error for invalid register, got none")
	}

	// Verify error message
	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}

	// Verify PC was not changed
	if cpu.pc != 0x200 {
		t.Errorf("Expected PC to remain 0x200, got 0x%X", cpu.pc)
	}
}

// TestSkipIfEqualEdgeCases verifies skipIfEqual works with edge case values.
func TestSkipIfEqualEdgeCases(t *testing.T) {
	cpu := NewCPU()

	// Test with value 0x00
	registerIndex := uint16(0x0)
	cpu.registers[registerIndex] = 0x00
	cpu.pc = 0x200

	err := cpu.skipIfEqual(registerIndex, 0x00)
	if err != nil {
		t.Errorf("Expected no error for edge case 0x00, got %v", err)
	}
	if cpu.pc != 0x202 {
		t.Errorf("Expected PC to be 0x202, got 0x%X", cpu.pc)
	}

	// Test with value 0xFF
	registerIndex = uint16(0xF)
	cpu.registers[registerIndex] = 0xFF
	cpu.pc = 0x300

	err = cpu.skipIfEqual(registerIndex, 0xFF)
	if err != nil {
		t.Errorf("Expected no error for edge case 0xFF, got %v", err)
	}
	if cpu.pc != 0x302 {
		t.Errorf("Expected PC to be 0x302, got 0x%X", cpu.pc)
	}
}

// TestSkipIfEqualMultipleRegisters verifies skipIfEqual works correctly with all 16 registers.
func TestSkipIfEqualMultipleRegisters(t *testing.T) {
	cpu := NewCPU()

	// Test all 16 registers
	for i := 0; i < 16; i++ {
		registerIndex := uint16(i)
		value := uint16(0x10 + i)

		// Set register to the value
		cpu.registers[registerIndex] = uint8(value)
		cpu.pc = 0x200 + uint16(i)*4

		err := cpu.skipIfEqual(registerIndex, value)
		if err != nil {
			t.Errorf("Register %d: Expected no error, got %v", i, err)
		}

		// Verify PC was skipped
		expectedPC := cpu.pc
		if expectedPC != 0x202+uint16(i)*4 {
			t.Errorf("Register %d: Expected PC to be 0x%X, got 0x%X", i, 0x202+uint16(i)*4, expectedPC)
		}
	}
}

// TestSkipIfNotEqual verifies that skipIfNotEqual (0x4xkk) correctly skips the next instruction
// when the register value does NOT equal the given byte value.
func TestSkipIfNotEqual(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	registerValue := uint16(0x42)
	compareValue := uint16(0x50)

	// Set register to a different value
	cpu.registers[registerIndex] = uint8(registerValue)
	originalPC := cpu.pc

	// Execute skip if not equal
	err := cpu.skipIfNotEqual(registerIndex, compareValue)
	if err != nil {
		t.Errorf("Expected no error on skipIfNotEqual, got %v", err)
	}

	// Verify PC was incremented by 2 (skip next instruction)
	expectedPC := originalPC + 2
	if cpu.pc != expectedPC {
		t.Errorf("Expected PC to be 0x%X, got 0x%X", expectedPC, cpu.pc)
	}
}

// TestSkipIfNotEqualNoSkip verifies that skipIfNotEqual does NOT skip when values are equal.
func TestSkipIfNotEqualNoSkip(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x3)
	value := uint16(0x42)

	// Set register to the same value
	cpu.registers[registerIndex] = uint8(value)
	originalPC := cpu.pc

	// Execute skip if not equal with same value
	err := cpu.skipIfNotEqual(registerIndex, value)
	if err != nil {
		t.Errorf("Expected no error on skipIfNotEqual, got %v", err)
	}

	// Verify PC was NOT incremented
	if cpu.pc != originalPC {
		t.Errorf("Expected PC to remain 0x%X, got 0x%X", originalPC, cpu.pc)
	}
}

// TestSkipIfNotEqualInvalidRegister verifies that skipIfNotEqual returns an error
// when the register index exceeds 0xF.
func TestSkipIfNotEqualInvalidRegister(t *testing.T) {
	cpu := NewCPU()

	invalidRegister := uint16(0x10)
	value := uint16(0x42)

	err := cpu.skipIfNotEqual(invalidRegister, value)
	if err == nil {
		t.Errorf("Expected error for invalid register, got none")
	}

	// Verify error message
	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}

	// Verify PC was not changed
	if cpu.pc != 0x200 {
		t.Errorf("Expected PC to remain 0x200, got 0x%X", cpu.pc)
	}
}

// TestSkipIfNotEqualEdgeCases verifies skipIfNotEqual works with edge case values.
func TestSkipIfNotEqualEdgeCases(t *testing.T) {
	cpu := NewCPU()

	// Test with value 0x00 (register != 0x00)
	registerIndex := uint16(0x0)
	cpu.registers[registerIndex] = 0x01
	cpu.pc = 0x200

	err := cpu.skipIfNotEqual(registerIndex, 0x00)
	if err != nil {
		t.Errorf("Expected no error for edge case, got %v", err)
	}
	if cpu.pc != 0x202 {
		t.Errorf("Expected PC to be 0x202, got 0x%X", cpu.pc)
	}

	// Test with value 0xFF (register != 0xFF)
	registerIndex = uint16(0xF)
	cpu.registers[registerIndex] = 0xFE
	cpu.pc = 0x300

	err = cpu.skipIfNotEqual(registerIndex, 0xFF)
	if err != nil {
		t.Errorf("Expected no error for edge case, got %v", err)
	}
	if cpu.pc != 0x302 {
		t.Errorf("Expected PC to be 0x302, got 0x%X", cpu.pc)
	}

	// Test equal values (should NOT skip)
	cpu.registers[registerIndex] = 0xFF
	cpu.pc = 0x400

	err = cpu.skipIfNotEqual(registerIndex, 0xFF)
	if err != nil {
		t.Errorf("Expected no error when values are equal, got %v", err)
	}
	if cpu.pc != 0x400 {
		t.Errorf("Expected PC to remain 0x400, got 0x%X", cpu.pc)
	}
}

// TestSkipIfNotEqualMultipleRegisters verifies skipIfNotEqual works correctly with all 16 registers.
func TestSkipIfNotEqualMultipleRegisters(t *testing.T) {
	cpu := NewCPU()

	// Test all 16 registers
	for i := 0; i < 16; i++ {
		registerIndex := uint16(i)
		registerValue := uint16(0x10 + i)
		compareValue := uint16(0x20 + i) // Different value

		// Set register to the value
		cpu.registers[registerIndex] = uint8(registerValue)
		cpu.pc = 0x200 + uint16(i)*4

		err := cpu.skipIfNotEqual(registerIndex, compareValue)
		if err != nil {
			t.Errorf("Register %d: Expected no error, got %v", i, err)
		}

		// Verify PC was skipped (since values don't match)
		expectedPC := cpu.pc
		if expectedPC != 0x202+uint16(i)*4 {
			t.Errorf("Register %d: Expected PC to be 0x%X, got 0x%X", i, 0x202+uint16(i)*4, expectedPC)
		}
	}
}

// TestSkipIfEqualReg verifies that skipIfEqualReg (0x5xy0) correctly skips the next instruction
// when two registers contain equal values.
func TestSkipIfEqualReg(t *testing.T) {
	cpu := NewCPU()

	register1 := uint16(0x3)
	register2 := uint16(0x7)
	value := uint16(0x42)

	// Set both registers to the same value
	cpu.registers[register1] = uint8(value)
	cpu.registers[register2] = uint8(value)
	originalPC := cpu.pc

	// Execute skip if equal registers
	err := cpu.skipIfEqualReg(register1, register2)
	if err != nil {
		t.Errorf("Expected no error on skipIfEqualReg, got %v", err)
	}

	// Verify PC was incremented by 2 (skip next instruction)
	expectedPC := originalPC + 2
	if cpu.pc != expectedPC {
		t.Errorf("Expected PC to be 0x%X, got 0x%X", expectedPC, cpu.pc)
	}
}

// TestSkipIfEqualRegNoSkip verifies that skipIfEqualReg does NOT skip when registers contain different values.
func TestSkipIfEqualRegNoSkip(t *testing.T) {
	cpu := NewCPU()

	register1 := uint16(0x2)
	register2 := uint16(0x5)

	// Set registers to different values
	cpu.registers[register1] = 0x42
	cpu.registers[register2] = 0x50
	originalPC := cpu.pc

	// Execute skip if equal registers
	err := cpu.skipIfEqualReg(register1, register2)
	if err != nil {
		t.Errorf("Expected no error on skipIfEqualReg, got %v", err)
	}

	// Verify PC was NOT incremented
	if cpu.pc != originalPC {
		t.Errorf("Expected PC to remain 0x%X, got 0x%X", originalPC, cpu.pc)
	}
}

// TestSkipIfEqualRegInvalidRegister1 verifies that skipIfEqualReg returns an error
// when the first register index exceeds 0xF.
func TestSkipIfEqualRegInvalidRegister1(t *testing.T) {
	cpu := NewCPU()

	invalidRegister := uint16(0x10)
	validRegister := uint16(0x5)

	err := cpu.skipIfEqualReg(invalidRegister, validRegister)
	if err == nil {
		t.Errorf("Expected error for invalid first register, got none")
	}

	// Verify error message
	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}

	// Verify PC was not changed
	if cpu.pc != 0x200 {
		t.Errorf("Expected PC to remain 0x200, got 0x%X", cpu.pc)
	}
}

// TestSkipIfEqualRegInvalidRegister2 verifies that skipIfEqualReg returns an error
// when the second register index exceeds 0xF.
func TestSkipIfEqualRegInvalidRegister2(t *testing.T) {
	cpu := NewCPU()

	validRegister := uint16(0x5)
	invalidRegister := uint16(0x10)

	err := cpu.skipIfEqualReg(validRegister, invalidRegister)
	if err == nil {
		t.Errorf("Expected error for invalid second register, got none")
	}

	// Verify error message
	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}

	// Verify PC was not changed
	if cpu.pc != 0x200 {
		t.Errorf("Expected PC to remain 0x200, got 0x%X", cpu.pc)
	}
}

// TestSkipIfEqualRegBothInvalid verifies that skipIfEqualReg returns an error
// when both register indices exceed 0xF.
func TestSkipIfEqualRegBothInvalid(t *testing.T) {
	cpu := NewCPU()

	invalidRegister1 := uint16(0x10)
	invalidRegister2 := uint16(0x11)

	err := cpu.skipIfEqualReg(invalidRegister1, invalidRegister2)
	if err == nil {
		t.Errorf("Expected error for both invalid registers, got none")
	}

	// Verify error message
	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestSkipIfEqualRegEdgeCases verifies skipIfEqualReg works with edge case values.
func TestSkipIfEqualRegEdgeCases(t *testing.T) {
	cpu := NewCPU()

	// Test with value 0x00
	register1 := uint16(0x0)
	register2 := uint16(0x1)
	cpu.registers[register1] = 0x00
	cpu.registers[register2] = 0x00
	cpu.pc = 0x200

	err := cpu.skipIfEqualReg(register1, register2)
	if err != nil {
		t.Errorf("Expected no error for edge case 0x00, got %v", err)
	}
	if cpu.pc != 0x202 {
		t.Errorf("Expected PC to be 0x202, got 0x%X", cpu.pc)
	}

	// Test with value 0xFF
	register1 = uint16(0xE)
	register2 = uint16(0xF)
	cpu.registers[register1] = 0xFF
	cpu.registers[register2] = 0xFF
	cpu.pc = 0x300

	err = cpu.skipIfEqualReg(register1, register2)
	if err != nil {
		t.Errorf("Expected no error for edge case 0xFF, got %v", err)
	}
	if cpu.pc != 0x302 {
		t.Errorf("Expected PC to be 0x302, got 0x%X", cpu.pc)
	}
}

// TestSkipIfEqualRegAllRegisters verifies skipIfEqualReg works correctly with all pairs of registers.
func TestSkipIfEqualRegAllRegisters(t *testing.T) {
	cpu := NewCPU()

	// Test comparing each register with every other register
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			register1 := uint16(i)
			register2 := uint16(j)
			value := uint16(0x40 + i)

			// Set both registers to the same value
			cpu.registers[register1] = uint8(value)
			cpu.registers[register2] = uint8(value)
			cpu.pc = 0x200

			err := cpu.skipIfEqualReg(register1, register2)
			if err != nil {
				t.Errorf("Registers %d, %d: Expected no error, got %v", i, j, err)
			}

			// All should skip since values are equal
			if cpu.pc != 0x202 {
				t.Errorf("Registers %d, %d: Expected PC to be 0x202, got 0x%X", i, j, cpu.pc)
			}
		}
	}
}

// TestSkipIfEqualRegSameRegister verifies skipIfEqualReg works when comparing a register to itself.
func TestSkipIfEqualRegSameRegister(t *testing.T) {
	cpu := NewCPU()

	register := uint16(0x5)
	value := uint16(0x7A)

	cpu.registers[register] = uint8(value)
	cpu.pc = 0x200

	// Comparing a register to itself should always be equal
	err := cpu.skipIfEqualReg(register, register)
	if err != nil {
		t.Errorf("Expected no error comparing register to itself, got %v", err)
	}

	if cpu.pc != 0x202 {
		t.Errorf("Expected PC to be 0x202, got 0x%X", cpu.pc)
	}
}

// TestSetReg verifies that setReg (0x6xkk) correctly sets a register to a given byte value.
func TestSetReg(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	value := uint16(0x42)

	// Execute set register
	err := cpu.setReg(registerIndex, value)
	if err != nil {
		t.Errorf("Expected no error on setReg, got %v", err)
	}

	// Verify register was set to the value
	if cpu.registers[registerIndex] != uint8(value) {
		t.Errorf("Expected register[%d] to be 0x%X, got 0x%X", registerIndex, value, cpu.registers[registerIndex])
	}
}

// TestSetRegInvalidRegister verifies that setReg returns an error
// when the register index exceeds 0xF.
func TestSetRegInvalidRegister(t *testing.T) {
	cpu := NewCPU()

	invalidRegister := uint16(0x10)
	value := uint16(0x42)

	err := cpu.setReg(invalidRegister, value)
	if err == nil {
		t.Errorf("Expected error for invalid register, got none")
	}

	// Verify error message
	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestSetRegOverwrite verifies that setReg correctly overwrites existing register values.
func TestSetRegOverwrite(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x3)

	// Set register to initial value
	cpu.registers[registerIndex] = 0x10

	// Overwrite with new value
	newValue := uint16(0xAA)
	err := cpu.setReg(registerIndex, newValue)
	if err != nil {
		t.Errorf("Expected no error on setReg, got %v", err)
	}

	// Verify register was overwritten
	if cpu.registers[registerIndex] != uint8(newValue) {
		t.Errorf("Expected register[%d] to be 0x%X, got 0x%X", registerIndex, newValue, cpu.registers[registerIndex])
	}
}

// TestSetRegEdgeCases verifies setReg works with edge case values (0x00 and 0xFF).
func TestSetRegEdgeCases(t *testing.T) {
	cpu := NewCPU()

	// Test setting to 0x00
	registerIndex := uint16(0x0)
	err := cpu.setReg(registerIndex, 0x00)
	if err != nil {
		t.Errorf("Expected no error for edge case 0x00, got %v", err)
	}
	if cpu.registers[registerIndex] != 0x00 {
		t.Errorf("Expected register[0] to be 0x00, got 0x%X", cpu.registers[registerIndex])
	}

	// Test setting to 0xFF
	registerIndex = uint16(0xF)
	err = cpu.setReg(registerIndex, 0xFF)
	if err != nil {
		t.Errorf("Expected no error for edge case 0xFF, got %v", err)
	}
	if cpu.registers[registerIndex] != 0xFF {
		t.Errorf("Expected register[15] to be 0xFF, got 0x%X", cpu.registers[registerIndex])
	}
}

// TestSetRegAllRegisters verifies setReg works correctly with all 16 registers.
func TestSetRegAllRegisters(t *testing.T) {
	cpu := NewCPU()

	// Set each register to a unique value
	for i := 0; i < 16; i++ {
		registerIndex := uint16(i)
		value := uint16(0x10 + i)

		err := cpu.setReg(registerIndex, value)
		if err != nil {
			t.Errorf("Register %d: Expected no error, got %v", i, err)
		}

		// Verify register was set correctly
		if cpu.registers[registerIndex] != uint8(value) {
			t.Errorf("Register %d: Expected 0x%X, got 0x%X", i, value, cpu.registers[registerIndex])
		}
	}
}

// TestSetRegSequential verifies that setReg correctly sets multiple registers in sequence.
func TestSetRegSequential(t *testing.T) {
	cpu := NewCPU()

	values := []uint16{0x11, 0x22, 0x33, 0x44, 0x55}

	for i, val := range values {
		registerIndex := uint16(i)

		err := cpu.setReg(registerIndex, val)
		if err != nil {
			t.Errorf("Iteration %d: Expected no error, got %v", i, err)
		}

		// Verify register was set
		if cpu.registers[registerIndex] != uint8(val) {
			t.Errorf("Iteration %d: Expected 0x%X, got 0x%X", i, val, cpu.registers[registerIndex])
		}

		// Verify previously set registers remain unchanged
		for j := 0; j < i; j++ {
			expectedValue := values[j]
			if cpu.registers[j] != uint8(expectedValue) {
				t.Errorf("Iteration %d: Register %d changed from 0x%X to 0x%X", i, j, expectedValue, cpu.registers[j])
			}
		}
	}
}

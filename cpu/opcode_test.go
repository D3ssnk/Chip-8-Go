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

// TestAddVal verifies that addVal (0x7xkk) correctly adds a byte value to a register.
func TestAddVal(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	cpu.registers[registerIndex] = 0x10
	valueToAdd := uint16(0x20)

	// Execute add value
	err := cpu.addVal(registerIndex, valueToAdd)
	if err != nil {
		t.Errorf("Expected no error on addVal, got %v", err)
	}

	// Verify register was incremented correctly
	expectedValue := uint8(0x10 + 0x20)
	if cpu.registers[registerIndex] != expectedValue {
		t.Errorf("Expected register[%d] to be 0x%X, got 0x%X", registerIndex, expectedValue, cpu.registers[registerIndex])
	}
}

// TestAddValInvalidRegister verifies that addVal returns an error
// when the register index exceeds 0xF.
func TestAddValInvalidRegister(t *testing.T) {
	cpu := NewCPU()

	invalidRegister := uint16(0x10)
	value := uint16(0x42)

	err := cpu.addVal(invalidRegister, value)
	if err == nil {
		t.Errorf("Expected error for invalid register, got none")
	}

	// Verify error message
	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestAddValOverflow verifies that addVal handles overflow correctly (wraps around).
func TestAddValOverflow(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x3)

	// Test overflow: 0xFF + 0x02 = 0x01 (with wrapping)
	cpu.registers[registerIndex] = 0xFF
	err := cpu.addVal(registerIndex, 0x02)
	if err != nil {
		t.Errorf("Expected no error on overflow, got %v", err)
	}

	// In Go, uint8 overflow wraps around
	expectedValue := uint8(0x01)
	if cpu.registers[registerIndex] != expectedValue {
		t.Errorf("Expected register[%d] to be 0x%X (wrapped), got 0x%X", registerIndex, expectedValue, cpu.registers[registerIndex])
	}
}

// TestAddValZero verifies that adding zero doesn't change the register value.
func TestAddValZero(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x7)
	initialValue := uint8(0x42)
	cpu.registers[registerIndex] = initialValue

	// Add zero
	err := cpu.addVal(registerIndex, 0x00)
	if err != nil {
		t.Errorf("Expected no error when adding zero, got %v", err)
	}

	// Verify register value unchanged
	if cpu.registers[registerIndex] != initialValue {
		t.Errorf("Expected register[%d] to remain 0x%X, got 0x%X", registerIndex, initialValue, cpu.registers[registerIndex])
	}
}

// TestAddValEdgeCases verifies addVal works with edge case values.
func TestAddValEdgeCases(t *testing.T) {
	cpu := NewCPU()

	// Test adding to a zero register
	registerIndex := uint16(0x0)
	cpu.registers[registerIndex] = 0x00
	err := cpu.addVal(registerIndex, 0xFF)
	if err != nil {
		t.Errorf("Expected no error for edge case, got %v", err)
	}
	if cpu.registers[registerIndex] != 0xFF {
		t.Errorf("Expected register[0] to be 0xFF, got 0x%X", cpu.registers[registerIndex])
	}

	// Test adding max value to max value (overflow)
	registerIndex = uint16(0xF)
	cpu.registers[registerIndex] = 0xFF
	err = cpu.addVal(registerIndex, 0xFF)
	if err != nil {
		t.Errorf("Expected no error for edge case, got %v", err)
	}
	expectedValue := uint8(0xFE) // 0xFF + 0xFF wraps to 0xFE
	if cpu.registers[registerIndex] != expectedValue {
		t.Errorf("Expected register[15] to be 0x%X, got 0x%X", expectedValue, cpu.registers[registerIndex])
	}
}

// TestAddValAllRegisters verifies addVal works correctly with all 16 registers.
func TestAddValAllRegisters(t *testing.T) {
	cpu := NewCPU()

	// Add different values to each register
	for i := 0; i < 16; i++ {
		registerIndex := uint16(i)
		baseValue := uint16(0x10 + i)
		addValue := uint16(0x20)

		cpu.registers[registerIndex] = uint8(baseValue)

		err := cpu.addVal(registerIndex, addValue)
		if err != nil {
			t.Errorf("Register %d: Expected no error, got %v", i, err)
		}

		// Verify register was incremented correctly
		expectedValue := uint8(baseValue + addValue)
		if cpu.registers[registerIndex] != expectedValue {
			t.Errorf("Register %d: Expected 0x%X, got 0x%X", i, expectedValue, cpu.registers[registerIndex])
		}
	}
}

// TestAddValSequential verifies that multiple additions accumulate correctly.
func TestAddValSequential(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x4)
	cpu.registers[registerIndex] = 0x10

	// Add values sequentially
	addValues := []uint16{0x05, 0x10, 0x20, 0x30}
	expectedResult := uint8(0x10)

	for i, val := range addValues {
		err := cpu.addVal(registerIndex, val)
		if err != nil {
			t.Errorf("Addition %d: Expected no error, got %v", i, err)
		}

		expectedResult += uint8(val)
		if cpu.registers[registerIndex] != expectedResult {
			t.Errorf("Addition %d: Expected 0x%X, got 0x%X", i, expectedResult, cpu.registers[registerIndex])
		}
	}
}

// TestAddValMultipleRegisters verifies that addVal correctly updates individual registers
// without affecting others.
func TestAddValMultipleRegisters(t *testing.T) {
	cpu := NewCPU()

	// Initialize all registers with distinct values
	for i := 0; i < 16; i++ {
		cpu.registers[i] = uint8(i * 10)
	}

	// Add to register 5
	err := cpu.addVal(5, 0x50)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify register 5 was updated
	expectedValue5 := uint8(5*10 + 0x50)
	if cpu.registers[5] != expectedValue5 {
		t.Errorf("Expected register[5] to be 0x%X, got 0x%X", expectedValue5, cpu.registers[5])
	}

	// Verify other registers remain unchanged
	for i := 0; i < 16; i++ {
		if i == 5 {
			continue // Skip the one we modified
		}
		expectedValue := uint8(i * 10)
		if cpu.registers[i] != expectedValue {
			t.Errorf("Register %d: Expected 0x%X, got 0x%X", i, expectedValue, cpu.registers[i])
		}
	}
}

// TestSetRegReg verifies that setRegReg (0x8xy0) correctly copies a value
// from one register to another.
func TestSetRegReg(t *testing.T) {
	cpu := NewCPU()

	sourceRegister := uint16(0x3)
	destRegister := uint16(0x7)

	cpu.registers[sourceRegister] = 0x42
	cpu.registers[destRegister] = 0x00

	// Execute set register from register
	err := cpu.setRegReg(destRegister, sourceRegister)
	if err != nil {
		t.Errorf("Expected no error on setRegReg, got %v", err)
	}

	// Verify destination register was set to source value
	if cpu.registers[destRegister] != cpu.registers[sourceRegister] {
		t.Errorf("Expected register[%d] to be 0x%X, got 0x%X", destRegister, cpu.registers[sourceRegister], cpu.registers[destRegister])
	}

	// Verify source register was not modified
	if cpu.registers[sourceRegister] != 0x42 {
		t.Errorf("Expected source register[%d] to remain 0x42, got 0x%X", sourceRegister, cpu.registers[sourceRegister])
	}
}

// TestSetRegRegInvalidRegister1 verifies that setRegReg returns an error
// when the destination register index exceeds 0xF.
func TestSetRegRegInvalidRegister1(t *testing.T) {
	cpu := NewCPU()

	invalidRegister := uint16(0x10)
	sourceRegister := uint16(0x5)

	cpu.registers[sourceRegister] = 0x42

	err := cpu.setRegReg(invalidRegister, sourceRegister)
	if err == nil {
		t.Errorf("Expected error for invalid destination register, got none")
	}

	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestSetRegRegInvalidRegister2 verifies that setRegReg returns an error
// when the source register index exceeds 0xF.
func TestSetRegRegInvalidRegister2(t *testing.T) {
	cpu := NewCPU()

	destRegister := uint16(0x5)
	invalidRegister := uint16(0x10)

	err := cpu.setRegReg(destRegister, invalidRegister)
	if err == nil {
		t.Errorf("Expected error for invalid source register, got none")
	}

	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestSetRegRegBothInvalid verifies that setRegReg returns an error
// when both register indices exceed 0xF.
func TestSetRegRegBothInvalid(t *testing.T) {
	cpu := NewCPU()

	invalidReg1 := uint16(0x10)
	invalidReg2 := uint16(0x11)

	err := cpu.setRegReg(invalidReg1, invalidReg2)
	if err == nil {
		t.Errorf("Expected error for both invalid registers, got none")
	}

	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestSetRegRegEdgeCases verifies setRegReg works with edge case values.
func TestSetRegRegEdgeCases(t *testing.T) {
	cpu := NewCPU()

	// Test copying 0x00
	cpu.registers[0x0] = 0x00
	cpu.registers[0x1] = 0xFF
	err := cpu.setRegReg(0x1, 0x0)
	if err != nil {
		t.Errorf("Expected no error for edge case, got %v", err)
	}
	if cpu.registers[0x1] != 0x00 {
		t.Errorf("Expected register[1] to be 0x00, got 0x%X", cpu.registers[0x1])
	}

	// Test copying 0xFF
	cpu.registers[0xE] = 0xFF
	cpu.registers[0xF] = 0x00
	err = cpu.setRegReg(0xF, 0xE)
	if err != nil {
		t.Errorf("Expected no error for edge case, got %v", err)
	}
	if cpu.registers[0xF] != 0xFF {
		t.Errorf("Expected register[15] to be 0xFF, got 0x%X", cpu.registers[0xF])
	}
}

// TestSetRegRegAllRegisters verifies setRegReg works correctly with all register pairs.
func TestSetRegRegAllRegisters(t *testing.T) {
	cpu := NewCPU()

	// Initialize all registers with distinct values
	for i := 0; i < 16; i++ {
		cpu.registers[i] = uint8(i * 17) // Use 17 to create distinct values
	}

	// Test copying from each register to every other register
	for src := 0; src < 16; src++ {
		for dst := 0; dst < 16; dst++ {
			originalDestValue := cpu.registers[dst]
			sourceValue := cpu.registers[src]

			err := cpu.setRegReg(uint16(dst), uint16(src))
			if err != nil {
				t.Errorf("Copy src[%d] to dst[%d]: Expected no error, got %v", src, dst, err)
			}

			// Verify destination was updated to source value
			if cpu.registers[dst] != sourceValue {
				t.Errorf("Copy src[%d] to dst[%d]: Expected 0x%X, got 0x%X", src, dst, sourceValue, cpu.registers[dst])
			}

			// Verify source register was not modified
			if cpu.registers[src] != sourceValue {
				t.Errorf("Copy src[%d] to dst[%d]: Source register was modified to 0x%X", src, dst, cpu.registers[src])
			}

			// Restore destination for next test
			cpu.registers[dst] = originalDestValue
		}
	}
}

// TestSetRegRegSameRegister verifies that setRegReg works when source and destination are the same register.
func TestSetRegRegSameRegister(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	cpu.registers[registerIndex] = 0x42

	err := cpu.setRegReg(registerIndex, registerIndex)
	if err != nil {
		t.Errorf("Expected no error when source and dest are same, got %v", err)
	}

	// Register value should remain unchanged
	if cpu.registers[registerIndex] != 0x42 {
		t.Errorf("Expected register[%d] to remain 0x42, got 0x%X", registerIndex, cpu.registers[registerIndex])
	}
}

// TestSetRegRegIsolation verifies that copying to one register doesn't affect other registers.
func TestSetRegRegIsolation(t *testing.T) {
	cpu := NewCPU()

	// Initialize all registers with distinct values
	for i := 0; i < 16; i++ {
		cpu.registers[i] = uint8(i * 10)
	}

	// Copy register 3 to register 7
	sourceValue := cpu.registers[3]
	err := cpu.setRegReg(7, 3)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify register 7 was updated
	if cpu.registers[7] != sourceValue {
		t.Errorf("Expected register[7] to be 0x%X, got 0x%X", sourceValue, cpu.registers[7])
	}

	// Verify all other registers remain unchanged
	for i := 0; i < 16; i++ {
		if i == 7 {
			continue // Skip the one we modified
		}
		expectedValue := uint8(i * 10)
		if cpu.registers[i] != expectedValue {
			t.Errorf("Register %d: Expected 0x%X, got 0x%X", i, expectedValue, cpu.registers[i])
		}
	}
}

// TestSetRegRegSequential verifies that multiple copy operations work correctly in sequence.
func TestSetRegRegSequential(t *testing.T) {
	cpu := NewCPU()

	cpu.registers[0] = 0x11
	cpu.registers[1] = 0x22
	cpu.registers[2] = 0x33

	// Copy 0 to 5
	err := cpu.setRegReg(5, 0)
	if err != nil || cpu.registers[5] != 0x11 {
		t.Errorf("First copy failed: expected 0x11, got 0x%X", cpu.registers[5])
	}

	// Copy 1 to 6
	err = cpu.setRegReg(6, 1)
	if err != nil || cpu.registers[6] != 0x22 {
		t.Errorf("Second copy failed: expected 0x22, got 0x%X", cpu.registers[6])
	}

	// Copy 2 to 7
	err = cpu.setRegReg(7, 2)
	if err != nil || cpu.registers[7] != 0x33 {
		t.Errorf("Third copy failed: expected 0x33, got 0x%X", cpu.registers[7])
	}

	// Verify original values unchanged
	if cpu.registers[0] != 0x11 || cpu.registers[1] != 0x22 || cpu.registers[2] != 0x33 {
		t.Errorf("Original registers were modified during sequential copies")
	}
}

// TestOrReg verifies that orReg (0x8xy1) correctly performs a bitwise OR
// between two registers and stores the result in the first register.
func TestOrReg(t *testing.T) {
	cpu := NewCPU()

	register1 := uint16(0x3)
	register2 := uint16(0x7)

	cpu.registers[register1] = 0x0F // 0000 1111
	cpu.registers[register2] = 0xF0 // 1111 0000

	// Execute OR register
	err := cpu.orReg(register1, register2)
	if err != nil {
		t.Errorf("Expected no error on orReg, got %v", err)
	}

	// Verify destination register was set to bitwise OR result
	expectedValue := uint8(0xFF) // 1111 1111
	if cpu.registers[register1] != expectedValue {
		t.Errorf("Expected register[%d] to be 0x%X, got 0x%X", register1, expectedValue, cpu.registers[register1])
	}

	// Verify source register was not modified
	if cpu.registers[register2] != 0xF0 {
		t.Errorf("Expected source register[%d] to remain 0xF0, got 0x%X", register2, cpu.registers[register2])
	}
}

// TestOrRegInvalidRegister1 verifies that orReg returns an error
// when the first register index exceeds 0xF.
func TestOrRegInvalidRegister1(t *testing.T) {
	cpu := NewCPU()

	invalidRegister := uint16(0x10)
	validRegister := uint16(0x5)

	err := cpu.orReg(invalidRegister, validRegister)
	if err == nil {
		t.Errorf("Expected error for invalid first register, got none")
	}

	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestOrRegInvalidRegister2 verifies that orReg returns an error
// when the second register index exceeds 0xF.
func TestOrRegInvalidRegister2(t *testing.T) {
	cpu := NewCPU()

	validRegister := uint16(0x5)
	invalidRegister := uint16(0x10)

	err := cpu.orReg(validRegister, invalidRegister)
	if err == nil {
		t.Errorf("Expected error for invalid second register, got none")
	}

	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestOrRegBothInvalid verifies that orReg returns an error
// when both register indices exceed 0xF.
func TestOrRegBothInvalid(t *testing.T) {
	cpu := NewCPU()

	invalidReg1 := uint16(0x10)
	invalidReg2 := uint16(0x11)

	err := cpu.orReg(invalidReg1, invalidReg2)
	if err == nil {
		t.Errorf("Expected error for both invalid registers, got none")
	}

	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestOrRegEdgeCases verifies orReg works correctly with bitwise edge case values.
func TestOrRegEdgeCases(t *testing.T) {
	cpu := NewCPU()

	// Test 0x00 | 0x00 = 0x00
	cpu.registers[0x0] = 0x00
	cpu.registers[0x1] = 0x00
	err := cpu.orReg(0x0, 0x1)
	if err != nil {
		t.Errorf("Expected no error for edge case 0x00 | 0x00, got %v", err)
	}
	if cpu.registers[0x0] != 0x00 {
		t.Errorf("Expected register[0] to be 0x00, got 0x%X", cpu.registers[0x0])
	}

	// Test 0xFF | 0x00 = 0xFF
	cpu.registers[0x2] = 0xFF
	cpu.registers[0x3] = 0x00
	err = cpu.orReg(0x2, 0x3)
	if err != nil {
		t.Errorf("Expected no error for edge case 0xFF | 0x00, got %v", err)
	}
	if cpu.registers[0x2] != 0xFF {
		t.Errorf("Expected register[2] to be 0xFF, got 0x%X", cpu.registers[0x2])
	}

	// Test alternating bits: 0xAA | 0x55 = 0xFF
	// 0xAA = 10101010, 0x55 = 01010101
	cpu.registers[0x4] = 0xAA
	cpu.registers[0x5] = 0x55
	err = cpu.orReg(0x4, 0x5)
	if err != nil {
		t.Errorf("Expected no error for edge case 0xAA | 0x55, got %v", err)
	}
	if cpu.registers[0x4] != 0xFF {
		t.Errorf("Expected register[4] to be 0xFF, got 0x%X", cpu.registers[0x4])
	}
}

// TestOrRegSameRegister verifies that orReg works when source and destination are the same register.
func TestOrRegSameRegister(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	value := uint8(0x42)
	cpu.registers[registerIndex] = value

	err := cpu.orReg(registerIndex, registerIndex)
	if err != nil {
		t.Errorf("Expected no error when source and dest are same, got %v", err)
	}

	// A value OR'd with itself should remain unchanged
	if cpu.registers[registerIndex] != value {
		t.Errorf("Expected register[%d] to remain 0x%X, got 0x%X", registerIndex, value, cpu.registers[registerIndex])
	}
}

// TestOrRegIsolation verifies that performing an OR operation on one register doesn't affect others.
func TestOrRegIsolation(t *testing.T) {
	cpu := NewCPU()

	// Initialize all registers with distinct values
	for i := 0; i < 16; i++ {
		cpu.registers[i] = uint8(i * 10)
	}

	// Override specific registers for the OR test
	cpu.registers[3] = 0x01 // 0000 0001
	cpu.registers[7] = 0x02 // 0000 0010

	err := cpu.orReg(3, 7)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify register 3 was updated (0x01 | 0x02 = 0x03)
	if cpu.registers[3] != 0x03 {
		t.Errorf("Expected register[3] to be 0x03, got 0x%X", cpu.registers[3])
	}

	// Verify all other registers remain unchanged
	for i := 0; i < 16; i++ {
		if i == 3 {
			continue // Skip the destination register
		}
		
		expectedValue := uint8(i * 10)
		if i == 7 {
			expectedValue = 0x02 // Source register maintains its specific value
		}
		
		if cpu.registers[i] != expectedValue {
			t.Errorf("Register %d: Expected 0x%X, got 0x%X", i, expectedValue, cpu.registers[i])
		}
	}
}

// TestAndReg verifies that andReg (0x8xy2) correctly performs a bitwise AND
// between two registers and stores the result in the first register.
func TestAndReg(t *testing.T) {
	cpu := NewCPU()

	register1 := uint16(0x3)
	register2 := uint16(0x7)

	cpu.registers[register1] = 0xAF // 1010 1111
	cpu.registers[register2] = 0xF0 // 1111 0000

	// Execute AND register
	err := cpu.andReg(register1, register2)
	if err != nil {
		t.Errorf("Expected no error on andReg, got %v", err)
	}

	// Verify destination register was set to bitwise AND result
	expectedValue := uint8(0xA0) // 1010 0000
	if cpu.registers[register1] != expectedValue {
		t.Errorf("Expected register[%d] to be 0x%X, got 0x%X", register1, expectedValue, cpu.registers[register1])
	}

	// Verify source register was not modified
	if cpu.registers[register2] != 0xF0 {
		t.Errorf("Expected source register[%d] to remain 0xF0, got 0x%X", register2, cpu.registers[register2])
	}
}

// TestAndRegInvalidRegister1 verifies that andReg returns an error
// when the first register index exceeds 0xF.
func TestAndRegInvalidRegister1(t *testing.T) {
	cpu := NewCPU()

	invalidRegister := uint16(0x10)
	validRegister := uint16(0x5)

	err := cpu.andReg(invalidRegister, validRegister)
	if err == nil {
		t.Errorf("Expected error for invalid first register, got none")
	}

	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestAndRegInvalidRegister2 verifies that andReg returns an error
// when the second register index exceeds 0xF.
func TestAndRegInvalidRegister2(t *testing.T) {
	cpu := NewCPU()

	validRegister := uint16(0x5)
	invalidRegister := uint16(0x10)

	err := cpu.andReg(validRegister, invalidRegister)
	if err == nil {
		t.Errorf("Expected error for invalid second register, got none")
	}

	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestAndRegBothInvalid verifies that andReg returns an error
// when both register indices exceed 0xF.
func TestAndRegBothInvalid(t *testing.T) {
	cpu := NewCPU()

	invalidReg1 := uint16(0x10)
	invalidReg2 := uint16(0x11)

	err := cpu.andReg(invalidReg1, invalidReg2)
	if err == nil {
		t.Errorf("Expected error for both invalid registers, got none")
	}

	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestAndRegEdgeCases verifies andReg works correctly with bitwise edge case values.
func TestAndRegEdgeCases(t *testing.T) {
	cpu := NewCPU()

	// Test 0x00 & 0x00 = 0x00
	cpu.registers[0x0] = 0x00
	cpu.registers[0x1] = 0x00
	err := cpu.andReg(0x0, 0x1)
	if err != nil {
		t.Errorf("Expected no error for edge case 0x00 & 0x00, got %v", err)
	}
	if cpu.registers[0x0] != 0x00 {
		t.Errorf("Expected register[0] to be 0x00, got 0x%X", cpu.registers[0x0])
	}

	// Test 0xFF & 0x00 = 0x00
	cpu.registers[0x2] = 0xFF
	cpu.registers[0x3] = 0x00
	err = cpu.andReg(0x2, 0x3)
	if err != nil {
		t.Errorf("Expected no error for edge case 0xFF & 0x00, got %v", err)
	}
	if cpu.registers[0x2] != 0x00 {
		t.Errorf("Expected register[2] to be 0x00, got 0x%X", cpu.registers[0x2])
	}

	// Test alternating bits: 0xAA & 0x55 = 0x00
	// 0xAA = 10101010, 0x55 = 01010101
	cpu.registers[0x4] = 0xAA
	cpu.registers[0x5] = 0x55
	err = cpu.andReg(0x4, 0x5)
	if err != nil {
		t.Errorf("Expected no error for edge case 0xAA & 0x55, got %v", err)
	}
	if cpu.registers[0x4] != 0x00 {
		t.Errorf("Expected register[4] to be 0x00, got 0x%X", cpu.registers[0x4])
	}
}

// TestAndRegSameRegister verifies that andReg works when source and destination are the same register.
func TestAndRegSameRegister(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	value := uint8(0x42)
	cpu.registers[registerIndex] = value

	err := cpu.andReg(registerIndex, registerIndex)
	if err != nil {
		t.Errorf("Expected no error when source and dest are same, got %v", err)
	}

	// A value AND'd with itself should remain unchanged
	if cpu.registers[registerIndex] != value {
		t.Errorf("Expected register[%d] to remain 0x%X, got 0x%X", registerIndex, value, cpu.registers[registerIndex])
	}
}

// TestAndRegIsolation verifies that performing an AND operation on one register doesn't affect others.
func TestAndRegIsolation(t *testing.T) {
	cpu := NewCPU()

	// Initialize all registers with distinct values
	for i := 0; i < 16; i++ {
		cpu.registers[i] = uint8(i * 10)
	}

	// Override specific registers for the AND test
	cpu.registers[3] = 0x0F // 0000 1111
	cpu.registers[7] = 0x55 // 0101 0101

	err := cpu.andReg(3, 7)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify register 3 was updated (0x0F & 0x55 = 0x05)
	if cpu.registers[3] != 0x05 {
		t.Errorf("Expected register[3] to be 0x05, got 0x%X", cpu.registers[3])
	}

	// Verify all other registers remain unchanged
	for i := 0; i < 16; i++ {
		if i == 3 {
			continue // Skip the destination register
		}
		
		expectedValue := uint8(i * 10)
		if i == 7 {
			expectedValue = 0x55 // Source register maintains its specific value
		}
		
		if cpu.registers[i] != expectedValue {
			t.Errorf("Register %d: Expected 0x%X, got 0x%X", i, expectedValue, cpu.registers[i])
		}
	}
}
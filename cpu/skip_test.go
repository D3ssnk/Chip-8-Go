// Package cpu implements tests for CHIP-8 opcodes.
package cpu

import (
	"testing"
)

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

// TestSkipIfNotEqualRegSkip verifies that skipIfNotEqualReg (0x9xy0) correctly skips
// the next instruction when two registers contain different values.
func TestSkipIfNotEqualRegSkip(t *testing.T) {
	cpu := NewCPU()

	register1 := uint16(0x3)
	register2 := uint16(0x7)

	// Set registers to different values
	cpu.registers[register1] = 0x42
	cpu.registers[register2] = 0x50
	cpu.pc = 0x200
	originalPC := cpu.pc

	err := cpu.skipIfNotEqualReg(register1, register2)
	if err != nil {
		t.Errorf("Expected no error on skipIfNotEqualReg, got %v", err)
	}

	// Verify PC was incremented by 2 (skip next instruction)
	expectedPC := originalPC + 2
	if cpu.pc != expectedPC {
		t.Errorf("Expected PC to be 0x%X, got 0x%X", expectedPC, cpu.pc)
	}
}

// TestSkipIfNotEqualRegNoSkip verifies that skipIfNotEqualReg does NOT skip
// when registers contain equal values.
func TestSkipIfNotEqualRegNoSkip(t *testing.T) {
	cpu := NewCPU()

	register1 := uint16(0x2)
	register2 := uint16(0x5)

	// Set registers to the exact same value
	cpu.registers[register1] = 0x42
	cpu.registers[register2] = 0x42
	cpu.pc = 0x200
	originalPC := cpu.pc

	err := cpu.skipIfNotEqualReg(register1, register2)
	if err != nil {
		t.Errorf("Expected no error on skipIfNotEqualReg, got %v", err)
	}

	// Verify PC was NOT incremented
	if cpu.pc != originalPC {
		t.Errorf("Expected PC to remain 0x%X, got 0x%X", originalPC, cpu.pc)
	}
}

// TestSkipIfNotEqualRegInvalidRegister1 verifies that skipIfNotEqualReg returns an error
// when the first register index exceeds 0xF.
func TestSkipIfNotEqualRegInvalidRegister1(t *testing.T) {
	cpu := NewCPU()

	invalidRegister := uint16(0x10)
	validRegister := uint16(0x5)

	err := cpu.skipIfNotEqualReg(invalidRegister, validRegister)
	if err == nil {
		t.Errorf("Expected error for invalid first register, got none")
	}

	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestSkipIfNotEqualRegInvalidRegister2 verifies that skipIfNotEqualReg returns an error
// when the second register index exceeds 0xF.
func TestSkipIfNotEqualRegInvalidRegister2(t *testing.T) {
	cpu := NewCPU()

	validRegister := uint16(0x5)
	invalidRegister := uint16(0x10)

	err := cpu.skipIfNotEqualReg(validRegister, invalidRegister)
	if err == nil {
		t.Errorf("Expected error for invalid second register, got none")
	}

	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestSkipIfNotEqualRegBothInvalid verifies that skipIfNotEqualReg returns an error
// when both register indices exceed 0xF.
func TestSkipIfNotEqualRegBothInvalid(t *testing.T) {
	cpu := NewCPU()

	invalidReg1 := uint16(0x10)
	invalidReg2 := uint16(0x11)

	err := cpu.skipIfNotEqualReg(invalidReg1, invalidReg2)
	if err == nil {
		t.Errorf("Expected error for both invalid registers, got none")
	}
}

// TestSkipIfNotEqualRegSameRegister verifies that skipIfNotEqualReg behaves correctly
// when comparing a register to itself (it should never skip).
func TestSkipIfNotEqualRegSameRegister(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	cpu.registers[registerIndex] = 0x7A
	cpu.pc = 0x200
	originalPC := cpu.pc

	// Comparing a register to itself means they are equal, so it should NOT skip
	err := cpu.skipIfNotEqualReg(registerIndex, registerIndex)
	if err != nil {
		t.Errorf("Expected no error comparing register to itself, got %v", err)
	}

	if cpu.pc != originalPC {
		t.Errorf("Expected PC to remain 0x%X, got 0x%X", originalPC, cpu.pc)
	}
}

// TestSkipIfKeyPressedSkip verifies that skipIfKeyPressed (0xEx9E) correctly skips
// the next instruction when the specified key is currently pressed.
func TestSkipIfKeyPressedSkip(t *testing.T) {
	cpu := NewCPU()

	keyIndex := uint16(0xA)

	// Simulate the key being pressed
	cpu.keypad[keyIndex] = true

	originalPC := cpu.pc

	err := cpu.skipIfKeyPressed(keyIndex)
	if err != nil {
		t.Errorf("Expected no error on skipIfKeyPressed, got %v", err)
	}

	// Verify PC was incremented by 2 (skip next instruction)
	expectedPC := originalPC + 2
	if cpu.pc != expectedPC {
		t.Errorf("Expected PC to be 0x%X, got 0x%X", expectedPC, cpu.pc)
	}
}

// TestSkipIfKeyPressedNoSkip verifies that skipIfKeyPressed does NOT skip
// when the specified key is not pressed.
func TestSkipIfKeyPressedNoSkip(t *testing.T) {
	cpu := NewCPU()

	keyIndex := uint16(0x5)

	// Explicitly ensure the key is NOT pressed
	cpu.keypad[keyIndex] = false

	originalPC := cpu.pc

	err := cpu.skipIfKeyPressed(keyIndex)
	if err != nil {
		t.Errorf("Expected no error on skipIfKeyPressed, got %v", err)
	}

	// Verify PC was NOT incremented
	if cpu.pc != originalPC {
		t.Errorf("Expected PC to remain 0x%X, got 0x%X", originalPC, cpu.pc)
	}
}

// TestSkipIfKeyPressedOutOfBounds verifies that skipIfKeyPressed returns an error
// when the requested key index exceeds 0xF (15).
func TestSkipIfKeyPressedOutOfBounds(t *testing.T) {
	cpu := NewCPU()

	originalPC := cpu.pc
	invalidKey := uint16(0x10) // 16 is out of bounds for a 16-key keypad (0-15)

	err := cpu.skipIfKeyPressed(invalidKey)
	if err == nil {
		t.Errorf("Expected error for out of bounds key, got none")
	}

	if err.Error() != "Key press is out of bounds" {
		t.Errorf("Expected error message 'Key press is out of bounds', got '%v'", err.Error())
	}

	// Verify PC was not changed
	if cpu.pc != originalPC {
		t.Errorf("Expected PC to remain 0x%X, got 0x%X", originalPC, cpu.pc)
	}
}

// TestSkipIfKeyPressedIsolation verifies that checking one key doesn't accidentally
// read the state of another key.
func TestSkipIfKeyPressedIsolation(t *testing.T) {
	cpu := NewCPU()

	// Press key 0x2
	cpu.keypad[0x2] = true
	// Ensure key 0x3 is NOT pressed
	cpu.keypad[0x3] = false

	originalPC := cpu.pc

	// Check key 0x3 (which is adjacent to the pressed key 0x2)
	err := cpu.skipIfKeyPressed(0x3)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// It should NOT skip
	if cpu.pc != originalPC {
		t.Errorf("Expected PC to remain 0x%X, got 0x%X", originalPC, cpu.pc)
	}
}

// TestSkipIfKeyNotPressedSkip verifies that skipIfKeyNotPressed (0xExA1) correctly skips
// the next instruction when the specified key is NOT currently pressed.
func TestSkipIfKeyNotPressedSkip(t *testing.T) {
	cpu := NewCPU()

	keyIndex := uint16(0xA)

	// Ensure the key is NOT pressed
	cpu.keypad[keyIndex] = false

	originalPC := cpu.pc

	err := cpu.skipIfKeyNotPressed(keyIndex)
	if err != nil {
		t.Errorf("Expected no error on skipIfKeyNotPressed, got %v", err)
	}

	// Verify PC was incremented by 2 (skip next instruction)
	expectedPC := originalPC + 2
	if cpu.pc != expectedPC {
		t.Errorf("Expected PC to be 0x%X, got 0x%X", expectedPC, cpu.pc)
	}
}

// TestSkipIfKeyNotPressedNoSkip verifies that skipIfKeyNotPressed does NOT skip
// when the specified key IS pressed.
func TestSkipIfKeyNotPressedNoSkip(t *testing.T) {
	cpu := NewCPU()

	keyIndex := uint16(0x5)

	// Simulate the key being pressed
	cpu.keypad[keyIndex] = true

	originalPC := cpu.pc

	err := cpu.skipIfKeyNotPressed(keyIndex)
	if err != nil {
		t.Errorf("Expected no error on skipIfKeyNotPressed, got %v", err)
	}

	// Verify PC was NOT incremented
	if cpu.pc != originalPC {
		t.Errorf("Expected PC to remain 0x%X, got 0x%X", originalPC, cpu.pc)
	}
}

// TestSkipIfKeyNotPressedOutOfBounds verifies that skipIfKeyNotPressed returns an error
// when the requested key index exceeds 0xF (15).
func TestSkipIfKeyNotPressedOutOfBounds(t *testing.T) {
	cpu := NewCPU()

	originalPC := cpu.pc
	invalidKey := uint16(0x10) // Out of bounds

	err := cpu.skipIfKeyNotPressed(invalidKey)
	if err == nil {
		t.Errorf("Expected error for out of bounds key, got none")
	}

	if err.Error() != "Key press is out of bounds" {
		t.Errorf("Expected error message 'Key press is out of bounds', got '%v'", err.Error())
	}

	// Verify PC was not changed
	if cpu.pc != originalPC {
		t.Errorf("Expected PC to remain 0x%X, got 0x%X", originalPC, cpu.pc)
	}
}

// TestSkipIfKeyNotPressedIsolation verifies that checking one key doesn't accidentally
// rely on the state of another key.
func TestSkipIfKeyNotPressedIsolation(t *testing.T) {
	cpu := NewCPU()

	// Press key 0x2
	cpu.keypad[0x2] = true
	// Ensure key 0x3 is NOT pressed
	cpu.keypad[0x3] = false

	originalPC := cpu.pc

	// Check key 0x2 (which IS pressed, so it should NOT skip)
	err := cpu.skipIfKeyNotPressed(0x2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// It should NOT skip
	if cpu.pc != originalPC {
		t.Errorf("Expected PC to remain 0x%X, got 0x%X", originalPC, cpu.pc)
	}
}

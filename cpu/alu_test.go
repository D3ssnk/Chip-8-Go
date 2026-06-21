// Package cpu implements tests for CHIP-8 opcodes.
package cpu

import (
	"testing"
)

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

// TestSubRegNoBorrow verifies that subReg (0x8xy5) correctly subtracts Vy from Vx
// when Vx > Vy, setting VF to 1 (no borrow).
func TestSubRegNoBorrow(t *testing.T) {
	cpu := NewCPU()

	register1 := uint16(0x3)
	register2 := uint16(0x7)

	cpu.registers[register1] = 0x10 // 16
	cpu.registers[register2] = 0x05 // 5

	err := cpu.subReg(register1, register2)
	if err != nil {
		t.Errorf("Expected no error on subReg, got %v", err)
	}

	// Verify math result
	if cpu.registers[register1] != 0x0B { // 11
		t.Errorf("Expected register[%d] to be 0x0B, got 0x%X", register1, cpu.registers[register1])
	}

	// Verify VF flag (No borrow = 1)
	if cpu.registers[0xF] != 1 {
		t.Errorf("Expected VF (register[15]) to be 1, got %d", cpu.registers[0xF])
	}

	// Verify source register unchanged
	if cpu.registers[register2] != 0x05 {
		t.Errorf("Expected source register[%d] to remain 0x05, got 0x%X", register2, cpu.registers[register2])
	}
}

// TestSubRegBorrow verifies that subReg correctly subtracts Vy from Vx
// when Vx < Vy, allowing underflow and setting VF to 0 (borrow).
func TestSubRegBorrow(t *testing.T) {
	cpu := NewCPU()

	register1 := uint16(0x3)
	register2 := uint16(0x7)

	cpu.registers[register1] = 0x05 // 5
	cpu.registers[register2] = 0x07 // 7

	err := cpu.subReg(register1, register2)
	if err != nil {
		t.Errorf("Expected no error on subReg, got %v", err)
	}

	// Verify underflow math result (5 - 7 = 254 in uint8)
	expectedValue := uint8(0xFE) 
	if cpu.registers[register1] != expectedValue {
		t.Errorf("Expected register[%d] to be 0xFE (underflow), got 0x%X", register1, cpu.registers[register1])
	}

	// Verify VF flag (Borrow = 0)
	if cpu.registers[0xF] != 0 {
		t.Errorf("Expected VF (register[15]) to be 0, got %d", cpu.registers[0xF])
	}
}

// TestSubRegEqual verifies subtraction when Vx == Vy, 
// expecting a result of 0 and VF set to 1 (no borrow).
func TestSubRegEqual(t *testing.T) {
	cpu := NewCPU()

	register1 := uint16(0x4)
	register2 := uint16(0x5)

	cpu.registers[register1] = 0x42
	cpu.registers[register2] = 0x42

	err := cpu.subReg(register1, register2)
	if err != nil {
		t.Errorf("Expected no error on subReg, got %v", err)
	}

	if cpu.registers[register1] != 0x00 {
		t.Errorf("Expected register[%d] to be 0x00, got 0x%X", register1, cpu.registers[register1])
	}

	if cpu.registers[0xF] != 1 {
		t.Errorf("Expected VF to be 1 when values are equal, got %d", cpu.registers[0xF])
	}
}

// TestSubRegInvalidRegister1 verifies that subReg returns an error
// when the first register index exceeds 0xF.
func TestSubRegInvalidRegister1(t *testing.T) {
	cpu := NewCPU()

	err := cpu.subReg(0x10, 0x5)
	if err == nil {
		t.Errorf("Expected error for invalid first register, got none")
	}
}

// TestSubRegInvalidRegister2 verifies that subReg returns an error
// when the second register index exceeds 0xF.
func TestSubRegInvalidRegister2(t *testing.T) {
	cpu := NewCPU()

	err := cpu.subReg(0x5, 0x10)
	if err == nil {
		t.Errorf("Expected error for invalid second register, got none")
	}
}

// TestSubRegBothInvalid verifies that subReg returns an error
// when both register indices exceed 0xF.
func TestSubRegBothInvalid(t *testing.T) {
	cpu := NewCPU()

	err := cpu.subReg(0x10, 0x11)
	if err == nil {
		t.Errorf("Expected error for both invalid registers, got none")
	}
}

// TestSubRegSameRegister verifies subtracting a register from itself.
func TestSubRegSameRegister(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	cpu.registers[registerIndex] = 0x42

	err := cpu.subReg(registerIndex, registerIndex)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if cpu.registers[registerIndex] != 0x00 {
		t.Errorf("Expected register[%d] to be 0x00, got 0x%X", registerIndex, cpu.registers[registerIndex])
	}

	if cpu.registers[0xF] != 1 {
		t.Errorf("Expected VF to be 1 when subtracting register from itself, got %d", cpu.registers[0xF])
	}
}

// TestSubRegWithVFAsDestination verifies the CHIP-8 quirk where if VF (0xF)
// is the destination register, the math result overwrites the flag calculation.
func TestSubRegWithVFAsDestination(t *testing.T) {
	cpu := NewCPU()
	
	// Set VF as the destination (Vx)
	destRegister := uint16(0xF)
	sourceRegister := uint16(0x2)
	
	cpu.registers[destRegister] = 0x10 // 16
	cpu.registers[sourceRegister] = 0x05 // 5
	
	err := cpu.subReg(destRegister, sourceRegister)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// VF should be overwritten by the mathematical result (0x10 - 0x05 = 0x0B)
	// rather than storing the 'no borrow' flag of 1.
	if cpu.registers[0xF] != 0x0B {
		t.Errorf("Expected VF to be overwritten by math result 0x0B, got 0x%X", cpu.registers[0xF])
	}
}

// TestSubRegIsolation verifies that performing a SUB operation
// on one register doesn't affect others, except for VF.
func TestSubRegIsolation(t *testing.T) {
	cpu := NewCPU()

	// Initialize all registers with distinct values
	for i := 0; i < 15; i++ { // Skip F
		cpu.registers[i] = uint8(i * 10)
	}
	cpu.registers[0xF] = 0

	// Override specific registers
	cpu.registers[3] = 0x20 // 32
	cpu.registers[7] = 0x10 // 16

	err := cpu.subReg(3, 7)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all other registers remain unchanged
	for i := 0; i < 15; i++ {
		if i == 3 || i == 7 {
			continue 
		}
		
		expectedValue := uint8(i * 10)
		
		if cpu.registers[i] != expectedValue {
			t.Errorf("Register %d: Expected 0x%X, got 0x%X", i, expectedValue, cpu.registers[i])
		}
	}
}

// TestShiftRightLSBOne verifies that shiftRight (0x8xy6) correctly shifts Vx right by 1
// and sets VF to 1 when the least significant bit of Vx is 1.
func TestShiftRightLSBOne(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	cpu.registers[registerIndex] = 0x15 // 21 (Binary: 0001 0101)

	err := cpu.shiftRight(registerIndex)
	if err != nil {
		t.Errorf("Expected no error on shiftRight, got %v", err)
	}

	// Verify shifted result (0001 0101 >> 1 = 0000 1010 = 0x0A)
	if cpu.registers[registerIndex] != 0x0A {
		t.Errorf("Expected register[%d] to be 0x0A, got 0x%X", registerIndex, cpu.registers[registerIndex])
	}

	// Verify VF flag (LSB was 1)
	if cpu.registers[0xF] != 1 {
		t.Errorf("Expected VF (register[15]) to be 1, got %d", cpu.registers[0xF])
	}
}

// TestShiftRightLSBZero verifies that shiftRight correctly shifts Vx right by 1
// and sets VF to 0 when the least significant bit of Vx is 0.
func TestShiftRightLSBZero(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	cpu.registers[registerIndex] = 0x16 // 22 (Binary: 0001 0110)

	err := cpu.shiftRight(registerIndex)
	if err != nil {
		t.Errorf("Expected no error on shiftRight, got %v", err)
	}

	// Verify shifted result (0001 0110 >> 1 = 0000 1011 = 0x0B)
	if cpu.registers[registerIndex] != 0x0B {
		t.Errorf("Expected register[%d] to be 0x0B, got 0x%X", registerIndex, cpu.registers[registerIndex])
	}

	// Verify VF flag (LSB was 0)
	if cpu.registers[0xF] != 0 {
		t.Errorf("Expected VF (register[15]) to be 0, got %d", cpu.registers[0xF])
	}
}

// TestShiftRightInvalidRegister verifies that shiftRight returns an error
// when the register index exceeds 0xF.
func TestShiftRightInvalidRegister(t *testing.T) {
	cpu := NewCPU()

	err := cpu.shiftRight(0x10)
	if err == nil {
		t.Errorf("Expected error for invalid register, got none")
	}
	
	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestShiftRightWithVFAsDestination verifies that if VF (0xF) is the target register,
// the mathematical shifted result overwrites the LSB flag calculation.
func TestShiftRightWithVFAsDestination(t *testing.T) {
	cpu := NewCPU()
	
	destRegister := uint16(0xF)
	cpu.registers[destRegister] = 0x03 // Binary: 0000 0011 (LSB is 1)
	
	err := cpu.shiftRight(destRegister)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// VF should first be set to 1 (the LSB), but then immediately overwritten
	// by the shifted result (0x03 >> 1 = 0x01).
	// While the end result is coincidentally the same here (1), testing 0x05 ensures
	// we see the math result (2) rather than just the flag (1).
	
	cpu.registers[destRegister] = 0x05 // Binary: 0000 0101 (LSB is 1)
	cpu.shiftRight(destRegister)
	
	if cpu.registers[0xF] != 0x02 { // 0x05 >> 1 = 0x02
		t.Errorf("Expected VF to be overwritten by math result 0x02, got 0x%X", cpu.registers[0xF])
	}
}

// TestShiftRightIsolation verifies that performing a shift operation
// on one register doesn't affect others, except for VF.
func TestShiftRightIsolation(t *testing.T) {
	cpu := NewCPU()

	// Initialize all registers with distinct values
	for i := 0; i < 15; i++ { // Skip F
		cpu.registers[i] = uint8(i * 10)
	}
	cpu.registers[0xF] = 0

	// Override a specific register
	targetRegister := uint16(3)
	cpu.registers[targetRegister] = 0x08 // 0000 1000

	err := cpu.shiftRight(targetRegister)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all other registers remain unchanged
	for i := 0; i < 15; i++ {
		if i == int(targetRegister) {
			continue 
		}
		
		expectedValue := uint8(i * 10)
		
		if cpu.registers[i] != expectedValue {
			t.Errorf("Register %d: Expected 0x%X, got 0x%X", i, expectedValue, cpu.registers[i])
		}
	}
}

// TestSubRegReverseNoBorrow verifies that subRegReverse (0x8xy7) correctly subtracts Vx from Vy
// when Vy > Vx, setting VF to 1 (no borrow).
func TestSubRegReverseNoBorrow(t *testing.T) {
	cpu := NewCPU()

	register1 := uint16(0x3) // Vx
	register2 := uint16(0x7) // Vy

	cpu.registers[register1] = 0x05 // 5
	cpu.registers[register2] = 0x10 // 16

	err := cpu.subRegReverse(register1, register2)
	if err != nil {
		t.Errorf("Expected no error on subRegReverse, got %v", err)
	}

	// Verify math result (Vy - Vx = 16 - 5 = 11)
	if cpu.registers[register1] != 0x0B { 
		t.Errorf("Expected register[%d] to be 0x0B, got 0x%X", register1, cpu.registers[register1])
	}

	// Verify VF flag (No borrow = 1)
	if cpu.registers[0xF] != 1 {
		t.Errorf("Expected VF (register[15]) to be 1, got %d", cpu.registers[0xF])
	}

	// Verify source register unchanged
	if cpu.registers[register2] != 0x10 {
		t.Errorf("Expected source register[%d] to remain 0x10, got 0x%X", register2, cpu.registers[register2])
	}
}

// TestSubRegReverseBorrow verifies that subRegReverse correctly subtracts Vx from Vy
// when Vy < Vx, allowing underflow and setting VF to 0 (borrow).
func TestSubRegReverseBorrow(t *testing.T) {
	cpu := NewCPU()

	register1 := uint16(0x3) // Vx
	register2 := uint16(0x7) // Vy

	cpu.registers[register1] = 0x07 // 7
	cpu.registers[register2] = 0x05 // 5

	err := cpu.subRegReverse(register1, register2)
	if err != nil {
		t.Errorf("Expected no error on subRegReverse, got %v", err)
	}

	// Verify underflow math result (5 - 7 = 254 in uint8)
	expectedValue := uint8(0xFE) 
	if cpu.registers[register1] != expectedValue {
		t.Errorf("Expected register[%d] to be 0xFE (underflow), got 0x%X", register1, cpu.registers[register1])
	}

	// Verify VF flag (Borrow = 0)
	if cpu.registers[0xF] != 0 {
		t.Errorf("Expected VF (register[15]) to be 0, got %d", cpu.registers[0xF])
	}
}

// TestSubRegReverseEqual verifies subtraction when Vy == Vx, 
// expecting a result of 0 and VF set to 1 (no borrow).
func TestSubRegReverseEqual(t *testing.T) {
	cpu := NewCPU()

	register1 := uint16(0x4)
	register2 := uint16(0x5)

	cpu.registers[register1] = 0x42
	cpu.registers[register2] = 0x42

	err := cpu.subRegReverse(register1, register2)
	if err != nil {
		t.Errorf("Expected no error on subRegReverse, got %v", err)
	}

	if cpu.registers[register1] != 0x00 {
		t.Errorf("Expected register[%d] to be 0x00, got 0x%X", register1, cpu.registers[register1])
	}

	if cpu.registers[0xF] != 1 {
		t.Errorf("Expected VF to be 1 when values are equal, got %d", cpu.registers[0xF])
	}
}

// TestSubRegReverseInvalidRegister1 verifies that subRegReverse returns an error
// when the first register index exceeds 0xF.
func TestSubRegReverseInvalidRegister1(t *testing.T) {
	cpu := NewCPU()

	err := cpu.subRegReverse(0x10, 0x5)
	if err == nil {
		t.Errorf("Expected error for invalid first register, got none")
	}
}

// TestSubRegReverseInvalidRegister2 verifies that subRegReverse returns an error
// when the second register index exceeds 0xF.
func TestSubRegReverseInvalidRegister2(t *testing.T) {
	cpu := NewCPU()

	err := cpu.subRegReverse(0x5, 0x10)
	if err == nil {
		t.Errorf("Expected error for invalid second register, got none")
	}
}

// TestSubRegReverseBothInvalid verifies that subRegReverse returns an error
// when both register indices exceed 0xF.
func TestSubRegReverseBothInvalid(t *testing.T) {
	cpu := NewCPU()

	err := cpu.subRegReverse(0x10, 0x11)
	if err == nil {
		t.Errorf("Expected error for both invalid registers, got none")
	}
}

// TestSubRegReverseSameRegister verifies subtracting a register from itself.
func TestSubRegReverseSameRegister(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	cpu.registers[registerIndex] = 0x42

	err := cpu.subRegReverse(registerIndex, registerIndex)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if cpu.registers[registerIndex] != 0x00 {
		t.Errorf("Expected register[%d] to be 0x00, got 0x%X", registerIndex, cpu.registers[registerIndex])
	}

	if cpu.registers[0xF] != 1 {
		t.Errorf("Expected VF to be 1 when subtracting register from itself, got %d", cpu.registers[0xF])
	}
}

// TestSubRegReverseWithVFAsDestination verifies the CHIP-8 quirk where if VF (0xF)
// is the destination register, the math result overwrites the flag calculation.
func TestSubRegReverseWithVFAsDestination(t *testing.T) {
	cpu := NewCPU()
	
	// Set VF as the destination (Vx)
	destRegister := uint16(0xF)
	sourceRegister := uint16(0x2)
	
	cpu.registers[destRegister] = 0x05   // Vx (5)
	cpu.registers[sourceRegister] = 0x10 // Vy (16)
	
	err := cpu.subRegReverse(destRegister, sourceRegister)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// VF should evaluate the flag (16 >= 5, so flag = 1), 
	// but then immediately be overwritten by the math result (16 - 5 = 11 / 0x0B).
	if cpu.registers[0xF] != 0x0B {
		t.Errorf("Expected VF to be overwritten by math result 0x0B, got 0x%X", cpu.registers[0xF])
	}
}

// TestSubRegReverseIsolation verifies that performing a SUBN operation
// on one register doesn't affect others, except for VF.
func TestSubRegReverseIsolation(t *testing.T) {
	cpu := NewCPU()

	// Initialize all registers with distinct values
	for i := 0; i < 15; i++ { // Skip F
		cpu.registers[i] = uint8(i * 10)
	}
	cpu.registers[0xF] = 0

	// Override specific registers
	cpu.registers[3] = 0x10 // Vx (16)
	cpu.registers[7] = 0x20 // Vy (32)

	err := cpu.subRegReverse(3, 7)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all other registers remain unchanged
	for i := 0; i < 15; i++ {
		if i == 3 || i == 7 {
			continue 
		}
		
		expectedValue := uint8(i * 10)
		
		if cpu.registers[i] != expectedValue {
			t.Errorf("Register %d: Expected 0x%X, got 0x%X", i, expectedValue, cpu.registers[i])
		}
	}
}

// TestShiftLeftMSBOne verifies that shiftLeft (0x8xyE) correctly shifts Vx left by 1
// and sets VF to 1 when the most significant bit of Vx is 1.
func TestShiftLeftMSBOne(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	cpu.registers[registerIndex] = 0x85 // 133 (Binary: 1000 0101)

	err := cpu.shiftLeft(registerIndex)
	if err != nil {
		t.Errorf("Expected no error on shiftLeft, got %v", err)
	}

	// Verify shifted result (1000 0101 << 1 = 0000 1010 = 0x0A)
	if cpu.registers[registerIndex] != 0x0A {
		t.Errorf("Expected register[%d] to be 0x0A, got 0x%X", registerIndex, cpu.registers[registerIndex])
	}

	// Verify VF flag (MSB was 1)
	if cpu.registers[0xF] != 1 {
		t.Errorf("Expected VF (register[15]) to be 1, got %d", cpu.registers[0xF])
	}
}

// TestShiftLeftMSBZero verifies that shiftLeft correctly shifts Vx left by 1
// and sets VF to 0 when the most significant bit of Vx is 0.
func TestShiftLeftMSBZero(t *testing.T) {
	cpu := NewCPU()

	registerIndex := uint16(0x5)
	cpu.registers[registerIndex] = 0x45 // 69 (Binary: 0100 0101)

	err := cpu.shiftLeft(registerIndex)
	if err != nil {
		t.Errorf("Expected no error on shiftLeft, got %v", err)
	}

	// Verify shifted result (0100 0101 << 1 = 1000 1010 = 0x8A)
	if cpu.registers[registerIndex] != 0x8A {
		t.Errorf("Expected register[%d] to be 0x8A, got 0x%X", registerIndex, cpu.registers[registerIndex])
	}

	// Verify VF flag (MSB was 0)
	if cpu.registers[0xF] != 0 {
		t.Errorf("Expected VF (register[15]) to be 0, got %d", cpu.registers[0xF])
	}
}

// TestShiftLeftInvalidRegister verifies that shiftLeft returns an error
// when the register index exceeds 0xF.
func TestShiftLeftInvalidRegister(t *testing.T) {
	cpu := NewCPU()

	err := cpu.shiftLeft(0x10)
	if err == nil {
		t.Errorf("Expected error for invalid register, got none")
	}
	
	if err.Error() != "Invalid Register" {
		t.Errorf("Expected error message 'Invalid Register', got '%v'", err.Error())
	}
}

// TestShiftLeftWithVFAsDestination verifies that if VF (0xF) is the target register,
// the mathematical shifted result overwrites the MSB flag calculation.
func TestShiftLeftWithVFAsDestination(t *testing.T) {
	cpu := NewCPU()
	
	destRegister := uint16(0xF)
	cpu.registers[destRegister] = 0x82 // Binary: 1000 0010 (MSB is 1)
	
	err := cpu.shiftLeft(destRegister)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// VF should first be set to 1 (the MSB), but then immediately overwritten
	// by the shifted result (0x82 << 1 = 0x04).
	if cpu.registers[0xF] != 0x04 { 
		t.Errorf("Expected VF to be overwritten by math result 0x04, got 0x%X", cpu.registers[0xF])
	}
}

// TestShiftLeftIsolation verifies that performing a shift operation
// on one register doesn't affect others, except for VF.
func TestShiftLeftIsolation(t *testing.T) {
	cpu := NewCPU()

	// Initialize all registers with distinct values
	for i := 0; i < 15; i++ { // Skip F
		cpu.registers[i] = uint8(i * 10)
	}
	cpu.registers[0xF] = 0

	// Override a specific register
	targetRegister := uint16(3)
	cpu.registers[targetRegister] = 0x10 // 0001 0000

	err := cpu.shiftLeft(targetRegister)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all other registers remain unchanged
	for i := 0; i < 15; i++ {
		if i == int(targetRegister) {
			continue 
		}
		
		expectedValue := uint8(i * 10)
		
		if cpu.registers[i] != expectedValue {
			t.Errorf("Register %d: Expected 0x%X, got 0x%X", i, expectedValue, cpu.registers[i])
		}
	}
}
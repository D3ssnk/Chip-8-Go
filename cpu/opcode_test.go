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

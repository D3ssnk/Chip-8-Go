package cpu

import (
	"errors"
	"math/rand/v2"
)

func (cpu *CPU) clear() {
	cpu.display = [32][64]bool{}
}

func (cpu *CPU) ret() error {
	if cpu.sp == 0 {
		return errors.New("Stack is empty")
	}

	cpu.pc = cpu.stack[cpu.sp]
	cpu.sp--
	return nil
}

func (cpu *CPU) jump(address uint16) error {
	if address > 0xFFF {
		return errors.New("Address is out of bounds")
	}
	cpu.pc = address
	return nil
}

func (cpu *CPU) call(address uint16) error {
	if address > 0xFFF {
		return errors.New("Address is out of bounds")
	}
	if cpu.sp == 0xF {
		return errors.New("Stack is full")
	}

	cpu.sp++
	cpu.stack[cpu.sp] = cpu.pc
	cpu.pc = address
	return nil
}

func (cpu *CPU) skipIfEqual(registerIndex uint16, lastByte uint16) error {
	if registerIndex > 0xF {
		return errors.New("Invalid Register")
	}

	if cpu.registers[registerIndex] == uint8(lastByte) {
		cpu.pc += 2
	}

	return nil
}

func (cpu *CPU) skipIfNotEqual(registerIndex uint16, lastByte uint16) error {
	if registerIndex > 0xF {
		return errors.New("Invalid Register")
	}

	if cpu.registers[registerIndex] != uint8(lastByte) {
		cpu.pc += 2
	}

	return nil
}

func (cpu *CPU) skipIfEqualReg(register1Index uint16, register2Index uint16) error {
	if register1Index > 0xF || register2Index > 0xF {
		return errors.New("Invalid Register")
	}

	if cpu.registers[register1Index] == cpu.registers[register2Index] {
		cpu.pc += 2
	}

	return nil
}

func (cpu *CPU) setReg(registerIndex uint16, lastByte uint16) error {
	if registerIndex > 0xF {
		return errors.New("Invalid Register")
	}

	cpu.registers[registerIndex] = uint8(lastByte)

	return nil
}

func (cpu *CPU) addVal(registerIndex uint16, lastByte uint16) error {
	if registerIndex > 0xF {
		return errors.New("Invalid Register")
	}

	cpu.registers[registerIndex] += uint8(lastByte)

	return nil
}

func (cpu *CPU) setRegReg(register1Index uint16, register2Index uint16) error {
	if register1Index > 0xF || register2Index > 0xF {
		return errors.New("Invalid Register")
	}

	cpu.registers[register1Index] = cpu.registers[register2Index]

	return nil
}

func (cpu *CPU) orReg(register1Index uint16, register2Index uint16) error {
	if register1Index > 0xF || register2Index > 0xF {
		return errors.New("Invalid Register")
	}

	cpu.registers[register1Index] = cpu.registers[register2Index] | cpu.registers[register1Index]

	return nil
}

func (cpu *CPU) andReg(register1Index uint16, register2Index uint16) error {
	if register1Index > 0xF || register2Index > 0xF {
		return errors.New("Invalid Register")
	}

	cpu.registers[register1Index] = cpu.registers[register2Index] & cpu.registers[register1Index]

	return nil
}

func (cpu *CPU) subReg(register1Index uint16, register2Index uint16) error {
	if register1Index > 0xF || register2Index > 0xF {
		return errors.New("Invalid Register")
	}
	register1Value := cpu.registers[register1Index]
	register2Value := cpu.registers[register2Index]

	if cpu.registers[register1Index] >= cpu.registers[register2Index] {
		cpu.registers[0xF] = 1
	} else {
		cpu.registers[0xF] = 0
	}
	cpu.registers[register1Index] = register1Value - register2Value

	return nil
}

func (cpu *CPU) shiftRight(register1Index uint16) error {
	if register1Index > 0xF {
		return errors.New("Invalid Register")
	}
	register1Value := cpu.registers[register1Index]

	if register1Value&0x01 == 0x01 {
		cpu.registers[0xF] = 1
	} else {
		cpu.registers[0xF] = 0
	}
	cpu.registers[register1Index] = register1Value >> 1

	return nil
}

func (cpu *CPU) subRegReverse(register1Index uint16, register2Index uint16) error {
	if register1Index > 0xF || register2Index > 0xF {
		return errors.New("Invalid Register")
	}
	register1Value := cpu.registers[register1Index]
	register2Value := cpu.registers[register2Index]

	if cpu.registers[register2Index] >= cpu.registers[register1Index] {
		cpu.registers[0xF] = 1
	} else {
		cpu.registers[0xF] = 0
	}
	cpu.registers[register1Index] = register2Value - register1Value

	return nil
}

func (cpu *CPU) shiftLeft(register1Index uint16) error {
	if register1Index > 0xF {
		return errors.New("Invalid Register")
	}
	register1Value := cpu.registers[register1Index]

	if register1Value>>7 == 0x01 {
		cpu.registers[0xF] = 1
	} else {
		cpu.registers[0xF] = 0
	}
	cpu.registers[register1Index] = register1Value << 1

	return nil
}

func (cpu *CPU) setI(address uint16) error {
	if address > 0xFFF {
		return errors.New("Address is out of bounds")
	}

	cpu.i = address

	return nil
}

func (cpu *CPU) skipIfNotEqualReg(register1Index uint16, register2Index uint16) error {
	if register1Index > 0xF || register2Index > 0xF {
		return errors.New("Invalid Register")
	}

	if cpu.registers[register1Index] != cpu.registers[register2Index] {
		cpu.pc += 2
	}

	return nil
}

func (cpu *CPU) jumpV0(address uint16) error {
	if uint16(cpu.registers[0x0])+address > 0xFFF {
		return errors.New("Address is out of bounds")
	}

	cpu.pc = uint16(cpu.registers[0x0]) + address

	return nil
}

func (cpu *CPU) randReg(registerIndex uint16, lastByte uint16) error {
	if registerIndex > 0xF {
		return errors.New("Invalid Register")
	}
	randVal := rand.IntN(256)
	cpu.registers[registerIndex] = (uint8(lastByte) & uint8(randVal))

	return nil
}

func (cpu *CPU) draw(xcord uint16, ycord uint16, counter uint16) error {
	if xcord >= 0x40 || ycord >= 0x20 {
		return errors.New("Cordinates out of bounds")
	}
	if cpu.i+counter > 0xFFF {
		return errors.New("Address out of bounds")
	}

	cpu.registers[0xF] = 0
	for i, currentByte := range cpu.memory[cpu.i : cpu.i+counter] {
		for j := 0; j < 8; j++ {
			currentPixelX := (int(xcord) + j) % 64
			currentPixelY := (int(ycord) + i) % 32
			if (currentByte>>(7-j))&0x01 == 0x01 {
				if cpu.display[currentPixelY][currentPixelX] {
					cpu.registers[0xF] = 1
				}
				cpu.display[currentPixelY][currentPixelX] = !cpu.display[currentPixelY][currentPixelX]
			}
		}
	}
	return nil
}

func (cpu *CPU) skipIfKeyPressed(key uint16) error {
	if key > 0xF {
		return errors.New("Key press is out of bounds")
	}

	if cpu.keypad[key] {
		cpu.pc += 2
	}

	return nil
}

func (cpu *CPU) skipIfKeyNotPressed(key uint16) error {
	if key > 0xF {
		return errors.New("Key press is out of bounds")
	}

	if !cpu.keypad[key] {
		cpu.pc += 2
	}

	return nil
}

func (cpu *CPU) getDelayTimer(registerIndex uint16) error {
	if registerIndex > 0xF {
		return errors.New("Invalid Register")
	}
	
	cpu.registers[registerIndex] = cpu.delayTimer

	return nil
}
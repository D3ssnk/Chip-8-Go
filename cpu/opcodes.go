package cpu

import "errors"

func (cpu *CPU) clear() {
	cpu.display = [32][8]uint8{}
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
	if register1Index  > 0xF {
		return errors.New("Invalid Register")
	}
	register1Value := cpu.registers[register1Index]

	if register1Value & 0x01 == 0x01 {
		cpu.registers[0xF] = 1
	} else {
		cpu.registers[0xF] = 0
	}
	cpu.registers[register1Index] = register1Value >> 1

	return nil
}
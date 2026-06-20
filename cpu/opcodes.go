package cpu

import "errors"

func (cpu *CPU) clear() {
	cpu.display = [32][8]uint8{}
}

func (cpu *CPU) ret() error{
	if cpu.sp == 0 {
		return errors.New("Stack is empty")
	}

	cpu.pc = cpu.stack[cpu.sp]
	cpu.sp--
	return nil
}

func (cpu *CPU) jump(address uint16) error{
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
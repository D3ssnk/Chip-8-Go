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
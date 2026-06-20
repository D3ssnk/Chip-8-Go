package cpu

func (cpu *CPU) clear() {
	cpu.display = [32][8]uint8{}
}
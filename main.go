package main

import (
	"github.com/D3ssnk/Chip-8-Go/cpu"
)

func main() {
	var cpu = cpu.NewCPU()
	cpu.LoadROM("Maze (alt) [David Winter, 199x].ch8")
}

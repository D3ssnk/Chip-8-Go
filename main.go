package main

import (
	"fmt"
	"github.com/D3ssnk/Chip-8-Go/cpu"
)

func main() {
	var cpu = cpu.NewCPU()
	fmt.Println(cpu.Memory)
}

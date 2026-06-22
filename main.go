package main

import (
	"log"
	"os"

	"github.com/D3ssnk/Chip-8-Go/cpu"
	"github.com/D3ssnk/Chip-8-Go/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	var cpu = cpu.NewCPU()
	cpu.LoadROM(os.Args[1])
	var game = game.InitGame(cpu)
	ebiten.SetWindowSize(640, 320)
	ebiten.SetWindowTitle("Chip 8")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

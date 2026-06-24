package game

import (
	"image/color"

	"github.com/D3ssnk/Chip-8-Go/cpu"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game defines the Ebiten application structure, holding an instance of the CHIP-8 CPU.
type Game struct {
	cpu cpu.CPU
}

// checkKeypad maps standard Ebiten keyboard inputs to the classic CHIP-8 16-key keypad layout.
// It resets the CPU keypad array on each call, then iterates through mapped keys (1-4, Q-R, A-F, Z-V),
// updating the CPU's internal keypad state if any are pressed.
// Returns true if at least one mapped key is currently depressed.
func (g *Game) checkKeypad() bool {
	g.cpu.ResetKeypad()
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.cpu.SetKeypad(0, true)
	}  
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.cpu.SetKeypad(1, true)
	}  
	if ebiten.IsKeyPressed(ebiten.Key2) {
		g.cpu.SetKeypad(2, true)
	} 
	if ebiten.IsKeyPressed(ebiten.Key3) {
		g.cpu.SetKeypad(3, true)
	}  
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.cpu.SetKeypad(4, true)
	}  
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.cpu.SetKeypad(5, true)
	}  
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.cpu.SetKeypad(6, true)
	} 
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.cpu.SetKeypad(7, true)
	} 
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.cpu.SetKeypad(8, true)
	}  
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.cpu.SetKeypad(9, true)
	}  
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		g.cpu.SetKeypad(10, true)
	}  
	if ebiten.IsKeyPressed(ebiten.KeyC) {
		g.cpu.SetKeypad(11, true)
	}  
	if ebiten.IsKeyPressed(ebiten.Key4) {
		g.cpu.SetKeypad(12, true)
	}  
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.cpu.SetKeypad(13, true)
	}  
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.cpu.SetKeypad(14, true)
	}  
	if ebiten.IsKeyPressed(ebiten.KeyV) {
		g.cpu.SetKeypad(15, true)
	}

	if g.cpu.GetKeypad() == [16]bool{} {
		return false
	}
	return true
}

// Update acts as the central game loop, executing 10 CPU instructions per frame to simulate an appropriate clock speed.
// It halts execution if the CPU is waiting for user input and handles the 60Hz decrementing of delay and sound timers.
func (g *Game) Update() error {
	for i := 0; i < 10; i++ {
		keyWasPressed := g.checkKeypad()
		if keyWasPressed && g.cpu.GeWaitingForKeypad() {
			g.cpu.SetWaitingForKeypad(false)
		}
		if g.cpu.GeWaitingForKeypad() {
			break
		}

		instruction, err := g.cpu.Fetch()
		if err != nil {
			return err
		}

		err = g.cpu.Execute(instruction)
		if err != nil {
			return err
		}
	}

	if g.cpu.GetDelayTimer() > 0 {
		g.cpu.DecrementDelayTimer()
	}

	if g.cpu.GetSoundTimer() > 0 {
		g.cpu.DecrementSoundTimer()
	}

	return nil
}

// Draw renders the internal CHIP-8 64x32 boolean display matrix to the Ebiten screen window.
// It iterates through the display array, drawing 10x10 white pixels for active true bits,
// and 10x10 pink pixels for inactive false bits, translating the geometry matrix accordingly.
func (g *Game) Draw(screen *ebiten.Image) {
	var pinkPixel *ebiten.Image = ebiten.NewImage(10, 10)
	var whitePixel *ebiten.Image = ebiten.NewImage(10, 10)
	var options ebiten.DrawImageOptions
	display := g.cpu.GetDisplay()

	pinkPixel.Fill(color.RGBA{255, 105, 180, 255})
	whitePixel.Fill(color.RGBA{255, 255, 255, 255})

	for i := 0; i < 32; i++ {
		for j := 0; j < 64; j++ {
			if display[i][j] {
				screen.DrawImage(whitePixel, &options)
			} else {
				screen.DrawImage(pinkPixel, &options)
			}
			options.GeoM.Translate(10, 0)
		}
		options.GeoM.Translate(-640, 10)
	}
}

// Layout dictates the rendering resolution of the application window.
// Given the CHIP-8's native 64x32 resolution and the 10x10 pixel scaling implemented in Draw,
// the layout dimensions are hardcoded to 640x320.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 320
}

// InitGame acts as a constructor, injecting an instantiated CPU object into the newly created Game struct.
func InitGame(cpu cpu.CPU) Game {
	var game Game
	game.cpu = cpu

	return game
}

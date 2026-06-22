package game

import (
	//"image/color"

	"image/color"
	"github.com/D3ssnk/Chip-8-Go/cpu"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	cpu cpu.CPU
}

func (g *Game) checkKeypad() bool{
	g.cpu.ResetKeypad()
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.cpu.SetKeypad(0,true)
	} else if ebiten.IsKeyPressed(ebiten.Key1) {
		g.cpu.SetKeypad(1,true)
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		g.cpu.SetKeypad(2,true)
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		g.cpu.SetKeypad(3,true)
	} else if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.cpu.SetKeypad(4,true)
	} else if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.cpu.SetKeypad(5,true)
	} else if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.cpu.SetKeypad(6,true)
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.cpu.SetKeypad(7,true)
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.cpu.SetKeypad(8,true)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.cpu.SetKeypad(9,true)
	} else if ebiten.IsKeyPressed(ebiten.KeyZ) {
		g.cpu.SetKeypad(10,true)
	} else if ebiten.IsKeyPressed(ebiten.KeyC) {
		g.cpu.SetKeypad(11,true)
	} else if ebiten.IsKeyPressed(ebiten.Key4) {
		g.cpu.SetKeypad(12,true)
	} else if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.cpu.SetKeypad(13,true)
	} else if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.cpu.SetKeypad(14,true)
	} else if ebiten.IsKeyPressed(ebiten.KeyV) {
		g.cpu.SetKeypad(15,true)
	}

	if g.cpu.GetKeypad() == [16]bool{} {
		return false
	}
	return true
}

func (g *Game) Update() error {
	keyWasPressed := g.checkKeypad()
	if keyWasPressed && g.cpu.GeWaitingForKeypad() {
		g.cpu.SetWaitingForKeypad(false)
	}

	if g.cpu.GetDelayTimer() > 0 {
		g.cpu.DecrementDelayTimer()
	}

	if g.cpu.GetSoundTimer() > 0 {
		g.cpu.DecrementSoundTimer()
	}

	instruction, err := g.cpu.Fetch()
	if err != nil {
		return err
	}

	err = g.cpu.Execute(instruction)
	if err != nil {
		return err
	}
	
	return nil
}

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
		options.GeoM.Translate(-640,10)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 320
}

func InitGame(cpu cpu.CPU) Game {
	var game Game
	game.cpu = cpu

	return game
}

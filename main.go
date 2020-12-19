package main

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Canvas struct {
	image *image.RGBA
}

func (c *Canvas) Update() error {
	const l = screenWidth * screenHeight
	for i := 0; i < l; i++ {
		x := 0
		c.image.Pix[4*i] = uint8(x >> 24)
		c.image.Pix[4*i+1] = uint8(x >> 16)
		c.image.Pix[4*i+2] = uint8(x >> 8)
		c.image.Pix[4*i+3] = 0xff
	}
	return nil
}

func (c *Canvas) Draw(screen *ebiten.Image) {
	screen.ReplacePixels(c.image.Pix)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (c *Canvas) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Noise (Ebiten Demo)")
	c := &Canvas{
		image: image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	}
	if err := ebiten.RunGame(c); err != nil {
		log.Fatal(err)
	}
}

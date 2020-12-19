package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Vertex struct {
	x, y float32
}

func randomVertex() Vertex {
	x, y := rand.Intn(screenWidth), rand.Intn(screenHeight)
	return Vertex{float32(x - screenWidth/2), float32(screenHeight/2 - y)}
}

type Canvas struct {
	image *image.RGBA
}

func (c *Canvas) PutPixel(x, y float32) {
	bounds := c.image.Bounds()
	cX := bounds.Max.X/2 + int(x)
	cY := bounds.Max.Y/2 - int(y)
	c.image.Set(cX, cY, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
}

func (c *Canvas) DrawLine(p0, p1 Vertex) {
	dy, dx := p1.y-p0.y, p1.x-p0.x
	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		if p0.x > p1.x {
			p0, p1 = p1, p0
		}
		a := dy / dx
		y := p0.y
		for x := p0.x; x <= p1.x; x++ {
			c.PutPixel(x, y)
			y += a
		}
	} else {
		if p0.y > p1.y {
			p0, p1 = p1, p0
		}
		a := dx / dy
		x := p0.x
		for y := p0.y; y <= p1.y; y++ {
			c.PutPixel(x, y)
			x += a
		}
	}
}

func (c *Canvas) Update() error {
	c.Clear()
	c.DrawLine(randomVertex(), randomVertex())
	return nil
}

func (c *Canvas) Clear() error {
	bounds := c.image.Bounds()
	for i := 0; i < bounds.Max.X; i++ {
		for j := 0; j < bounds.Max.Y; j++ {
			c.image.Set(i, j, color.RGBA{0, 0, 0, 0xFF})
		}
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
	ebiten.SetMaxTPS(10) // Slow enough to visualise the random lines
	c := &Canvas{
		image: image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	}
	if err := ebiten.RunGame(c); err != nil {
		log.Fatal(err)
	}
}

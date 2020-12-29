package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"rasterizer/canvas"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type game struct {
	canv    canvas.Canvas
	polygon []canvas.Vertex
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Rasterizer")
	ebiten.SetMaxTPS(60) // Slow enough to visualise the random lines
	g := game{canv: *canvas.NewCanvas(screenWidth, screenHeight)}
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}

func randomVertex(width, height int) canvas.Vertex {
	x, y := rand.Intn(width), rand.Intn(height)
	return canvas.Vertex{X: float32(x - width/2), Y: float32(height/2 - y)}
}

// Update does a thing
func (g *game) Update() error {
	if g.polygon == nil {
		width, height := g.canv.Dimensions()
		p0 := randomVertex(width, height)
		p1 := randomVertex(width, height)
		p2 := randomVertex(width, height)
		g.polygon = []canvas.Vertex{p0, p1, p2}
	}
	g.canv.Fill(color.RGBA{0, 0, 0, 0xFF})
	g.canv.ShadeTriangle(g.polygon[0], g.polygon[1], g.polygon[2], color.RGBA{170, 240, 209, 0xFF})
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.ReplacePixels(g.canv.Buffer())
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.canv.Dimensions()
}

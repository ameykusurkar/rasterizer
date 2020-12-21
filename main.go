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
	image   *image.RGBA
	polygon []Vertex
}

func (c *Canvas) PutPixel(x, y float32, color color.Color) {
	bounds := c.image.Bounds()
	cX := bounds.Max.X/2 + int(x)
	cY := bounds.Max.Y/2 - int(y)
	c.image.Set(cX, cY, color)
}

// TODO: Is this the most efficient way?
func Interpolate(p0, p1 Vertex) []Vertex {
	dy, dx := p1.y-p0.y, p1.x-p0.x
	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		return InterpolateHorizontal(p0, p1)
	}
  return InterpolateVertical(p0, p1)
}

func InterpolateHorizontal(p0, p1 Vertex) []Vertex {
	verts := make([]Vertex, 0)
	dy, dx := p1.y-p0.y, p1.x-p0.x
	if p0.x > p1.x {
		p0, p1 = p1, p0
	}
	a := dy / dx
	y := p0.y
	for x := p0.x; x <= p1.x; x++ {
		verts = append(verts, Vertex{x, y})
		y += a
	}
	return verts
}

func InterpolateVertical(p0, p1 Vertex) []Vertex {
	verts := make([]Vertex, 0)
	dy, dx := p1.y-p0.y, p1.x-p0.x
	if p0.y > p1.y {
		p0, p1 = p1, p0
	}
	a := dx / dy
	x := p0.x
	for y := p0.y; y <= p1.y; y++ {
		verts = append(verts, Vertex{x, y})
		x += a
	}
	return verts
}

func (c *Canvas) DrawLine(p0, p1 Vertex) {
	for _, vert := range Interpolate(p0, p1) {
		c.PutPixel(vert.x, vert.y, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
	}
}

func (c *Canvas) DrawTriangle(p0, p1, p2 Vertex) {
	c.DrawLine(p0, p1)
	c.DrawLine(p1, p2)
	c.DrawLine(p2, p0)
}

func (c *Canvas) FillTriangle(p0, p1, p2 Vertex) {
	if p1.y < p0.y {
		p0, p1 = p1, p0
	}
	if p2.y < p0.y {
		p0, p2 = p2, p0
	}
	if p2.y < p1.y {
		p1, p2 = p2, p1
	}

	x01 := InterpolateVertical(p0, p1)
	x01 = x01[:len(x01)-1] // Last value overlaps with x12
	x12 := InterpolateVertical(p1, p2)
	x02 := InterpolateVertical(p0, p2)
	x012 := append(x01, x12...)

	var xLeft, xRight []Vertex
	if x01[len(x01)-1].x < x02[len(x01)-1].x {
		xLeft, xRight = x012, x02
	} else {
		xLeft, xRight = x02, x012
	}

	// TODO: Understand why we cannot use DrawLine here
	for i := 0; i < len(x02); i++ {
		left, right := xLeft[i].x, xRight[i].x
		for x := left; x <= right; x++ {
			c.PutPixel(x, xLeft[i].y, color.RGBA{170, 240, 209, 0xFF})
		}
	}
}

func (c *Canvas) Update() error {
	if c.polygon == nil {
		p0, p1, p2 := randomVertex(), randomVertex(), randomVertex()
		c.polygon = []Vertex{p0, p1, p2}
	}
	c.Clear()
	c.DrawTriangle(c.polygon[0], c.polygon[1], c.polygon[2])
	c.FillTriangle(c.polygon[0], c.polygon[1], c.polygon[2])
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
	ebiten.SetWindowTitle("Rasterizer")
	ebiten.SetMaxTPS(10) // Slow enough to visualise the random lines
	c := &Canvas{
		image: image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	}
	if err := ebiten.RunGame(c); err != nil {
		log.Fatal(err)
	}
}

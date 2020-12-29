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

type vertex struct {
	x, y float32
}

func randomVertex() vertex {
	x, y := rand.Intn(screenWidth), rand.Intn(screenHeight)
	return vertex{float32(x - screenWidth/2), float32(screenHeight/2 - y)}
}

type canvas struct {
	image   *image.RGBA
	polygon []vertex
}

func (c *canvas) PutPixel(x, y float32, color color.Color) {
	bounds := c.image.Bounds()
	cX := bounds.Max.X/2 + int(x)
	cY := bounds.Max.Y/2 - int(y)
	c.image.Set(cX, cY, color)
}

// TODO: Is this the most efficient way?
func interpolate(p0, p1 vertex) []vertex {
	dy, dx := p1.y-p0.y, p1.x-p0.x
	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		return interpolateHorizontal(p0, p1)
	}
	return interpolateVertical(p0, p1)
}

func interpolateHorizontal(p0, p1 vertex) []vertex {
	verts := make([]vertex, 0)
	dy, dx := p1.y-p0.y, p1.x-p0.x
	if p0.x > p1.x {
		p0, p1 = p1, p0
	}
	a := dy / dx
	y := p0.y
	for x := p0.x; x <= p1.x; x++ {
		verts = append(verts, vertex{x, y})
		y += a
	}
	return verts
}

func interpolateVertical(p0, p1 vertex) []vertex {
	verts := make([]vertex, 0)
	dy, dx := p1.y-p0.y, p1.x-p0.x
	if p0.y > p1.y {
		p0, p1 = p1, p0
	}
	a := dx / dy
	x := p0.x
	for y := p0.y; y <= p1.y; y++ {
		verts = append(verts, vertex{x, y})
		x += a
	}
	return verts
}

func (c *canvas) DrawLine(p0, p1 vertex) {
	for _, vert := range interpolate(p0, p1) {
		c.PutPixel(vert.x, vert.y, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
	}
}

func (c *canvas) DrawTriangle(p0, p1, p2 vertex) {
	c.DrawLine(p0, p1)
	c.DrawLine(p1, p2)
	c.DrawLine(p2, p0)
}

func (c *canvas) FillTriangle(p0, p1, p2 vertex) {
	if p1.y < p0.y {
		p0, p1 = p1, p0
	}
	if p2.y < p0.y {
		p0, p2 = p2, p0
	}
	if p2.y < p1.y {
		p1, p2 = p2, p1
	}

	x01 := interpolateVertical(p0, p1)
	x01 = x01[:len(x01)-1] // Last value overlaps with x12
	x12 := interpolateVertical(p1, p2)
	x02 := interpolateVertical(p0, p2)
	x012 := append(x01, x12...)

	var xLeft, xRight []vertex
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

// TODO: Refactor
func (c *canvas) ShadeTriangle(p0, p1, p2 vertex) {
	if p1.y < p0.y {
		p0, p1 = p1, p0
	}
	if p2.y < p0.y {
		p0, p2 = p2, p0
	}
	if p2.y < p1.y {
		p1, p2 = p2, p1
	}

	h0, h1, h2 := float32(0.0), float32(0.5), float32(1.0)

	x01 := interpolateVertical(p0, p1)
	h01 := interpolateVertical(vertex{x: h0, y: p0.y}, vertex{x: h1, y: p1.y})
	x01 = x01[:len(x01)-1] // Last value overlaps with x12
	h01 = h01[:len(h01)-1] // Last value overlaps with x12

	x12 := interpolateVertical(p1, p2)
	h12 := interpolateVertical(vertex{x: h1, y: p1.y}, vertex{x: h2, y: p2.y})

	x02 := interpolateVertical(p0, p2)
	h02 := interpolateVertical(vertex{x: h0, y: p0.y}, vertex{x: h2, y: p2.y})

	x012 := append(x01, x12...)
	h012 := append(h01, h12...)

	var xLefts, xRights []vertex
	var hLefts, hRights []vertex
	if x01[len(x01)-1].x < x02[len(x01)-1].x {
		xLefts, xRights = x012, x02
		hLefts, hRights = h012, h02
	} else {
		xLefts, xRights = x02, x012
		hLefts, hRights = h02, h012
	}

	for i := 0; i < len(x02); i++ {
		xLeft, xRight := xLefts[i].x, xRights[i].x
		hLeft, hRight := hLefts[i].x, hRights[i].x
		hh := interpolateHorizontal(vertex{x: xLeft, y: hLeft}, vertex{x: xRight, y: hRight})
		for x := xLeft; x <= xRight; x++ {
			hGrad := hh[int(x-xLeft)].y
			c.PutPixel(x, xLefts[i].y, color.RGBA{uint8(hGrad * 170.0), uint8(hGrad * 240.0), uint8(hGrad * 209.0), 0xFF})
		}
	}
}

func (c *canvas) Update() error {
	if c.polygon == nil {
		p0, p1, p2 := randomVertex(), randomVertex(), randomVertex()
		c.polygon = []vertex{p0, p1, p2}
	}
	c.Clear()
	// c.DrawTriangle(c.polygon[0], c.polygon[1], c.polygon[2])
	c.ShadeTriangle(c.polygon[0], c.polygon[1], c.polygon[2])
	return nil
}

func (c *canvas) Clear() error {
	bounds := c.image.Bounds()
	for i := 0; i < bounds.Max.X; i++ {
		for j := 0; j < bounds.Max.Y; j++ {
			c.image.Set(i, j, color.RGBA{0, 0, 0, 0xFF})
		}
	}
	return nil
}

func (c *canvas) Draw(screen *ebiten.Image) {
	screen.ReplacePixels(c.image.Pix)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (c *canvas) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Rasterizer")
	ebiten.SetMaxTPS(10) // Slow enough to visualise the random lines
	c := &canvas{
		image: image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	}
	if err := ebiten.RunGame(c); err != nil {
		log.Fatal(err)
	}
}

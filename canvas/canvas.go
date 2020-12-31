package canvas

import (
	"image"
	"image/color"
	"math"
)

// Point is a thing
type Point struct {
	X, Y float32
}

// Canvas is a buffer on which we can draw lines, triangles etc.
type Canvas struct {
	image *image.RGBA
}

// NewCanvas returns a new Canvas with dimensions (width, height)
func NewCanvas(width, height int) *Canvas {
	return &Canvas{
		image: image.NewRGBA(image.Rect(0, 0, width, height)),
	}
}

// Dimensions returns the width and height of the Canvas.
func (c *Canvas) Dimensions() (int, int) {
	bounds := c.image.Bounds()
	return bounds.Max.X, bounds.Max.Y
}

// Buffer returns the raw RGBA buffer of the Canvas.
func (c *Canvas) Buffer() []uint8 {
	return c.image.Pix
}

// Fill sets all the pixels on the Canvas to the specified color.
func (c *Canvas) Fill(clr color.Color) error {
	bounds := c.image.Bounds()
	for i := 0; i < bounds.Max.X; i++ {
		for j := 0; j < bounds.Max.Y; j++ {
			c.image.Set(i, j, clr)
		}
	}
	return nil
}

// PutPixel puts at pixel at (x, y) on the Canvas, with (0, 0) as the top-left corner.
func (c *Canvas) PutPixel(x, y float32, color color.Color) {
	c.image.Set(int(x), int(y), color)
}

// DrawLine draws a line with a specified color between two points.
func (c *Canvas) DrawLine(p0, p1 Point, clr color.Color) {
	for _, vert := range interpolate(p0, p1) {
		c.PutPixel(vert.X, vert.Y, clr)
	}
}

// FillTriangle fills the triangle formed by the given three points with the specified color.
func (c *Canvas) FillTriangle(p0, p1, p2 Point, clr color.Color) {
	if p1.Y < p0.Y {
		p0, p1 = p1, p0
	}
	if p2.Y < p0.Y {
		p0, p2 = p2, p0
	}
	if p2.Y < p1.Y {
		p1, p2 = p2, p1
	}

	x01 := interpolateVertical(p0, p1)
	x01 = x01[:len(x01)-1] // Last value overlaps with x12
	x12 := interpolateVertical(p1, p2)
	x02 := interpolateVertical(p0, p2)
	x012 := append(x01, x12...)

	var xLeft, xRight []Point
	// TODO: Clean up this if condition
	if len(x01) > 0 && x01[len(x01)-1].X < x02[len(x01)-1].X {
		xLeft, xRight = x012, x02
	} else {
		xLeft, xRight = x02, x012
	}

	for i := 0; i < len(x02); i++ {
		left, right := xLeft[i].X, xRight[i].X
		for x := left; x <= right; x++ {
			c.PutPixel(x, xLeft[i].Y, clr)
		}
	}
}

// ShadeTriangle shades the triangle formed by the given three points with the specified color, with a gradient.
func (c *Canvas) ShadeTriangle(p0, p1, p2 Point, clr color.RGBA) {
	// TODO: Refactor
	if p1.Y < p0.Y {
		p0, p1 = p1, p0
	}
	if p2.Y < p0.Y {
		p0, p2 = p2, p0
	}
	if p2.Y < p1.Y {
		p1, p2 = p2, p1
	}

	h0, h1, h2 := float32(0.0), float32(0.5), float32(1.0)

	x01 := interpolateVertical(p0, p1)
	h01 := interpolateVertical(Point{X: h0, Y: p0.Y}, Point{X: h1, Y: p1.Y})
	x01 = x01[:len(x01)-1] // Last value overlaps with x12
	h01 = h01[:len(h01)-1] // Last value overlaps with h12

	x12 := interpolateVertical(p1, p2)
	h12 := interpolateVertical(Point{X: h1, Y: p1.Y}, Point{X: h2, Y: p2.Y})

	x02 := interpolateVertical(p0, p2)
	h02 := interpolateVertical(Point{X: h0, Y: p0.Y}, Point{X: h2, Y: p2.Y})

	x012 := append(x01, x12...)
	h012 := append(h01, h12...)

	var xLefts, xRights []Point
	var hLefts, hRights []Point
	if x01[len(x01)-1].X < x02[len(x01)-1].X {
		xLefts, xRights = x012, x02
		hLefts, hRights = h012, h02
	} else {
		xLefts, xRights = x02, x012
		hLefts, hRights = h02, h012
	}

	for i := 0; i < len(x02); i++ {
		xLeft, xRight := xLefts[i].X, xRights[i].X
		hLeft, hRight := hLefts[i].X, hRights[i].X
		hh := interpolateHorizontal(Point{X: xLeft, Y: hLeft}, Point{X: xRight, Y: hRight})
		for x := xLeft; x <= xRight; x++ {
			hGrad := hh[int(x-xLeft)].Y
			gradColor := color.RGBA{
				R: uint8(hGrad * float32(clr.R)),
				G: uint8(hGrad * float32(clr.G)),
				B: uint8(hGrad * float32(clr.B)),
				A: 0xFF,
			}
			c.PutPixel(x, xLefts[i].Y, gradColor)
		}
	}
}

func interpolate(p0, p1 Point) []Point {
	dy, dx := p1.Y-p0.Y, p1.X-p0.X
	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		return interpolateHorizontal(p0, p1)
	}
	return interpolateVertical(p0, p1)
}

func interpolateHorizontal(p0, p1 Point) []Point {
	if p0.X > p1.X {
		p0, p1 = p1, p0
	}
	dy, dx := p1.Y-p0.Y, p1.X-p0.X
	a := dy / dx
	y := p0.Y
	verts := make([]Point, 0)
	for x := p0.X; x <= p1.X; x++ {
		verts = append(verts, Point{x, y})
		y += a
	}
	return verts
}

func interpolateVertical(p0, p1 Point) []Point {
	if p0.Y > p1.Y {
		p0, p1 = p1, p0
	}
	dy, dx := p1.Y-p0.Y, p1.X-p0.X
	a := dx / dy
	x := p0.X
	verts := make([]Point, 0)
	for y := p0.Y; y <= p1.Y; y++ {
		verts = append(verts, Point{x, y})
		x += a
	}
	return verts
}

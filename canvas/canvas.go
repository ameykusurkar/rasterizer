package canvas

import (
	"image"
	"image/color"
	"math"

	geom "rasterizer/geometry"
)

// IndexedTriangleList represents shapes using triangles.
type IndexedTriangleList struct {
	Vertices []geom.Vec3
	Indices  []int
}

// Canvas is a buffer on which we can draw lines, triangles etc.
type Canvas struct {
	image       *image.RGBA
	depthBuffer [][]float32
}

// NewCanvas returns a new Canvas with dimensions (width, height).
func NewCanvas(width, height int) *Canvas {
	return &Canvas{
		image: image.NewRGBA(image.Rect(0, 0, width, height)),
		// Default depth is positive infinity
		depthBuffer: make2dBuffer(width, height, float32(math.Inf(1))),
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

// Clear resets all the canvas pixels and depth buffer.
func (c *Canvas) Clear() {
	bounds := c.image.Bounds()
	for i := 0; i < bounds.Max.X; i++ {
		for j := 0; j < bounds.Max.Y; j++ {
			c.PutPixel(i, j, color.RGBA{0, 0, 0, 0xFF})
			c.depthBuffer[j][i] = float32(math.Inf(1))
		}
	}
}

// PutPixel puts at pixel at (x, y) on the Canvas, with (0, 0) as the top-left corner.
func (c *Canvas) PutPixel(x, y int, color color.Color) {
	c.image.Set(x, y, color)
}

// TestAndSet sets the depth value at (x, y) if it is smallest than the existing,
// and returns whether the depth was set.
func (c *Canvas) TestAndSet(x, y int, depth float32) bool {
	if point := (image.Point{X: x, Y: y}); !point.In(c.image.Bounds()) {
		return false
	}

	if depth < c.depthBuffer[y][x] {
		c.depthBuffer[y][x] = depth
		return true
	}
	return false
}

// FillTriangle fills the triangle formed by the given three points with the
// specified color, using the top-left rule.
func (c *Canvas) FillTriangle(v0, v1, v2 TexVertex, tex Texture) {
	// Sort points by their Y-coordinate
	if v1.Pos.Y < v0.Pos.Y {
		v0, v1 = v1, v0
	}
	if v2.Pos.Y < v0.Pos.Y {
		v0, v2 = v2, v0
	}
	if v2.Pos.Y < v1.Pos.Y {
		v1, v2 = v2, v1
	}
	vTop, vMid, vBottom := v0, v1, v2

	switch {
	case vTop.Pos.Y == vMid.Pos.Y:
		c.fillTriangleFlatTop(vTop, vMid, vBottom, tex)
	case vMid.Pos.Y == vBottom.Pos.Y:
		c.fillTriangleFlatBottom(vTop, vMid, vBottom, tex)
	default:
		alpha := (vMid.Pos.Y - vTop.Pos.Y) / (vBottom.Pos.Y - vTop.Pos.Y)
		vSplit := vTop.InterpolateTo(vBottom, alpha)

		c.fillTriangleFlatBottom(vTop, vMid, vSplit, tex)
		c.fillTriangleFlatTop(vMid, vSplit, vBottom, tex)
	}
}

func (c *Canvas) fillTriangleFlatTop(vLeft, vRight, vBottom TexVertex, tex Texture) {
	if vRight.Pos.X < vLeft.Pos.X {
		vLeft, vRight = vRight, vLeft
	}

	deltaY := vBottom.Pos.Y - vLeft.Pos.Y
	stepLeft := vBottom.Sub(vLeft).Scale(1 / deltaY)
	stepRight := vBottom.Sub(vRight).Scale(1 / deltaY)

	// Round half down to follow the top-left rule
	yStart, yEnd := int(roundHalfDown(vLeft.Pos.Y)), int(roundHalfDown(vBottom.Pos.Y))

	c.fillTriangleFlat(vLeft, vRight, stepLeft, stepRight, yStart, yEnd, tex)
}

func (c *Canvas) fillTriangleFlatBottom(vTop, vLeft, vRight TexVertex, tex Texture) {
	if vRight.Pos.X < vLeft.Pos.X {
		vLeft, vRight = vRight, vLeft
	}

	deltaY := vLeft.Pos.Y - vTop.Pos.Y
	stepLeft := vLeft.Sub(vTop).Scale(1 / deltaY)
	stepRight := vRight.Sub(vTop).Scale(1 / deltaY)

	// Round half down to follow the top-left rule
	yStart, yEnd := int(roundHalfDown(vTop.Pos.Y)), int(roundHalfDown(vLeft.Pos.Y))

	c.fillTriangleFlat(vLeft, vRight, stepLeft, stepRight, yStart, yEnd, tex)
}

func (c *Canvas) fillTriangleFlat(
	vLeft, vRight, stepLeft, stepRight TexVertex,
	yStart, yEnd int,
	tex Texture) {
	// Add 0.5 because we want to use the midpoint of the pixel
	scanLeft := vLeft.Add(stepLeft.Scale(float32(yStart) + 0.5 - vLeft.Pos.Y))
	scanRight := vRight.Add(stepRight.Scale(float32(yStart) + 0.5 - vRight.Pos.Y))

	for y := yStart; y < yEnd; y++ {
		// Round half down to follow the top-left rule
		xStart, xEnd := int(roundHalfDown(scanLeft.Pos.X)), int(roundHalfDown(scanRight.Pos.X))

		deltaX := scanRight.Pos.X - scanLeft.Pos.X
		step := scanRight.Sub(scanLeft).Scale(1 / deltaX)
		scanCoord := scanLeft.Add(step.Scale(float32(xStart) + 0.5 - scanLeft.Pos.X))

		for x := xStart; x < xEnd; x++ {
			// We stored 1/Z in the Z-component so that interpolation will preserve
			// depth perspective. We need to undo the multiplication to get the original
			// texture coordinates.
			depth := 1 / scanCoord.Pos.Z

			// We test the pixel to be drawn against the depth buffer; we only want to draw it
			// if it will be on top of anything already present.
			if c.TestAndSet(x, y, depth) {
				c.PutPixel(x, y, tex.shade(scanCoord.Scale(depth)))
			}
			scanCoord = scanCoord.Add(step)
		}

		scanLeft = scanLeft.Add(stepLeft)
		scanRight = scanRight.Add(stepRight)
	}
}

// roundHalfDown rounds x to the nearest integer, but 0.5 is rounded down.
func roundHalfDown(x float32) float32 {
	return float32(math.Ceil(float64(x) - 0.5))
}

// DrawLine draws a line with a specified color between two points.
func (c *Canvas) DrawLine(p0, p1 geom.Vec2, clr color.Color) {
	for _, vert := range interpolate(p0, p1) {
		c.PutPixel(int(vert.X), int(vert.Y), clr)
	}
}

func interpolate(p0, p1 geom.Vec2) []geom.Vec2 {
	dy, dx := p1.Y-p0.Y, p1.X-p0.X
	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		return interpolateHorizontal(p0, p1)
	}
	return interpolateVertical(p0, p1)
}

func interpolateHorizontal(p0, p1 geom.Vec2) []geom.Vec2 {
	if p0.X > p1.X {
		p0, p1 = p1, p0
	}
	dy, dx := p1.Y-p0.Y, p1.X-p0.X
	a := dy / dx
	y := p0.Y
	verts := make([]geom.Vec2, 0)
	for x := p0.X; x <= p1.X; x++ {
		verts = append(verts, geom.Vec2{X: x, Y: y})
		y += a
	}
	return verts
}

func interpolateVertical(p0, p1 geom.Vec2) []geom.Vec2 {
	if p0.Y > p1.Y {
		p0, p1 = p1, p0
	}
	dy, dx := p1.Y-p0.Y, p1.X-p0.X
	a := dx / dy
	x := p0.X
	verts := make([]geom.Vec2, 0)
	for y := p0.Y; y <= p1.Y; y++ {
		verts = append(verts, geom.Vec2{X: x, Y: y})
		x += a
	}
	return verts
}

func make2dBuffer(width, height int, val float32) [][]float32 {
	buffer := make([][]float32, height)
	for j := 0; j < len(buffer); j++ {
		buffer[j] = make([]float32, width)
		for i := 0; i < len(buffer); i++ {
			buffer[j][i] = val
		}
	}
	return buffer
}

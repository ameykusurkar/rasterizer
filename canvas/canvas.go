package canvas

import (
	"image"
	"image/color"
	"math"

	geom "rasterizer/geometry"
)

// Canvas is a buffer on which we can draw lines, triangles etc.
type Canvas struct {
	image *image.RGBA
}

// NewCanvas returns a new Canvas with dimensions (width, height).
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
func (c *Canvas) PutPixel(x, y int, color color.Color) {
	c.image.Set(x, y, color)
}

// FillTriangle fills the triangle formed by the given three points with the
// specified color, using the top-left rule.
func (c *Canvas) FillTriangle(v0, v1, v2 TexVertex, tex *Texture) {
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
		vSplit := TexVertex{
			Pos:    vTop.Pos.InterpolateTo(vBottom.Pos, alpha),
			TexPos: vTop.TexPos.InterpolateTo(vBottom.TexPos, alpha),
		}

		c.fillTriangleFlatBottom(vTop, vMid, vSplit, tex)
		c.fillTriangleFlatTop(vMid, vSplit, vBottom, tex)
	}
}

func (c *Canvas) fillTriangleFlatTop(vLeft, vRight, vBottom TexVertex, tex *Texture) {
	if vRight.Pos.X < vLeft.Pos.X {
		vLeft, vRight = vRight, vLeft
	}

	stepLeft := vBottom.Sub(vLeft).Scale(1 / (vBottom.Pos.Y - vLeft.Pos.Y))
	stepRight := vBottom.Sub(vRight).Scale(1 / (vBottom.Pos.Y - vRight.Pos.Y))

	// Round half down to follow the top-left rule
	yStart, yEnd := int(roundHalfDown(vLeft.Pos.Y)), int(roundHalfDown(vBottom.Pos.Y))

	// Add 0.5 because we want to use the midpoint of the pixel
	scanLeft := vLeft.Add(stepLeft.Scale(float32(yStart) + 0.5 - vLeft.Pos.Y))
	scanRight := vRight.Add(stepRight.Scale(float32(yStart) + 0.5 - vRight.Pos.Y))

	for y := yStart; y < yEnd; y++ {
		// Round half down to follow the top-left rule
		xStart, xEnd := int(roundHalfDown(scanLeft.Pos.X)), int(roundHalfDown(scanRight.Pos.X))

		stepTex := scanRight.TexPos.Sub(scanLeft.TexPos).Scale(1 / (scanRight.Pos.X - scanLeft.Pos.X))
		texCoord := scanLeft.TexPos.Add(stepTex.Scale(float32(xStart) + 0.5 - scanLeft.Pos.X))

		for x := xStart; x < xEnd; x++ {
			c.PutPixel(x, y, tex.shade(texCoord))
			texCoord = texCoord.Add(stepTex)
		}

		scanLeft = scanLeft.Add(stepLeft)
		scanRight = scanRight.Add(stepRight)
	}
}

func (c *Canvas) fillTriangleFlatBottom(vTop, vLeft, vRight TexVertex, tex *Texture) {
	if vRight.Pos.X < vLeft.Pos.X {
		vLeft, vRight = vRight, vLeft
	}

	stepLeft := vLeft.Sub(vTop).Scale(1 / (vLeft.Pos.Y - vTop.Pos.Y))
	stepRight := vRight.Sub(vTop).Scale(1 / (vRight.Pos.Y - vTop.Pos.Y))

	// Round half down to follow the top-left rule
	yStart, yEnd := int(roundHalfDown(vTop.Pos.Y)), int(roundHalfDown(vLeft.Pos.Y))

	// Add 0.5 because we want to use the midpoint of the pixel
	scanLeft := vLeft.Add(stepLeft.Scale(float32(yStart) + 0.5 - vLeft.Pos.Y))
	scanRight := vRight.Add(stepRight.Scale(float32(yStart) + 0.5 - vRight.Pos.Y))

	for y := yStart; y < yEnd; y++ {
		// Round half down to follow the top-left rule
		xStart, xEnd := int(roundHalfDown(scanLeft.Pos.X)), int(roundHalfDown(scanRight.Pos.X))

		stepTex := scanRight.TexPos.Sub(scanLeft.TexPos).Scale(1 / (scanRight.Pos.X - scanLeft.Pos.X))
		texCoord := scanLeft.TexPos.Add(stepTex.Scale(float32(xStart) + 0.5 - scanLeft.Pos.X))

		for x := xStart; x < xEnd; x++ {
			c.PutPixel(x, y, tex.shade(texCoord))
			texCoord = texCoord.Add(stepTex)
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

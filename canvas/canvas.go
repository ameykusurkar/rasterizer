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
			Pos:    geom.Interpolate2D(vTop.Pos, vBottom.Pos, alpha),
			TexPos: geom.Interpolate2D(vTop.TexPos, vBottom.TexPos, alpha),
		}

		c.fillTriangleFlatBottom(vTop, vMid, vSplit, tex)
		c.fillTriangleFlatTop(vMid, vSplit, vBottom, tex)
	}
}

func (c *Canvas) fillTriangleFlatTop(vLeft, vRight, vBottom TexVertex, tex *Texture) {
	if vRight.Pos.X < vLeft.Pos.X {
		vLeft, vRight = vRight, vLeft
	}
	texLeft, texRight, texBottom := vLeft.TexPos, vRight.TexPos, vBottom.TexPos
	pLeft, pRight, pBottom := vLeft.Pos, vRight.Pos, vBottom.Pos

	// Calculate the dx/dy slope because x is the dependent variable; ie. how much
	// to increment x by as we iterate down the Y-axis.
	mLeft := pBottom.Sub(pLeft).Scale(1 / (pBottom.Y - pLeft.Y))
	mRight := pBottom.Sub(pRight).Scale(1 / (pBottom.Y - pRight.Y))

	mLeftTex := texBottom.Sub(texLeft).Scale(1 / (pBottom.Y - pLeft.Y))
	mRightTex := texBottom.Sub(texRight).Scale(1 / (pBottom.Y - pRight.Y))

	// Round half down to follow the top-left rule
	yStart, yEnd := int(roundHalfDown(pLeft.Y)), int(roundHalfDown(pBottom.Y))

	// Add 0.5 because we want to use the midpoint of the pixel
	scanLeft := pLeft.Add(mLeft.Scale(float32(yStart) + 0.5 - pLeft.Y))
	scanRight := pRight.Add(mRight.Scale(float32(yStart) + 0.5 - pRight.Y))

	texLeft = texLeft.Add(mLeftTex.Scale(float32(yStart) + 0.5 - pLeft.Y))
	texRight = texRight.Add(mRightTex.Scale(float32(yStart) + 0.5 - pRight.Y))

	for y := yStart; y < yEnd; y++ {
		// Round half down to follow the top-left rule
		xStart, xEnd := int(roundHalfDown(scanLeft.X)), int(roundHalfDown(scanRight.X))

		mScanTex := texRight.Sub(texLeft).Scale(1 / (scanRight.X - scanLeft.X))
		texCoord := texLeft.Add(mScanTex.Scale(float32(xStart) + 0.5 - scanLeft.X))

		for x := xStart; x < xEnd; x++ {
			c.PutPixel(x, y, tex.shade(texCoord))
			texCoord = texCoord.Add(mScanTex)
		}

		scanLeft = scanLeft.Add(mLeft)
		scanRight = scanRight.Add(mRight)

		texLeft = texLeft.Add(mLeftTex)
		texRight = texRight.Add(mRightTex)
	}
}

func (c *Canvas) fillTriangleFlatBottom(vTop, vLeft, vRight TexVertex, tex *Texture) {
	if vRight.Pos.X < vLeft.Pos.X {
		vLeft, vRight = vRight, vLeft
	}
	texTop, texLeft, texRight := vTop.TexPos, vLeft.TexPos, vRight.TexPos
	pTop, pLeft, pRight := vTop.Pos, vLeft.Pos, vRight.Pos

	// Calculate the dx/dy slope because x is the dependent variable; ie. how much
	// to increment x by as we iterate down the Y-axis.
	mLeft := (pLeft.X - pTop.X) / (pLeft.Y - pTop.Y)
	mRight := (pRight.X - pTop.X) / (pRight.Y - pTop.Y)

	mLeftTex := texLeft.Sub(texTop).Scale(1 / (pLeft.Y - pTop.Y))
	mRightTex := texRight.Sub(texTop).Scale(1 / (pRight.Y - pTop.Y))

	// Round half down to follow the top-left rule
	yStart, yEnd := int(roundHalfDown(pTop.Y)), int(roundHalfDown(pLeft.Y))

	// Add 0.5 because we want to use the midpoint of the pixel
	scanLeft := mLeft*(float32(yStart)+0.5-pLeft.Y) + pLeft.X
	scanRight := mRight*(float32(yStart)+0.5-pRight.Y) + pRight.X

	texLeft = texLeft.Add(mLeftTex.Scale(float32(yStart) + 0.5 - pLeft.Y))
	texRight = texRight.Add(mRightTex.Scale(float32(yStart) + 0.5 - pRight.Y))

	for y := yStart; y < yEnd; y++ {
		// Round half down to follow the top-left rule
		xStart, xEnd := int(roundHalfDown(scanLeft)), int(roundHalfDown(scanRight))

		mScanTex := texRight.Sub(texLeft).Scale(1 / (scanRight - scanLeft))
		texCoord := texLeft.Add(mScanTex.Scale(float32(xStart) + 0.5 - scanLeft))

		for x := xStart; x < xEnd; x++ {
			c.PutPixel(x, y, tex.shade(texCoord))
			texCoord = texCoord.Add(mScanTex)
		}

		scanLeft += mLeft
		scanRight += mRight

		texLeft = texLeft.Add(mLeftTex)
		texRight = texRight.Add(mRightTex)
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

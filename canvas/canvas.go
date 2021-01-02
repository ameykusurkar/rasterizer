package canvas

import (
	"image"
	"image/color"
	"math"

	geom "rasterizer/geometry"
)

var defaultTexture Texture = Texture{
	points: []geom.Vec2{{X: 0, Y: 1}, {X: 1, Y: 0}, {X: 0.5, Y: 0.5}},
}

// Texture is a thing, exactly what it is TBD.
type Texture struct {
	points []geom.Vec2
}

func (tex *Texture) shade(texCoord geom.Vec2) color.Color {
	red := geom.Vec3{X: 255, Y: 0, Z: 0}
	blue := geom.Vec3{X: 0, Y: 0, Z: 255}
	colorVec := red.Scale(texCoord.X).Add(blue.Scale(texCoord.Y))
	return color.RGBA{uint8(colorVec.X), uint8(colorVec.Y), uint8(colorVec.Z), 255}
}

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

// DrawLine draws a line with a specified color between two points.
func (c *Canvas) DrawLine(p0, p1 geom.Vec2, clr color.Color) {
	for _, vert := range interpolate(p0, p1) {
		c.PutPixel(int(vert.X), int(vert.Y), clr)
	}
}

// FillTriangle fills the triangle formed by the given three points with the
// specified color, using the top-left rule.
func (c *Canvas) FillTriangle(p0, p1, p2 geom.Vec2, clr color.Color) {
	// Sort points by their Y-coordinate
	if p1.Y < p0.Y {
		p0, p1 = p1, p0
	}
	if p2.Y < p0.Y {
		p0, p2 = p2, p0
	}
	if p2.Y < p1.Y {
		p1, p2 = p2, p1
	}
	pTop, pMid, pBottom := p0, p1, p2

	switch {
	case pTop.Y == pMid.Y:
		c.fillTriangleFlatTop(pTop, pMid, pBottom, clr)
	case pMid.Y == pBottom.Y:
		c.fillTriangleFlatBottom(pTop, pMid, pBottom, clr)
	default:
		mSplit := (pBottom.X - pTop.X) / (pBottom.Y - pTop.Y)
		pSplit := geom.Vec2{
			X: mSplit*(pMid.Y-pTop.Y) + pTop.X,
			Y: pMid.Y,
		}

		c.fillTriangleFlatBottom(pTop, pMid, pSplit, clr)
		c.fillTriangleFlatTop(pMid, pSplit, pBottom, clr)
	}
}

func (c *Canvas) fillTriangleFlatTop(pLeft, pRight, pBottom geom.Vec2, clr color.Color) {
	if pRight.X < pLeft.X {
		pLeft, pRight = pRight, pLeft
	}
	texLeft, texRight, texBottom := defaultTexture.points[0], defaultTexture.points[1], defaultTexture.points[2]

	// Calculate the dx/dy slope because x is the dependent variable; ie. how much
	// to increment x by as we iterate down the Y-axis.
	mLeft := (pBottom.X - pLeft.X) / (pBottom.Y - pLeft.Y)
	mRight := (pBottom.X - pRight.X) / (pBottom.Y - pRight.Y)

	mLeftTex := texBottom.Sub(texLeft).Scale(1 / (pBottom.Y - pLeft.Y))
	mRightTex := texBottom.Sub(texRight).Scale(1 / (pBottom.Y - pRight.Y))

	// Round half down to follow the top-left rule
	yStart, yEnd := int(roundHalfDown(pLeft.Y)), int(roundHalfDown(pBottom.Y))

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
			c.PutPixel(x, y, defaultTexture.shade(texCoord))
			texCoord = texCoord.Add(mScanTex)
		}

		scanLeft += mLeft
		scanRight += mRight

		texLeft = texLeft.Add(mLeftTex)
		texRight = texRight.Add(mRightTex)
	}
}

func (c *Canvas) fillTriangleFlatBottom(pTop, pLeft, pRight geom.Vec2, clr color.Color) {
	if pRight.X < pLeft.X {
		pLeft, pRight = pRight, pLeft
	}

	// Calculate the dx/dy slope because x is the dependent variable; ie. how much
	// to increment x by as we iterate down the Y-axis.
	mLeft := (pLeft.X - pTop.X) / (pLeft.Y - pTop.Y)
	mRight := (pRight.X - pTop.X) / (pRight.Y - pTop.Y)

	// Round half down to follow the top-left rule
	yStart, yEnd := int(roundHalfDown(pTop.Y)), int(roundHalfDown(pLeft.Y))
	for y := yStart; y < yEnd; y++ {
		yF := float32(y)

		// Add 0.5 because we want to use the midpoint of the pixel
		xStartF := mLeft*(yF+0.5-pLeft.Y) + pLeft.X
		xEndF := mRight*(yF+0.5-pRight.Y) + pRight.X

		// Round half down to follow the top-left rule
		xStart, xEnd := int(roundHalfDown(xStartF)), int(roundHalfDown(xEndF))

		for x := xStart; x < xEnd; x++ {
			c.PutPixel(x, y, clr)
		}
	}
}

// roundHalfDown rounds x to the nearest integer, but 0.5 is rounded down.
func roundHalfDown(x float32) float32 {
	return float32(math.Ceil(float64(x) - 0.5))
}

// ShadeTriangle shades the triangle formed by the given three points with the specified color, with a gradient.
// TODO: This uses an old algorithm and needs to be updated.
func (c *Canvas) ShadeTriangle(p0, p1, p2 geom.Vec2, clr color.RGBA) {
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
	h01 := interpolateVertical(geom.Vec2{X: h0, Y: p0.Y}, geom.Vec2{X: h1, Y: p1.Y})
	x01 = x01[:len(x01)-1] // Last value overlaps with x12
	h01 = h01[:len(h01)-1] // Last value overlaps with h12

	x12 := interpolateVertical(p1, p2)
	h12 := interpolateVertical(geom.Vec2{X: h1, Y: p1.Y}, geom.Vec2{X: h2, Y: p2.Y})

	x02 := interpolateVertical(p0, p2)
	h02 := interpolateVertical(geom.Vec2{X: h0, Y: p0.Y}, geom.Vec2{X: h2, Y: p2.Y})

	x012 := append(x01, x12...)
	h012 := append(h01, h12...)

	var xLefts, xRights []geom.Vec2
	var hLefts, hRights []geom.Vec2
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
		hh := interpolateHorizontal(geom.Vec2{X: xLeft, Y: hLeft}, geom.Vec2{X: xRight, Y: hRight})
		for x := xLeft; x <= xRight; x++ {
			hGrad := hh[int(x-xLeft)].Y
			gradColor := color.RGBA{
				R: uint8(hGrad * float32(clr.R)),
				G: uint8(hGrad * float32(clr.G)),
				B: uint8(hGrad * float32(clr.B)),
				A: 0xFF,
			}
			c.PutPixel(int(x), int(xLefts[i].Y), gradColor)
		}
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

package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"rasterizer/canvas"
	geom "rasterizer/geometry"
)

const (
	screenWidth  = 400
	screenHeight = 400
)

type game struct {
	canv   canvas.Canvas
	cube   geom.IndexedTriangleList
	thetaX float32
	thetaY float32
	thetaZ float32
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Rasterizer")
	ebiten.SetMaxTPS(60)
	g := game{
		canv:   *canvas.NewCanvas(screenWidth, screenHeight),
		cube:   *buildCube(),
		thetaZ: 0,
	}
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}

func (g *game) paint() {
	var d float32 = 1
	w, h := g.canv.Dimensions()

	projectedPoints := make([]geom.Vec2, 0, len(g.cube.Vertices))
	rMat := geom.RotationX(g.thetaX).
		MatMul(geom.RotationY(g.thetaY)).
		MatMul(geom.RotationZ(g.thetaZ))
	for _, v := range g.cube.Vertices {
		rotated := rMat.VecMul(v)
		projected := geom.Project(rotated, d)
		point := vertexToPoint(projected, w, h)
		projectedPoints = append(projectedPoints, point)
	}

	g.canv.Fill(color.RGBA{0, 0, 0, 0xFF})

	red := color.RGBA{255, 0, 0, 255}
	for i := 0; i < len(g.cube.Indices); i += 3 {
		idx0, idx1, idx2 := g.cube.Indices[i], g.cube.Indices[i+1], g.cube.Indices[i+2]
		p0, p1, p2 := projectedPoints[idx0], projectedPoints[idx1], projectedPoints[idx2]

		blue := color.RGBA{
			0, 0, 255 - uint8(255*float32(i)/float32(len(g.cube.Indices))), 255,
		}

		g.canv.FillTriangle(p0, p1, p2, blue)
		g.canv.DrawLine(p0, p1, red)
		g.canv.DrawLine(p1, p2, red)
		g.canv.DrawLine(p2, p0, red)
	}
}

func buildCube() *geom.IndexedTriangleList {
	return &geom.IndexedTriangleList{
		Vertices: []geom.Vec3{
			// Front
			{X: -1, Y: 1, Z: 2},
			{X: 1, Y: 1, Z: 2},
			{X: 1, Y: -1, Z: 2},
			{X: -1, Y: -1, Z: 2},
			// Back
			{X: -1, Y: 1, Z: 3},
			{X: 1, Y: 1, Z: 3},
			{X: 1, Y: -1, Z: 3},
			{X: -1, Y: -1, Z: 3},
		},
		Indices: []int{
			// Front
			0, 1, 3,
			3, 1, 2,

			// Back
			4, 5, 7,
			7, 5, 6,

			// Left
			0, 4, 7,
			7, 0, 3,

			// Right
			2, 1, 5,
			2, 5, 6,

			// Top
			0, 4, 5,
			0, 5, 1,

			// Bottom
			3, 7, 6,
			3, 6, 2,
		},
	}
}

func vertexToPoint(v geom.Vec3, width int, height int) geom.Vec2 {
	halfWidth, halfHeight := float32(width)/2, float32(height)/2
	x := (1 + v.X) * halfWidth
	y := (1 - v.Y) * halfHeight
	return geom.Vec2{X: float32(int(x)), Y: float32(int(y))}
}

func (g *game) Update() error {
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.thetaZ += 0.05
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		g.thetaZ -= 0.05
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.thetaX += 0.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.thetaX -= 0.05
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.thetaY += 0.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.thetaY -= 0.05
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	g.paint()
	screen.ReplacePixels(g.canv.Buffer())
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.canv.Dimensions()
}

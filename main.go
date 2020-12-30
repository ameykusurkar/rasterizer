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
	screenWidth  = 320
	screenHeight = 320
)

type game struct {
	canv   canvas.Canvas
	cube   geom.IndexedLineList
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

	projectedPoints := make([]canvas.Point, 0, len(g.cube.Vertices))
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
	for i := 0; i < len(g.cube.Indices)/2; i++ {
		idx0, idx1 := g.cube.Indices[2*i], g.cube.Indices[2*i+1]
		p0, p1 := projectedPoints[idx0], projectedPoints[idx1]
		g.canv.DrawLine(p0, p1, red)
	}
}

func buildCube() *geom.IndexedLineList {
	return &geom.IndexedLineList{
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
			0, 1,
			1, 2,
			2, 3,
			3, 0,

			4, 5,
			5, 6,
			6, 7,
			7, 4,

			0, 4,
			1, 5,
			2, 6,
			3, 7,
		},
	}
}

func vertexToPoint(v geom.Vec3, width int, height int) canvas.Point {
	halfWidth, halfHeight := float32(width)/2, float32(height)/2
	x := (1 + v.X) * halfWidth
	y := (1 - v.Y) * halfHeight
	return canvas.Point{X: float32(int(x)), Y: float32(int(y))}
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

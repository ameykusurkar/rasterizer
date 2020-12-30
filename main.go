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
	canv canvas.Canvas
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Rasterizer")
	ebiten.SetMaxTPS(60)
	g := game{canv: *canvas.NewCanvas(screenWidth, screenHeight)}
	paint(&g.canv)
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}

func paint(canv *canvas.Canvas) {
	var d float32 = 0.3
	cube := buildCube()
	w, h := canv.Dimensions()

	projectedPoints := make([]canvas.Point, 0, len(cube.Vertices))
	for _, v := range cube.Vertices {
		projected := geom.Project(v, d)
		point := vertexToPoint(projected, w, h)
		projectedPoints = append(projectedPoints, point)
	}

	canv.Fill(color.RGBA{0, 0, 0, 0xFF})

	red := color.RGBA{255, 0, 0, 255}
	for i := 0; i < len(cube.Indices)/2; i++ {
		idx0, idx1 := cube.Indices[2*i], cube.Indices[2*i+1]
		p0, p1 := projectedPoints[idx0], projectedPoints[idx1]
		canv.DrawLine(p0, p1, red)
	}
}

func buildCube() *geom.IndexedLineList {
	return &geom.IndexedLineList{
		Vertices: []geom.Vec3{
			// Front
			{X: -2, Y: 1, Z: 1},
			{X: 0, Y: 1, Z: 1},
			{X: 0, Y: -1, Z: 1},
			{X: -2, Y: -1, Z: 1},
			// Back
			{X: -2, Y: 1, Z: 2},
			{X: 0, Y: 1, Z: 2},
			{X: 0, Y: -1, Z: 2},
			{X: -2, Y: -1, Z: 2},
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
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.ReplacePixels(g.canv.Buffer())
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.canv.Dimensions()
}

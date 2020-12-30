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
	screenHeight = 240
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
	w, h := canv.Dimensions()
	fmt.Printf("%d, %d\n", w, h)

	var d float32 = 0.3

	vAf := geom.Vertex{X: -2, Y: 1, Z: 1}
	vBf := geom.Vertex{X: 0, Y: 1, Z: 1}
	vCf := geom.Vertex{X: 0, Y: -1, Z: 1}
	vDf := geom.Vertex{X: -2, Y: -1, Z: 1}

	vAb := geom.Vertex{X: -2, Y: 1, Z: 2}
	vBb := geom.Vertex{X: 0, Y: 1, Z: 2}
	vCb := geom.Vertex{X: 0, Y: -1, Z: 2}
	vDb := geom.Vertex{X: -2, Y: -1, Z: 2}

	pvAf := geom.Project(vAf, d)
	pvBf := geom.Project(vBf, d)
	pvCf := geom.Project(vCf, d)
	pvDf := geom.Project(vDf, d)
	fmt.Printf("%v, %v, %v, %v\n", pvAf, pvBf, pvCf, pvDf)

	pvAb := geom.Project(vAb, d)
	pvBb := geom.Project(vBb, d)
	pvCb := geom.Project(vCb, d)
	pvDb := geom.Project(vDb, d)
	fmt.Printf("%v, %v, %v, %v\n", pvAb, pvBb, pvCb, pvDb)

	cpvAf := vertexToPoint(pvAf, w, h)
	cpvBf := vertexToPoint(pvBf, w, h)
	cpvCf := vertexToPoint(pvCf, w, h)
	cpvDf := vertexToPoint(pvDf, w, h)
	fmt.Printf("%v, %v, %v, %v\n", cpvAf, cpvBf, cpvCf, cpvDf)

	cpvAb := vertexToPoint(pvAb, w, h)
	cpvBb := vertexToPoint(pvBb, w, h)
	cpvCb := vertexToPoint(pvCb, w, h)
	cpvDb := vertexToPoint(pvDb, w, h)
	fmt.Printf("%v, %v, %v, %v\n", cpvAb, cpvBb, cpvCb, cpvDb)

	canv.Fill(color.RGBA{0, 0, 0, 0xFF})

	red := color.RGBA{255, 0, 0, 255}
	green := color.RGBA{0, 255, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}

	canv.DrawLine(cpvAf, cpvBf, blue)
	canv.DrawLine(cpvBf, cpvCf, blue)
	canv.DrawLine(cpvCf, cpvDf, blue)
	canv.DrawLine(cpvDf, cpvAf, blue)

	canv.DrawLine(cpvAb, cpvBb, red)
	canv.DrawLine(cpvBb, cpvCb, red)
	canv.DrawLine(cpvCb, cpvDb, red)
	canv.DrawLine(cpvDb, cpvAb, red)

	canv.DrawLine(cpvAf, cpvAb, green)
	canv.DrawLine(cpvBf, cpvBb, green)
	canv.DrawLine(cpvCf, cpvCb, green)
	canv.DrawLine(cpvDf, cpvDb, green)
}

func vertexToPoint(v geom.Vertex, width int, height int) canvas.Point {
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

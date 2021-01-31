package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"rasterizer/canvas"
	geom "rasterizer/geometry"
)

const (
	screenWidth  = 400
	screenHeight = 400
)

var colors = []color.RGBA{
	{255, 0, 0, 255},
	{0, 255, 0, 255},
	{0, 0, 255, 255},
	{0, 255, 255, 255},
	{255, 255, 0, 255},
	{255, 0, 255, 255},
}

type game struct {
	pipeline     Pipeline
	cubes        []canvas.IndexedTriangleList
	tex          canvas.ImageTextureWrapped
	vertexShader *DefaultVertexShader
	thetaX       float32
	thetaY       float32
	thetaZ       float32
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Rasterizer")
	ebiten.SetMaxTPS(60)

	img, err := imageFromPath("insta.png")

	if err != nil {
		log.Fatal(err)
	}

	cubes := make([]canvas.IndexedTriangleList, 0)
	cubes = append(cubes, *buildCube(geom.Vec3{X: -0.5, Y: 1, Z: 4}, 2.0))
	cubes = append(cubes, *buildCube(geom.Vec3{X: 0.5, Y: 0, Z: 5}, 3.5))

	vertexShader := &DefaultVertexShader{
		rotation:       *geom.RotationZ(0),
		rotationCenter: cubes[0].Vertices[0].Pos.Add(cubes[0].Vertices[6].Pos).Scale(0.5),
	}

	g := game{
		pipeline: Pipeline{
			canv:         *canvas.NewCanvas(screenWidth, screenHeight),
			vertexShader: vertexShader,
		},
		vertexShader: vertexShader,
		tex:          canvas.ImageTextureWrapped{Img: img, Scale: 0.25},
		cubes:        cubes,
	}

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}

func imageFromPath(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	image, _, err := image.Decode(f)
	return image, err
}

func buildCube(center geom.Vec3, length float32) *canvas.IndexedTriangleList {
	/*  Orientation of vertices
	 *         4--------5
	 *        /|       /|
	 *       / |      / |
	 *      0--------1  |
	 *      |  7-----|--6
	 *      | /      | /
	 *      |/       |/
	 *      3--------2
	 */
	return &canvas.IndexedTriangleList{
		Vertices: []canvas.TexVertex{
			{
				Pos:    center.Add(geom.Vec3{X: -1, Y: 1, Z: -1}).Scale(length / 2),
				TexPos: geom.Vec2{X: 0, Y: 0},
			},
			{
				Pos:    center.Add(geom.Vec3{X: 1, Y: 1, Z: -1}).Scale(length / 2),
				TexPos: geom.Vec2{X: 1, Y: 0},
			},
			{
				Pos:    center.Add(geom.Vec3{X: 1, Y: -1, Z: -1}).Scale(length / 2),
				TexPos: geom.Vec2{X: 1, Y: 1},
			},
			{
				Pos:    center.Add(geom.Vec3{X: -1, Y: -1, Z: -1}).Scale(length / 2),
				TexPos: geom.Vec2{X: 0, Y: 1},
			},
			{
				Pos:    center.Add(geom.Vec3{X: -1, Y: 1, Z: 1}).Scale(length / 2),
				TexPos: geom.Vec2{X: 0, Y: 1},
			},
			{
				Pos:    center.Add(geom.Vec3{X: 1, Y: 1, Z: 1}).Scale(length / 2),
				TexPos: geom.Vec2{X: 1, Y: 1},
			},
			{
				Pos:    center.Add(geom.Vec3{X: 1, Y: -1, Z: 1}).Scale(length / 2),
				TexPos: geom.Vec2{X: 1, Y: 0},
			},
			{
				Pos:    center.Add(geom.Vec3{X: -1, Y: -1, Z: 1}).Scale(length / 2),
				TexPos: geom.Vec2{X: 0, Y: 0},
			},
		},
		Indices: []int{
			// Front
			3, 0, 1,
			3, 1, 2,

			// Back
			6, 5, 4,
			6, 4, 7,

			// Left
			7, 4, 0,
			7, 0, 3,

			// Right
			2, 1, 5,
			2, 5, 6,

			// Top
			0, 4, 5,
			0, 5, 1,

			// Bottom
			7, 3, 2,
			7, 2, 6,
		},
	}
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
	g.vertexShader.rotation = *geom.RotationX(g.thetaX).
		MatMul(geom.RotationY(g.thetaY)).
		MatMul(geom.RotationZ(g.thetaZ))
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	g.pipeline.canv.Clear()

	for _, cube := range g.cubes {
		g.pipeline.Draw(&cube, &g.tex)
	}

	screen.ReplacePixels(g.pipeline.canv.Buffer())
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.pipeline.canv.Dimensions()
}

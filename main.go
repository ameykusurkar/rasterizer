package main

import (
	"errors"
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
	vertexShader *VertexRotator
	thetaX       float32
	thetaY       float32
	thetaZ       float32
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Rasterizer")
	ebiten.SetMaxTPS(60)

	if len(os.Args) < 2 {
		log.Fatal(errors.New("Please provide an image for texture"))
	}

	img, err := imageFromPath(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	cubes := make([]canvas.IndexedTriangleList, 0)
	cubes = append(cubes, *buildCube(geom.Vec3{X: -0.5, Y: 1, Z: 4}, 2.0))
	cubes = append(cubes, *buildCube(geom.Vec3{X: 0.5, Y: 0, Z: 5}, 3.5))

	vertexShader := &VertexRotator{
		rotation:       *geom.RotationZ(0),
		rotationCenter: cubes[0].Vertices[0].Add(cubes[0].Vertices[6]).Scale(0.5),
	}

	g := game{
		pipeline: Pipeline{
			canv:           *canvas.NewCanvas(screenWidth, screenHeight),
			vertexShader:   vertexShader,
			geometryShader: &CubeShader{},
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
	return &canvas.IndexedTriangleList{
		Vertices: []geom.Vec3{
			center.Add(geom.Vec3{X: -1, Y: 1, Z: -1}).Scale(length / 2),
			center.Add(geom.Vec3{X: 1, Y: 1, Z: -1}).Scale(length / 2),
			center.Add(geom.Vec3{X: 1, Y: -1, Z: -1}).Scale(length / 2),
			center.Add(geom.Vec3{X: -1, Y: -1, Z: -1}).Scale(length / 2),
			center.Add(geom.Vec3{X: -1, Y: 1, Z: 1}).Scale(length / 2),
			center.Add(geom.Vec3{X: 1, Y: 1, Z: 1}).Scale(length / 2),
			center.Add(geom.Vec3{X: 1, Y: -1, Z: 1}).Scale(length / 2),
			center.Add(geom.Vec3{X: -1, Y: -1, Z: 1}).Scale(length / 2),
		},
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

// CubeShader tells you how to shade cubes.
type CubeShader struct{}

var topTriangle = []geom.Vec2{
	{X: 0, Y: 1},
	{X: 0, Y: 0},
	{X: 1, Y: 0},
}

var bottomTriangle = []geom.Vec2{
	{X: 0, Y: 1},
	{X: 1, Y: 0},
	{X: 1, Y: 1},
}

// Process returns the given vertices.
func (s *CubeShader) Process(vertices []geom.Vec3, index int) []canvas.TexVertex {
	processed := make([]canvas.TexVertex, 0, len(vertices))
	for i := 0; i < len(vertices); i++ {
		var texPos geom.Vec2
		if index%2 == 0 {
			texPos = topTriangle[i]
		} else {
			texPos = bottomTriangle[i]
		}
		processed = append(processed, canvas.TexVertex{
			Pos:    vertices[i],
			TexPos: texPos,
		})
	}
	return processed
}

// VertexRotator rotates vertices.
type VertexRotator struct {
	rotation       geom.Mat3
	rotationCenter geom.Vec3
}

// Process rotates the vertex.
func (s *VertexRotator) Process(v geom.Vec3) geom.Vec3 {
	return s.rotation.VecMul(v.Sub(s.rotationCenter)).Add(s.rotationCenter)
}

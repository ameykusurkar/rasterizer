package main

import (
	"image/color"
	"rasterizer/canvas"
	geom "rasterizer/geometry"
)

// Pipeline encapsulates the process of rendering a 3D scene to the screen.
type Pipeline struct {
	canv           canvas.Canvas
	rotation       geom.Mat3
	rotationCenter geom.Vec3
}

var defaultTexture canvas.Texture = canvas.Texture{
	Points: []geom.Vec2{{X: 0, Y: 1}, {X: 1, Y: 0}, {X: 0.5, Y: 0.5}},
}

// Draw renders the given triangles onto the screen.
func (p *Pipeline) Draw(triangleList *geom.IndexedTriangleList) {
	vertices := p.transformVertices(triangleList.Vertices)
	triangles3D := assembleTriangles(vertices, triangleList.Indices)

	white := color.RGBA{255, 255, 255, 255}
	for _, tri3D := range triangles3D {
		tri2D := [3]geom.Vec2{
			p.transformPerspective(tri3D[0]),
			p.transformPerspective(tri3D[1]),
			p.transformPerspective(tri3D[2]),
		}

		p.canv.FillTriangle(tri2D[0], tri2D[1], tri2D[2], &defaultTexture)
		p.canv.DrawLine(tri2D[0], tri2D[1], white)
		p.canv.DrawLine(tri2D[1], tri2D[2], white)
		p.canv.DrawLine(tri2D[2], tri2D[0], white)
	}
}

func colorToVec3(clr color.RGBA) geom.Vec3 {
	return geom.Vec3{
		X: float32(clr.R),
		Y: float32(clr.G),
		Z: float32(clr.B),
	}
}

// Apply any rotation or translation to the vertices if necessary.
func (p *Pipeline) transformVertices(vertices []geom.Vec3) []geom.Vec3 {
	rotatedVertices := make([]geom.Vec3, 0, len(vertices))
	for _, v := range vertices {
		rotated := p.rotation.VecMul(v.Sub(p.rotationCenter)).Add(p.rotationCenter)
		rotatedVertices = append(rotatedVertices, rotated)
	}
	return rotatedVertices
}

// Build triangles from the indexed list. Also applies backface culling.
func assembleTriangles(vertices []geom.Vec3, indices []int) [][3]geom.Vec3 {
	triangles := make([][3]geom.Vec3, 0)

	for i := 0; i < len(indices); i += 3 {
		idx0, idx1, idx2 := indices[i], indices[i+1], indices[i+2]

		v0, v1, v2 := vertices[idx0], vertices[idx1], vertices[idx2]
		// Assumes that the triangle's vertices are defined in clockwise order
		normal := v1.Sub(v0).Cross(v2.Sub(v0))
		if normal.Dot(v0) > 0 {
			// A positive dot-product indicates that the viewing vector is in the same
			// direcion as the triangle's normal. This means that we are looking at the
			// back-face of triangle, which should not be visible.
			continue
		}

		triangles = append(triangles, [3]geom.Vec3{v0, v1, v2})
	}

	return triangles
}

// Transforms the 3D scene to a 2D scene by applying perspective, that can then
// be drawn on a canvas.
func (p *Pipeline) transformPerspective(vertex geom.Vec3) geom.Vec2 {
	w, h := p.canv.Dimensions()
	projected := geom.Project(vertex, 1)
	return vertexToPoint(projected, w, h)
}

func vertexToPoint(v geom.Vec3, width int, height int) geom.Vec2 {
	halfWidth, halfHeight := float32(width)/2, float32(height)/2
	x := (1 + v.X) * halfWidth
	y := (1 - v.Y) * halfHeight
	return geom.Vec2{X: float32(int(x)), Y: float32(int(y))}
}

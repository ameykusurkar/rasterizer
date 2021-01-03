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

var texPoints = []geom.Vec2{{X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}}

// Draw renders the given triangles onto the screen.
func (p *Pipeline) Draw(triangleList *canvas.IndexedTriangleList, tex canvas.Texture) {
	vertices := p.transformVertices(triangleList.Vertices)

	triangles := assembleTriangles(vertices, triangleList.Indices)

	for _, tri := range triangles {
		triProjected := [3]canvas.TexVertex{
			p.transformPerspective(tri[0]),
			p.transformPerspective(tri[1]),
			p.transformPerspective(tri[2]),
		}

		p.canv.FillTriangle(triProjected[0], triProjected[1], triProjected[2], tex)

		tri2D := []geom.Vec2{
			vec3ToVec2(triProjected[0].Pos),
			vec3ToVec2(triProjected[1].Pos),
			vec3ToVec2(triProjected[2].Pos),
		}

		white := color.RGBA{255, 255, 255, 255}
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
func (p *Pipeline) transformVertices(vertices []canvas.TexVertex) []canvas.TexVertex {
	rotatedVertices := make([]canvas.TexVertex, 0, len(vertices))
	for _, v := range vertices {
		rotatedVertices = append(rotatedVertices, canvas.TexVertex{
			Pos:    p.rotation.VecMul(v.Pos.Sub(p.rotationCenter)).Add(p.rotationCenter),
			TexPos: v.TexPos,
		})
	}
	return rotatedVertices
}

// Build triangles from the indexed list. Also applies backface culling.
func assembleTriangles(vertices []canvas.TexVertex, indices []int) [][3]canvas.TexVertex {
	triangles := make([][3]canvas.TexVertex, 0)

	for i := 0; i < len(indices); i += 3 {
		idx0, idx1, idx2 := indices[i], indices[i+1], indices[i+2]

		v0, v1, v2 := vertices[idx0], vertices[idx1], vertices[idx2]

		// Assumes that the triangle's vertices are defined in clockwise order
		if triangleFacingAway(v0.Pos, v1.Pos, v2.Pos) {
			continue
		}

		triangles = append(triangles, [3]canvas.TexVertex{v0, v1, v2})
	}

	return triangles
}

func triangleFacingAway(v0, v1, v2 geom.Vec3) bool {
	normal := v1.Sub(v0).Cross(v2.Sub(v0))

	// A positive dot-product indicates that the viewing vector is in the same
	// direcion as the triangle's normal. This means that we are looking at the
	// back-face of triangle, which should not be visible.
	return normal.Dot(v0) > 0
}

// Transforms the 3D scene to a 2D scene by applying perspective, that can then
// be drawn on a canvas. We still return a 3D vertex: this allows us to retain
// depth information, which will be used later in the pipeline.
func (p *Pipeline) transformPerspective(vertex canvas.TexVertex) canvas.TexVertex {
	depth := vertex.Pos.Z
	vertex = vertex.Scale(1 / depth)
	vertex.Pos.Z = depth

	w, h := p.canv.Dimensions()
	vertex.Pos = vertexToPoint(vertex.Pos, w, h)
	return vertex
}

func vertexToPoint(v geom.Vec3, width int, height int) geom.Vec3 {
	halfWidth, halfHeight := float32(width)/2, float32(height)/2
	v.X = (1 + v.X) * halfWidth
	v.Y = (1 - v.Y) * halfHeight
	return v
}

func vec3ToVec2(v geom.Vec3) geom.Vec2 {
	return geom.Vec2{X: v.X, Y: v.Y}
}

package main

import (
	"image/color"
	"rasterizer/canvas"
	geom "rasterizer/geometry"
)

// VertexShader is a shader in the pipeline that processes individual vertices.
type VertexShader interface {
	Process(v geom.Vec3) geom.Vec3
}

// GeometryShader is a shader in the pipeline that processes assembled triangles.
type GeometryShader interface {
	Process(vertices []geom.Vec3, index int) []canvas.TexVertex
}

// Pipeline encapsulates the process of rendering a 3D scene to the screen.
type Pipeline struct {
	canv           canvas.Canvas
	vertexShader   VertexShader
	geometryShader GeometryShader
}

// Draw renders the given triangles onto the screen.
func (p *Pipeline) Draw(triangleList *canvas.IndexedTriangleList, tex canvas.Texture) {
	vertices := make([]geom.Vec3, 0, len(triangleList.Vertices))
	for _, vertex := range triangleList.Vertices {
		vertices = append(vertices, p.vertexShader.Process(vertex))
	}

	triangles, triangleIndices := assembleTriangles(vertices, triangleList.Indices)

	processedTriangles := make([][]canvas.TexVertex, 0, len(triangles))
	for i := 0; i < len(triangles); i++ {
		processedTriangles = append(processedTriangles, p.geometryShader.Process(triangles[i][:], triangleIndices[i]))
	}

	for _, tri := range processedTriangles {
		p.canv.FillTriangle(
			p.transformPerspective(tri[0]),
			p.transformPerspective(tri[1]),
			p.transformPerspective(tri[2]),
			tex,
		)
	}
}

func colorToVec3(clr color.RGBA) geom.Vec3 {
	return geom.Vec3{
		X: float32(clr.R),
		Y: float32(clr.G),
		Z: float32(clr.B),
	}
}

// Build triangles from the indexed list. Also applies backface culling.
func assembleTriangles(vertices []geom.Vec3, indices []int) ([][3]geom.Vec3, []int) {
	triangles := make([][3]geom.Vec3, 0)
	triangleIndices := make([]int, 0)

	for i := 0; i < len(indices); i += 3 {
		idx0, idx1, idx2 := indices[i], indices[i+1], indices[i+2]
		v0, v1, v2 := vertices[idx0], vertices[idx1], vertices[idx2]

		if triangleFacingAway(v0, v1, v2) {
			continue
		}

		triangles = append(triangles, [3]geom.Vec3{v0, v1, v2})
		triangleIndices = append(triangleIndices, i/3)
	}

	return triangles, triangleIndices
}

func triangleFacingAway(v0, v1, v2 geom.Vec3) bool {
	// Assumes that the triangle's vertices are defined in clockwise order
	normal := v1.Sub(v0).Cross(v2.Sub(v0))

	// A positive dot-product indicates that the viewing vector is in the same
	// direcion as the triangle's normal. This means that we are looking at the
	// back-face of triangle, which should not be visible.
	return normal.Dot(v0) > 0
}

// Transforms the 3D scene to a 2D scene by applying perspective, that can then
// be drawn on a canvas.
func (p *Pipeline) transformPerspective(vertex canvas.TexVertex) canvas.TexVertex {
	w, h := p.canv.Dimensions()
	zInv := 1 / vertex.Pos.Z

	// We also want to transform the texture coordinates so that perspective is
	// applied correctly to the texture. We will re-multiply the texture coordinates
	// by the Z component before drawing the pixel.
	projected := vertex.Scale(zInv)

	// Since the canvas is 2D, we use the Z component to store depth information.
	// We store 1/Z so that interpolation preserves depth perspective correctly.
	projected.Pos.Z = zInv
	return canvas.TexVertex{
		Pos:    vertexToPoint(projected.Pos, w, h),
		TexPos: projected.TexPos,
	}
}

func vertexToPoint(v geom.Vec3, width int, height int) geom.Vec3 {
	halfWidth, halfHeight := float32(width)/2, float32(height)/2
	return geom.Vec3{
		X: (1 + v.X) * halfWidth,
		Y: (1 - v.Y) * halfHeight,
		Z: v.Z,
	}
}

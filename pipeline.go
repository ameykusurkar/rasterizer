package main

import (
	"image/color"
	"rasterizer/canvas"
	geom "rasterizer/geometry"
)

// Draw renders the given triangles onto the canvas.
func Draw(canv *canvas.Canvas, triangleList *geom.IndexedTriangleList) {
	vertices := transformVertices(triangleList.Vertices)
	triangles3D := assembleTriangles(vertices, triangleList.Indices)

	width, height := canv.Dimensions()
	clr := color.RGBA{255, 0, 0, 255}

	for _, tri3D := range triangles3D {
		tri2D := [3]geom.Vec2{
			transformPerspective(tri3D[0], width, height),
			transformPerspective(tri3D[1], width, height),
			transformPerspective(tri3D[2], width, height),
		}
		canv.FillTriangle(tri2D[0], tri2D[1], tri2D[2], clr)
	}
}

// Apply any rotation or translation to the vertices if necessary.
func transformVertices(vertices []geom.Vec3) []geom.Vec3 {
	return vertices
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
func transformPerspective(vertex geom.Vec3, width, height int) geom.Vec2 {
	projected := geom.Project(vertex, 1)
	return vertexToPoint(projected, width, height)
}

func vertexToPoint(v geom.Vec3, width int, height int) geom.Vec2 {
	halfWidth, halfHeight := float32(width)/2, float32(height)/2
	x := (1 + v.X) * halfWidth
	y := (1 - v.Y) * halfHeight
	return geom.Vec2{X: float32(int(x)), Y: float32(int(y))}
}

package geometry

// Vertex represents a point in three-dimensional space.
type Vertex struct {
	X, Y, Z float32
}

// Project returns v projected onto the plane Z = d.
func Project(v Vertex, d float32) Vertex {
	return Vertex{
		X: v.X * d / v.Z,
		Y: v.Y * d / v.Z,
		Z: v.Z,
	}
}

// IndexedLineList uses a vertex buffer and an index buffer to represent three-dimensional shapes.
type IndexedLineList struct {
	Vertices []Vertex
	Indices  []int
}

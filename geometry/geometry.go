package geometry

// Vec3 is a three-dimensional geometric object that can be used to represent points in space or vectors.
type Vec3 struct {
	X, Y, Z float32
}

// Project returns v projected onto the plane Z = d.
func Project(v Vec3, d float32) Vec3 {
	return Vec3{
		X: v.X * d / v.Z,
		Y: v.Y * d / v.Z,
		Z: v.Z,
	}
}

// IndexedLineList uses a vertex buffer and an index buffer to represent three-dimensional shapes.
type IndexedLineList struct {
	Vertices []Vec3
	Indices  []int
}

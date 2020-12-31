package geometry

import "math"

// Vec2 represents a point or vector in two-dimensional space.
type Vec2 struct {
	X, Y float32
}

// Vec3 represents a point or vector in three-dimensional space.
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

// IndexedTriangleList represents shapes using triangles.
type IndexedTriangleList struct {
	Vertices []Vec3
	Indices  []int
}

// Mat3 is a 3x3 matrix.
type Mat3 [3][3]float32

// VecMul returns the product of the matrix with v.
func (m *Mat3) VecMul(v Vec3) Vec3 {
	return Vec3{
		X: m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z,
		Y: m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z,
		Z: m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z,
	}
}

// MatMul returns the matrix product of matriv with n.
func (m *Mat3) MatMul(n *Mat3) *Mat3 {
	var product Mat3

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			product[i][j] = m[i][0]*n[0][j] + m[i][1]*n[1][j] + m[i][2]*n[2][j]
		}
	}

	return &product
}

// RotationX returns the rotation matrix around the X-axis.
func RotationX(theta float32) *Mat3 {
	return &Mat3{
		{1, 0, 0},
		{0, cos(theta), -sin(theta)},
		{0, sin(theta), cos(theta)},
	}
}

// RotationY returns the rotation matrix around the Y-axis.
func RotationY(theta float32) *Mat3 {
	return &Mat3{
		{cos(theta), 0, sin(theta)},
		{0, 1, 0},
		{-sin(theta), 0, cos(theta)},
	}
}

// RotationZ returns the rotation matrix around the Z-axis.
func RotationZ(theta float32) *Mat3 {
	return &Mat3{
		{cos(theta), -sin(theta), 0},
		{sin(theta), cos(theta), 0},
		{0, 0, 1},
	}
}

func sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

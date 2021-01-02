package geometry

// Vec3 represents a point or vector in three-dimensional space.
type Vec3 struct {
	X, Y, Z float32
}

// Scale returns the scalar-vector product kv.
func (v Vec3) Scale(k float32) Vec3 {
	return Vec3{
		X: k * v.X,
		Y: k * v.Y,
		Z: k * v.Z,
	}
}

// Sub returns vector v - u.
func (v Vec3) Sub(u Vec3) Vec3 {
	return Vec3{
		X: v.X - u.X,
		Y: v.Y - u.Y,
		Z: v.Z - u.Z,
	}
}

// Add returns vector v + u.
func (v Vec3) Add(u Vec3) Vec3 {
	return Vec3{
		X: v.X + u.X,
		Y: v.Y + u.Y,
		Z: v.Z + u.Z,
	}
}

// Dot returns the dot product of the vector with u.
func (v Vec3) Dot(u Vec3) float32 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}

// Cross returns the cross product of the vector with u.
func (v Vec3) Cross(u Vec3) Vec3 {
	return Vec3{
		X: v.Y*u.Z - v.Z*u.Y,
		Y: v.Z*u.X - v.X*u.Z,
		Z: v.X*u.Y - v.Y*u.X,
	}
}

// InterpolateTo interpolates the vector towards another vector u by step alpha.
func (v Vec3) InterpolateTo(u Vec3, alpha float32) Vec3 {
	return u.Sub(v).Scale(alpha).Add(v)
}

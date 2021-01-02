package geometry

// Vec2 represents a point or vector in two-dimensional space.
type Vec2 struct {
	X, Y float32
}

// Scale returns the scalar-vector product kv.
func (v Vec2) Scale(k float32) Vec2 {
	return Vec2{
		X: k * v.X,
		Y: k * v.Y,
	}
}

// Sub returns vector v - u.
func (v Vec2) Sub(u Vec2) Vec2 {
	return Vec2{
		X: v.X - u.X,
		Y: v.Y - u.Y,
	}
}

// Add returns vector v + u.
func (v Vec2) Add(u Vec2) Vec2 {
	return Vec2{
		X: v.X + u.X,
		Y: v.Y + u.Y,
	}
}

// Interpolate2D returns the interpolation vectors u and v.
func Interpolate2D(u, v Vec2, alpha float32) Vec2 {
  return v.Sub(u).Scale(alpha).Add(u)
}

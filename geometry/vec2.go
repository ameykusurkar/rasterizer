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

// InterpolateTo interpolates the vector towards another vector u by step alpha.
func (v Vec2) InterpolateTo(u Vec2, alpha float32) Vec2 {
  return u.Sub(v).Scale(alpha).Add(v)
}

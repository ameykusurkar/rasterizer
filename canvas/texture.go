package canvas

import (
	"image"
	"image/color"

	geom "rasterizer/geometry"
)

// Texture is a thing, exactly what it is TBD.
type Texture struct {
	Img image.Image
}

func (tex *Texture) shade(v TexVertex) color.Color {
	max := tex.Img.Bounds().Max

	scaledX := int(v.TexPos.X * float32(max.X))
	if scaledX > max.X-1 {
		scaledX = max.X - 1
	}

	scaledY := int(v.TexPos.Y * float32(max.Y))
	if scaledY > max.Y-1 {
		scaledY = max.Y - 1
	}

	return tex.Img.At(scaledX, scaledY)
}

// TexVertex contains a vertex's position both on a two-dimensional surface (eg. a Canvas),
// and its corresponding position on a texture map. This structure is mostly for convinience,
// as we usually need to manipulate both the surface and texture map positions together.
type TexVertex struct {
	// Position of the vertex on a two-dimensional surface.
	Pos geom.Vec2
	// Position of the vertex on the texture map, where 0 <= X, Y < 1.
	TexPos geom.Vec2
}

// Scale returns the scalar-vector product kv.
func (v TexVertex) Scale(k float32) TexVertex {
	return TexVertex{
		Pos:    v.Pos.Scale(k),
		TexPos: v.TexPos.Scale(k),
	}
}

// Sub returns vector v - u.
func (v TexVertex) Sub(u TexVertex) TexVertex {
	return TexVertex{
		Pos:    v.Pos.Sub(u.Pos),
		TexPos: v.TexPos.Sub(u.TexPos),
	}
}

// Add returns vector v + u.
func (v TexVertex) Add(u TexVertex) TexVertex {
	return TexVertex{
		Pos:    v.Pos.Add(u.Pos),
		TexPos: v.TexPos.Add(u.TexPos),
	}
}

// InterpolateTo interpolates the vector towards another vector u by step alpha.
func (v TexVertex) InterpolateTo(u TexVertex, alpha float32) TexVertex {
	return TexVertex{
		Pos:    v.Pos.InterpolateTo(u.Pos, alpha),
		TexPos: v.TexPos.InterpolateTo(u.TexPos, alpha),
	}
}

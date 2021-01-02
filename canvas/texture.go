package canvas

import (
	"image/color"

	geom "rasterizer/geometry"
)

// Texture is a thing, exactly what it is TBD.
type Texture struct {
	Points []geom.Vec2
	Colors []geom.Vec3
}

func (tex *Texture) shade(texCoord geom.Vec2) color.Color {
	colorVec := tex.Colors[0].Scale(texCoord.X).Add(tex.Colors[1].Scale(texCoord.Y))
	return color.RGBA{uint8(colorVec.X), uint8(colorVec.Y), uint8(colorVec.Z), 255}
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

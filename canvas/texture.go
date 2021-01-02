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

// TexVertex represents a vertex than contains both its coordinates in
// 3D space, and its coordinates on a texture map. This allows us to manipulate
// the vertex while also updating its corresponding position on the texture map.
type TexVertex struct {
	// Position of the vertex in 3D space.
	// TODO: 2D or 3D?
	Pos geom.Vec2
	// Position of the vertex on the texture map, where 0 <= X, Y < 1.
	TexPos geom.Vec2
}

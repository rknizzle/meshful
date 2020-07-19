package meshful

import (
	"math"
)

// A triangle is made up of 3 vertices and has a surface normal
type Triangle struct {
	// the 3 vertices that make up a facet
	Vertices [3]Vec3

	// triangle normal vector
	Normal Vec3

	// color of the triangle
	Color *Color
}

func (t *Triangle) SignedVolume() float32 {
	return float32(t.Vertices[0].Dot(t.Vertices[1].Cross(t.Vertices[2])) / 6.0)
}

func (t *Triangle) Area() float32 {
	ab := t.Vertices[0].Diff(t.Vertices[1])
	ac := t.Vertices[0].Diff(t.Vertices[2])
	cross := ab.Cross(ac)

	return float32(0.5 * math.Sqrt(float64(cross.X*cross.X+cross.Y*cross.Y+cross.Z*cross.Z)))
}

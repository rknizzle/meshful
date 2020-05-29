package meshful

// A triangle is made up of 3 vertices and has a surface normal
type Triangle struct {
	// the 3 vertices that make up a facet
	Vertices [3]Vec3

	// triangle normal vector
	Normal Vec3

	// color of the triangle
	Color *Color
}

func (t *Triangle) GetSignedVolume() float32 {
	return float32(t.Vertices[0].Dot(t.Vertices[1].Cross(t.Vertices[2])) / 6.0)
}

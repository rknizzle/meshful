package meshful

// A triangle is made up of 3 vertices and has a surface normal
type Triangle struct {
	// the 3 vertices that make up a facet
	Vertices [3]Vec3

	// triangle normal vector
	Normal Vec3

	// extra info stored only in stl binary files
	Extra uint16
}

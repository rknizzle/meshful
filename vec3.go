package meshful

// A Vec3 represents a 3 dimensional vector for storing triangle vertices and normals
type Vec3 struct {
	X, Y, Z float32
}

// get the cross product of two vectors
func (vec Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		vec.Y*other.Z - vec.Z*other.Y,
		vec.Z*other.X - vec.X*other.Z,
		vec.X*other.Y - vec.Y*other.X,
	}
}

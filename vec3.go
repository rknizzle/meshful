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

// get the dot product of two vectors
func (vec Vec3) Dot(other Vec3) float64 {
	return float64(vec.X)*float64(other.X) +
		float64(vec.Y)*float64(other.Y) +
		float64(vec.Z)*float64(other.Z)
}

// diff returns the difference between two vectors
func (vec Vec3) Diff(other Vec3) Vec3 {
	return Vec3{
		vec.X - other.X,
		vec.Y - other.Y,
		vec.Z - other.Z,
	}
}

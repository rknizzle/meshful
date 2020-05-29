package meshful

// A mesh represents a collection of triangles
type Mesh struct {
	Triangles []Triangle
}

func (mesh *Mesh) GetBoundingBox() [3]float32 {
	minX := mesh.Triangles[0].Vertices[0].X
	maxX := mesh.Triangles[0].Vertices[0].X
	minY := mesh.Triangles[0].Vertices[0].Y
	maxY := mesh.Triangles[0].Vertices[0].Y
	minZ := mesh.Triangles[0].Vertices[0].Z
	maxZ := mesh.Triangles[0].Vertices[0].Z

	// loop through every vertex in the mesh
	for _, triangle := range mesh.Triangles {
		for _, vert := range triangle.Vertices {
			if vert.X < minX {
				minX = vert.X
			}
			if vert.X > maxX {
				maxX = vert.X
			}
			if vert.Y < minY {
				minY = vert.Y
			}
			if vert.Y > maxY {
				maxY = vert.Y
			}
			if vert.Z < minZ {
				minZ = vert.Z
			}
			if vert.Z > maxZ {
				maxZ = vert.Z
			}
		}
	}

	return [3]float32{maxX - minX, maxY - minY, maxZ - minZ}
}

func (mesh *Mesh) GetVolume() float32 {
	var volume float32 = 0.0
	// loop through each triangle and add up the volume of each to get the meshes volume
	for _, triangle := range mesh.Triangles {
		volume += triangle.GetSignedVolume()
	}
	return volume
}

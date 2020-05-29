package meshful

import (
	"testing"
)

func makeTestMesh() *Mesh {
	return &Mesh{
		Triangles: []Triangle{
			{
				Normal: Vec3{0, 0, -1},
				Vertices: [3]Vec3{
					{0, 0, 0},
					{0, 1, 0},
					{1, 0, 0},
				},
			},
			{
				Normal: Vec3{0, -1, 0},
				Vertices: [3]Vec3{
					{0, 0, 0},
					{1, 0, 0},
					{0, 0, 1},
				},
			},
			{
				Normal: Vec3{0.57735, 0.57735, 0.57735},
				Vertices: [3]Vec3{
					{0, 0, 1},
					{1, 0, 0},
					{0, 1, 0},
				},
			},
			{
				Normal: Vec3{-1, 0, 0},
				Vertices: [3]Vec3{
					{0, 0, 0},
					{0, 0, 1},
					{0, 1, 0},
				},
			},
		},
	}
}

func TestBoundingBox(t *testing.T) {
	mesh := makeTestMesh()
	bbox := mesh.GetBoundingBox()

	if bbox != ([3]float32{1, 1, 1}) {
		t.Errorf("Expected bounding box of [1 1 1], found: %v", bbox)
	}
}

func TestVolume(t *testing.T) {
	mesh := makeTestMesh()
	volume := mesh.GetVolume()
	if volume <= 0 {
		t.Errorf("Expected positive non-zero volume")
	}
}

func TestArea(t *testing.T) {
	mesh := makeTestMesh()
	area := mesh.GetSurfaceArea()
	if area <= 0 {
		t.Errorf("Expected positive non-zero surface area")
	}
}

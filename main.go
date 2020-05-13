package meshful

import (
	"github.com/rknizzle/meshful/io/stl"
	"github.com/rknizzle/meshful/mesh"
)

func Test() (*mesh.Mesh, error) {
	m, err := stl.Read()
	if err != nil {
		return nil, err
	}
	return m, nil
}

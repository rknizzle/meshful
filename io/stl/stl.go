package stl

import (
	"errors"
	"github.com/rknizzle/meshful/mesh"
)

func Read() (*mesh.Mesh, error) {
	return &mesh.Mesh{}, nil
}

func Write() error {
	return errors.New("Not yet implemented")
}

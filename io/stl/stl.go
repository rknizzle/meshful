package stl

import (
	"errors"
	"github.com/rknizzle/meshful"
)

func Read() (*meshful.Mesh, error) {
	return &meshful.Mesh{}, nil
}

func Write() error {
	return errors.New("Not yet implemented")
}

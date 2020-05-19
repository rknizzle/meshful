package main

import (
	"fmt"
	"github.com/rknizzle/meshful"
	"github.com/rknizzle/meshful/io/stl"
)

func main() {
	var m *meshful.Mesh
	m, err := stl.ReadFile("cmd/testdata/binary_cube.stl")
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
	err = stl.WriteFile("new.stl", m)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done")
}

package main

import (
	"fmt"
	"github.com/rknizzle/meshful"
	"github.com/rknizzle/meshful/io/stl"
)

func main() {
	var m *meshful.Mesh
	m, _ = stl.Read()
	fmt.Println(m)
}

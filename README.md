# meshful
#### Library for processing 3d triangle meshes

## Features
Read/Write STL and OBJ files and get mesh dimensions

## Demo snippet:
load in an STL file, get the dimensions and export as an OBJ.
``` go
package main

import (
	"fmt"
	"github.com/rknizzle/meshful"
	"github.com/rknizzle/meshful/io/obj"
	"github.com/rknizzle/meshful/io/stl"
)

func main() {
	var m *meshful.Mesh
	m, err := stl.ReadFile("tetrahedron.stl")
	if err != nil {
		panic(err)
	}

	fmt.Println("Dimensions of mesh:")
	fmt.Println("bbox:", m.BoundingBox())
	fmt.Println("volume:", m.Volume())
	fmt.Println("area:", m.SurfaceArea())

	err = obj.WriteFile("new.obj", m)
	if err != nil {
		panic(err)
	}
}
```

## WIP: This repo is a work in progress
#### TODO:
- transform mesh(scale, rotate, move)
- Support more complex OBJs
- Mesh health analysis and repair

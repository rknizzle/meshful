package obj

import (
	"bufio"
	"github.com/rknizzle/meshful"
	"io"
	"os"
)

// Readfile reads the contents of a Wavefront OBJ file into a new Mesh object
func ReadFile(filename string) (mesh *meshful.Mesh, err error) {
	file, openErr := os.Open(filename)
	if openErr != nil {
		err = openErr
		return
	}
	defer file.Close()

	return readAll(bufio.NewReader(file))
}

func readAll(r io.Reader) (mesh *meshful.Mesh, err error) {
	// read through each line and generate the mesh

	//vertices := []meshful.Vec3{}
	// read each vertex and add it to this vertices array

	//triangles := []meshful.Triangle{}
	// read through each face and the vertex number in the OBJ goes along with the index in the vertices array

	// and then just calculate each normal I think
	// and this triangles array becomes the mesh

	return &meshful.Mesh{}, nil
}

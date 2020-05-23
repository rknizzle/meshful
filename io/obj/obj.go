package obj

import (
	"bufio"
	"fmt"
	"github.com/rknizzle/meshful"
	"io"
	"os"
	"strings"
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
	scanner := bufio.NewScanner(r)

	// keep a list of all the vertices and faces specified in the file
	vertices := []meshful.Vec3{}
	faces := []meshful.Triangle{}

	// loop through each line of the file
	for scanner.Scan() {
		line := scanner.Text()

		// skip blank lines
		if line == "" {
			continue
		}

		// the first word of each line should be a token specifying the data type of that line
		words := strings.Fields(line)
		token := words[0]

		if token == "#" {
			// Its just a comment, continue to next line
			continue
		}
		if token == "v" {
			// new vertex -- get the coordinates and add it to the list of vertices
			v := meshful.Vec3{}
			vertices = append(vertices, v)
		}
		if token == "f" {
			// new face -- construct the face using the list of vertices
			f := meshful.Triangle{}
			faces = append(faces, f)
		}

		fmt.Println(token)
	}

	return &meshful.Mesh{}, nil
}
